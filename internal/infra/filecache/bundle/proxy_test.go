package bundle_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/suite"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	fcbundle "github.com/DanilaKorobkov/financial-analyst/internal/infra/filecache/bundle"
)

// fakeBundle — простая реализация data.Bundle для тестов Proxy. Позволяет
// подменять Fetch на скрипт ответов, не таская в тесты mockery — подмена
// идёт только на границе cache↔delegate, без HTTP.
type fakeBundle struct {
	bundleID string
	fields   []data.FieldDescriptor
	scripts  []fakeScript
	called   int
}

type fakeScript struct {
	values data.FieldValues
	err    error
}

type proxySuite struct {
	suite.Suite

	delegate *fakeBundle
	proxy    *fcbundle.Proxy
	dir      string
}

// envelopeFile — то, что лежит на диске в файле кеша. Дублирует
// приватный envelopeOnDisk Proxy ровно на ту глубину, на которой тестам
// нужно проверить раскладку JSON.
type envelopeFile struct {
	ExpiresAt time.Time                  `json:"expires_at"`
	Values    map[string]json.RawMessage `json:"values"`
}

func (b *fakeBundle) BundleID() string               { return b.bundleID }
func (b *fakeBundle) Fields() []data.FieldDescriptor { return b.fields }
func (b *fakeBundle) Fetch(_ context.Context, _ string) (data.FieldValues, error) {
	if b.called >= len(b.scripts) {
		return nil, fmt.Errorf("fakeBundle: unexpected call #%d", b.called+1)
	}
	s := b.scripts[b.called]
	b.called++
	return s.values, s.err
}

func newFakeBundle(scripts ...fakeScript) *fakeBundle {
	return &fakeBundle{
		bundleID: "stock-info",
		fields: []data.FieldDescriptor{
			{ID: domaincompany.FieldIssuerName, Type: data.TypeString, Description: "Название эмитента."},
			{ID: domaincompany.FieldIndustry, Type: data.TypeString, Description: "Отрасль."},
			{ID: domaincompany.FieldSectorID, Type: data.TypeInt64, Description: "Код сектора."},
			{ID: domaincompany.FieldExchange, Type: data.TypeExchange, Description: "Биржа."},
			{ID: domaincompany.FieldCurrency, Type: data.TypeCurrency, Description: "Валюта."},
			{ID: domaincompany.FieldReportFrequency, Type: data.TypeReportFrequency, Description: "Частота отчётности."},
			{ID: domaincompany.FieldSPB, Type: data.TypeBool, Description: "Листинг СПБ."},
			{ID: domaincompany.FieldIssueDate, Type: data.TypeDate, Description: "Дата начала торгов."},
			{ID: domaincompany.FieldSecurityType, Type: data.TypeSecurityType, Description: "Тип бумаги."},
			{ID: domaincompany.FieldListingLevel, Type: data.TypeListingLevel, Description: "Котировальный уровень."},
		},
		scripts: scripts,
	}
}

func sberValues() data.FieldValues {
	return data.FieldValues{
		domaincompany.FieldIssuerName:      "Сбербанк",
		domaincompany.FieldIndustry:        "Банковская деятельность",
		domaincompany.FieldSectorID:        int64(40),
		domaincompany.FieldExchange:        domaincompany.ExchangeMOEX,
		domaincompany.FieldCurrency:        domaincompany.CurrencyRUB,
		domaincompany.FieldReportFrequency: domaincompany.ReportFrequencyQuarterly,
		domaincompany.FieldSPB:             false,
		domaincompany.FieldIssueDate:       time.Date(2007, 7, 20, 0, 0, 0, 0, time.UTC),
		domaincompany.FieldSecurityType:    domaincompany.SecurityTypeCommonShare,
		domaincompany.FieldListingLevel:    domaincompany.ListingLevelFirst,
	}
}

func TestProxySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(proxySuite))
}

func (s *proxySuite) buildProxy(scripts ...fakeScript) {
	s.dir = s.T().TempDir()
	s.delegate = newFakeBundle(scripts...)
	s.proxy = fcbundle.NewProxy(fcbundle.ConfigProxy{
		Delegate: s.delegate,
		Dir:      s.dir,
	})
}

func (s *proxySuite) TestMetadataMirrorsDelegate() {
	s.buildProxy()
	s.Equal("stock-info", s.proxy.BundleID())
	s.Equal(s.delegate.Fields(), s.proxy.Fields())
}

func (s *proxySuite) TestFetchCacheMissPopulatesFile() {
	s.buildProxy(fakeScript{values: sberValues()})

	got, err := s.proxy.Fetch(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(sberValues(), got)
	s.FileExists(filepath.Join(s.dir, "SBER.json"))
}

func (s *proxySuite) TestFetchCacheHitSkipsDelegate() {
	s.buildProxy(fakeScript{values: sberValues()})

	_, err := s.proxy.Fetch(context.Background(), "SBER")
	s.Require().NoError(err)

	// Второй вызов идёт из кеша — delegate не должен быть позван повторно.
	got, err := s.proxy.Fetch(context.Background(), "SBER")
	s.Require().NoError(err)
	s.Equal(sberValues(), got)
	s.Equal(1, s.delegate.called)
}

func (s *proxySuite) TestFetchDelegateErrorNotCached() {
	s.buildProxy(
		fakeScript{err: domaincompany.ErrNotFound},
		fakeScript{err: domaincompany.ErrNotFound},
	)

	_, err := s.proxy.Fetch(context.Background(), "missing")
	s.Require().ErrorIs(err, domaincompany.ErrNotFound)

	// Повторный вызов должен снова уйти в delegate — отрицательных записей нет.
	_, err = s.proxy.Fetch(context.Background(), "missing")
	s.Require().ErrorIs(err, domaincompany.ErrNotFound)

	entries, readErr := os.ReadDir(s.dir)
	s.Require().NoError(readErr)
	s.Empty(entries, "ошибки делегата не пишутся в кеш")
}

func (s *proxySuite) TestFetchCorruptedCacheFallsBackToDelegate() {
	s.buildProxy(fakeScript{values: sberValues()})

	key := filepath.Join(s.dir, "SBER.json")
	s.Require().NoError(os.WriteFile(key, []byte("{not-json"), 0o600))

	got, err := s.proxy.Fetch(context.Background(), "SBER")
	s.Require().NoError(err)
	s.Equal(sberValues(), got)
}

// TestFetchUnsafeTickerStillCached проверяет, что тикер с разделителями
// пути и точками всё равно проходит через кеш: url.PathEscape делает
// ключ безопасным для diskv, файл лежит строго внутри Dir, а повторный
// запрос обслуживается без обращения к delegate.
func (s *proxySuite) TestFetchUnsafeTickerStillCached() {
	const unsafe = "../etc/passwd"
	s.buildProxy(fakeScript{values: sberValues()})

	_, err := s.proxy.Fetch(context.Background(), unsafe)
	s.Require().NoError(err)

	entries, readErr := os.ReadDir(s.dir)
	s.Require().NoError(readErr)
	s.Len(entries, 1)
	wantName := fmt.Sprintf("%s.json", url.PathEscape(unsafe))
	s.Equal(wantName, entries[0].Name())

	got, err := s.proxy.Fetch(context.Background(), unsafe)
	s.Require().NoError(err)
	s.Equal(sberValues(), got)
	s.Equal(1, s.delegate.called)
}

// TestFetchExpiredEntryRefreshed: при cache hit с истёкшим ExpiresAt
// запись считается просроченной — Proxy идёт в delegate и перезаписывает
// файл новой картой и сдвинутым ExpiresAt. Время контролируется
// testing/synctest.
func (s *proxySuite) TestFetchExpiredEntryRefreshed() {
	synctest.Test(s.T(), func(_ *testing.T) {
		const ttl = time.Hour
		dir := s.T().TempDir()
		first := sberValues()
		second := sberValues()
		second[domaincompany.FieldIndustry] = "Обновлённая отрасль"
		delegate := newFakeBundle(
			fakeScript{values: first},
			fakeScript{values: second},
		)
		proxy := fcbundle.NewProxy(fcbundle.ConfigProxy{
			Delegate: delegate,
			Dir:      dir,
			TTL:      ttl,
		})

		writtenAt := time.Now().UTC()
		got, err := proxy.Fetch(context.Background(), "SBER")
		s.Require().NoError(err)
		s.Equal(first, got)
		s.Equal(writtenAt.Add(ttl), readEnvelope(s.T(), dir).ExpiresAt)

		// Сдвигаемся за пределы TTL — следующий вызов должен идти в delegate.
		time.Sleep(ttl + time.Second)
		refreshedAt := time.Now().UTC()
		got, err = proxy.Fetch(context.Background(), "SBER")
		s.Require().NoError(err)
		s.Equal(second, got)
		envelope := readEnvelope(s.T(), dir)
		s.Equal(refreshedAt.Add(ttl), envelope.ExpiresAt)
	})
}

// TestFetchWithinTTLServesFromCache: пока now < ExpiresAt, кеш считается
// живым и delegate не зовётся.
func (s *proxySuite) TestFetchWithinTTLServesFromCache() {
	synctest.Test(s.T(), func(_ *testing.T) {
		dir := s.T().TempDir()
		delegate := newFakeBundle(fakeScript{values: sberValues()})
		proxy := fcbundle.NewProxy(fcbundle.ConfigProxy{
			Delegate: delegate,
			Dir:      dir,
			TTL:      time.Hour,
		})

		_, err := proxy.Fetch(context.Background(), "SBER")
		s.Require().NoError(err)

		time.Sleep(time.Hour - time.Second)
		got, err := proxy.Fetch(context.Background(), "SBER")
		s.Require().NoError(err)
		s.Equal(sberValues(), got)
	})
}

// TestFetchZeroTTLNeverExpires: TTL == 0 означает «без экспирации» —
// ExpiresAt в файле нулевой, и сколько бы часы ни шли, запись всегда
// считается живой.
func (s *proxySuite) TestFetchZeroTTLNeverExpires() {
	synctest.Test(s.T(), func(_ *testing.T) {
		dir := s.T().TempDir()
		delegate := newFakeBundle(fakeScript{values: sberValues()})
		proxy := fcbundle.NewProxy(fcbundle.ConfigProxy{
			Delegate: delegate,
			Dir:      dir,
		})

		_, err := proxy.Fetch(context.Background(), "SBER")
		s.Require().NoError(err)
		s.True(readEnvelope(s.T(), dir).ExpiresAt.IsZero())

		time.Sleep(100 * 365 * 24 * time.Hour)
		got, err := proxy.Fetch(context.Background(), "SBER")
		s.Require().NoError(err)
		s.Equal(sberValues(), got)
	})
}

// TestFetchUnsupportedFieldTypeFallsBackToDelegate: bundle декларирует
// FieldType вне известного набора — Encode пройдёт (jsoniter маршалит
// любое значение), но Decode при cache hit упадёт на default-ветке
// decodeByType. Proxy расценит это как «битый кеш» и снова сходит
// в delegate. Покрытие defensive default-ветки decodeByType через
// публичный API.
func (s *proxySuite) TestFetchUnsupportedFieldTypeFallsBackToDelegate() {
	const id = "synthetic::unknown-type"
	dir := s.T().TempDir()
	bundle := &fakeBundle{
		bundleID: "synthetic",
		fields:   []data.FieldDescriptor{{ID: id, Type: data.FieldType(9999)}},
		scripts: []fakeScript{
			{values: data.FieldValues{id: "value"}},
			{values: data.FieldValues{id: "value"}},
		},
	}
	proxy := fcbundle.NewProxy(fcbundle.ConfigProxy{Delegate: bundle, Dir: dir})

	_, err := proxy.Fetch(context.Background(), "any")
	s.Require().NoError(err)

	_, err = proxy.Fetch(context.Background(), "any")
	s.Require().NoError(err)
	// Оба раза должны были сходить в delegate: cache hit упал на decode,
	// proxy свалился на delegate.
	s.Equal(2, bundle.called)
}

// TestFetchCacheValueWrongTypeFallsBackToDelegate: envelope валидный,
// но конкретное поле в кеше имеет JSON неверного типа (строка вместо
// int). decodeByType вернёт ошибку конкретного поля, и proxy свалится
// на delegate.
func (s *proxySuite) TestFetchCacheValueWrongTypeFallsBackToDelegate() {
	s.buildProxy(fakeScript{values: sberValues()}, fakeScript{values: sberValues()})

	_, err := s.proxy.Fetch(context.Background(), "SBER")
	s.Require().NoError(err)

	key := filepath.Join(s.dir, "SBER.json")
	envelope := readEnvelope(s.T(), s.dir)
	envelope.Values[string(domaincompany.FieldSectorID)] = json.RawMessage(`"not-a-number"`)
	raw, mErr := json.Marshal(envelope)
	s.Require().NoError(mErr)
	s.Require().NoError(os.WriteFile(key, raw, 0o600))

	_, err = s.proxy.Fetch(context.Background(), "SBER")
	s.Require().NoError(err)
	s.Equal(2, s.delegate.called)
}

// TestFetchDelegateGenericError проверяет, что не-ErrNotFound-ошибка
// delegate так же не пишется в кеш и пробрасывается наверх.
func (s *proxySuite) TestFetchDelegateGenericError() {
	sentinel := errors.New("upstream boom")
	s.buildProxy(fakeScript{err: sentinel})

	_, err := s.proxy.Fetch(context.Background(), "any")
	s.Require().ErrorIs(err, sentinel)

	entries, readErr := os.ReadDir(s.dir)
	s.Require().NoError(readErr)
	s.Empty(entries)
}

func readEnvelope(t *testing.T, dir string) envelopeFile {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, "SBER.json")) //nolint:gosec // dir = t.TempDir()
	if err != nil {
		t.Fatalf("read envelope file: %v", err)
	}
	var envelope envelopeFile
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	return envelope
}
