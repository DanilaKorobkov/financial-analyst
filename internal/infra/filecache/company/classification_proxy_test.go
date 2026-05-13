package company_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/suite"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	fccompany "github.com/DanilaKorobkov/financial-analyst/internal/infra/filecache/company"
	company_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/company"
)

type classificationProxySuite struct {
	suite.Suite

	delegate *company_mock.ClassificationGateway
	proxy    *fccompany.ClassificationProxy
	dir      string
}

// envelopeJSON — JSON-раскладка файла кеша, симметричная prod-типам
// classificationEnvelope и classificationDTO. Описана здесь явно,
// чтобы тест ломался, если prod-раскладка молча разъедется.
type envelopeJSON struct {
	ExpiresAt      time.Time          `json:"expires_at"`
	Classification classificationJSON `json:"classification"`
}

type classificationJSON struct {
	Sector              string                 `json:"sector"`
	Industry            string                 `json:"industry"`
	Country             string                 `json:"country"`
	PrimaryReportTicker string                 `json:"primary_report_ticker"`
	Exchange            domaincompany.Exchange `json:"exchange"`
	Currency            domaincompany.Currency `json:"currency"`
}

// toDomain сворачивает прочитанную JSON-проекцию обратно в
// domain-значение для сравнений в s.Equal.
func (c *classificationJSON) toDomain() domaincompany.Classification {
	return domaincompany.Classification{
		Sector:              c.Sector,
		Industry:            c.Industry,
		Country:             c.Country,
		PrimaryReportTicker: c.PrimaryReportTicker,
		Exchange:            c.Exchange,
		Currency:            c.Currency,
	}
}

func TestClassificationProxySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(classificationProxySuite))
}

func (s *classificationProxySuite) SetupTest() {
	s.dir = s.T().TempDir()
	s.delegate = company_mock.NewClassificationGateway(s.T())
	s.proxy = fccompany.NewClassificationProxy(fccompany.ConfigClassificationProxy{
		Delegate: s.delegate,
		Dir:      s.dir,
	})
}

func (s *classificationProxySuite) TestFindByTickerCacheMissPopulatesFile() {
	cls := sberClassification()
	s.delegate.EXPECT().
		FindByTicker(context.Background(), "SBER").
		Return(cls, nil).
		Once()

	got, err := s.proxy.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(cls, got)
	s.FileExists(filepath.Join(s.dir, "SBER.json"))
}

func (s *classificationProxySuite) TestFindByTickerCacheHitSkipsDelegate() {
	cls := sberClassification()
	s.delegate.EXPECT().
		FindByTicker(context.Background(), "SBER").
		Return(cls, nil).
		Once()

	_, err := s.proxy.FindByTicker(context.Background(), "SBER")
	s.Require().NoError(err)

	// Второй вызов идёт из кеша — делегат не должен быть позван повторно.
	got, err := s.proxy.FindByTicker(context.Background(), "SBER")
	s.Require().NoError(err)
	s.Equal(cls, got)
}

func (s *classificationProxySuite) TestFindByTickerDelegateErrorNotCached() {
	s.delegate.EXPECT().
		FindByTicker(context.Background(), "GAZP").
		Return(domaincompany.Classification{}, domaincompany.ErrNotFound).
		Twice()

	_, err := s.proxy.FindByTicker(context.Background(), "GAZP")
	s.Require().ErrorIs(err, domaincompany.ErrNotFound)

	// Повторный вызов должен снова уйти в делегат — отрицательных записей нет.
	_, err = s.proxy.FindByTicker(context.Background(), "GAZP")
	s.Require().ErrorIs(err, domaincompany.ErrNotFound)

	entries, readErr := os.ReadDir(s.dir)
	s.Require().NoError(readErr)
	s.Empty(entries, "ошибки делегата не пишутся в кеш")
}

func (s *classificationProxySuite) TestFindByTickerCorruptedCacheFallsBackToDelegate() {
	key := filepath.Join(s.dir, "SBER.json")
	s.Require().NoError(os.WriteFile(key, []byte("{not-json"), 0o600))

	cls := sberClassification()
	s.delegate.EXPECT().
		FindByTicker(context.Background(), "SBER").
		Return(cls, nil).
		Once()

	got, err := s.proxy.FindByTicker(context.Background(), "SBER")
	s.Require().NoError(err)
	s.Equal(cls, got)
}

// TestFindByTickerUnsafeTickerStillCached проверяет, что тикер с
// разделителями пути и точками всё равно проходит через кеш:
// url.PathEscape делает ключ безопасным для diskv, файл лежит строго
// внутри Dir, а повторный запрос обслуживается без обращения к делегату.
func (s *classificationProxySuite) TestFindByTickerUnsafeTickerStillCached() {
	const unsafe = "../etc/passwd"
	cls := sberClassification()
	s.delegate.EXPECT().
		FindByTicker(context.Background(), unsafe).
		Return(cls, nil).
		Once()

	_, err := s.proxy.FindByTicker(context.Background(), unsafe)
	s.Require().NoError(err)

	// Файл материализовался строго внутри Dir — наружу не вылез.
	entries, readErr := os.ReadDir(s.dir)
	s.Require().NoError(readErr)
	s.Len(entries, 1)
	wantName := fmt.Sprintf("%s.json", url.PathEscape(unsafe))
	s.Equal(wantName, entries[0].Name())

	// Повторный вызов идёт из кеша — делегат не зовётся повторно.
	got, err := s.proxy.FindByTicker(context.Background(), unsafe)
	s.Require().NoError(err)
	s.Equal(cls, got)
}

// TestFindByTickerExpiredEntryRefreshed: при cache hit с истёкшим
// ExpiresAt запись считается просроченной — Proxy идёт в делегат
// и перезаписывает файл новой картой и сдвинутым ExpiresAt. Время
// контролируется testing/synctest: внутри bubble time.Now() и
// time.Sleep работают на синтетических часах bubble.
func (s *classificationProxySuite) TestFindByTickerExpiredEntryRefreshed() {
	synctest.Test(s.T(), func(_ *testing.T) {
		const ttl = time.Hour
		dir := s.T().TempDir()
		delegate := company_mock.NewClassificationGateway(s.T())
		proxy := fccompany.NewClassificationProxy(fccompany.ConfigClassificationProxy{
			Delegate: delegate,
			Dir:      dir,
			TTL:      ttl,
		})

		first := sberClassification()
		second := sberClassification()
		second.Industry = "Обновлённая отрасль"
		delegate.EXPECT().
			FindByTicker(context.Background(), "SBER").
			Return(first, nil).
			Once()
		delegate.EXPECT().
			FindByTicker(context.Background(), "SBER").
			Return(second, nil).
			Once()

		writtenAt := time.Now().UTC()
		got, err := proxy.FindByTicker(context.Background(), "SBER")
		s.Require().NoError(err)
		s.Equal(first, got)
		s.Equal(writtenAt.Add(ttl), readEnvelope(s.T(), dir, "SBER.json").ExpiresAt)

		// Сдвигаемся за пределы TTL — следующий вызов должен идти в делегат.
		time.Sleep(ttl + time.Second)
		refreshedAt := time.Now().UTC()
		got, err = proxy.FindByTicker(context.Background(), "SBER")
		s.Require().NoError(err)
		s.Equal(second, got)
		envelope := readEnvelope(s.T(), dir, "SBER.json")
		s.Equal(second, envelope.Classification.toDomain())
		s.Equal(refreshedAt.Add(ttl), envelope.ExpiresAt)
	})
}

// TestFindByTickerWithinTTLServesFromCache: пока now < ExpiresAt, кеш
// считается живым и делегат не зовётся.
func (s *classificationProxySuite) TestFindByTickerWithinTTLServesFromCache() {
	synctest.Test(s.T(), func(_ *testing.T) {
		dir := s.T().TempDir()
		delegate := company_mock.NewClassificationGateway(s.T())
		proxy := fccompany.NewClassificationProxy(fccompany.ConfigClassificationProxy{
			Delegate: delegate,
			Dir:      dir,
			TTL:      time.Hour,
		})

		cls := sberClassification()
		delegate.EXPECT().
			FindByTicker(context.Background(), "SBER").
			Return(cls, nil).
			Once()

		_, err := proxy.FindByTicker(context.Background(), "SBER")
		s.Require().NoError(err)

		// Доходим почти до конца TTL — запись должна продолжать жить.
		time.Sleep(time.Hour - time.Second)
		got, err := proxy.FindByTicker(context.Background(), "SBER")
		s.Require().NoError(err)
		s.Equal(cls, got)
	})
}

// TestFindByTickerZeroTTLNeverExpires: TTL == 0 означает «без
// экспирации» — ExpiresAt в файле нулевой, и сколько бы часы ни шли,
// запись всегда считается живой.
func (s *classificationProxySuite) TestFindByTickerZeroTTLNeverExpires() {
	synctest.Test(s.T(), func(_ *testing.T) {
		dir := s.T().TempDir()
		delegate := company_mock.NewClassificationGateway(s.T())
		proxy := fccompany.NewClassificationProxy(fccompany.ConfigClassificationProxy{
			Delegate: delegate,
			Dir:      dir,
		})

		cls := sberClassification()
		delegate.EXPECT().
			FindByTicker(context.Background(), "SBER").
			Return(cls, nil).
			Once()

		_, err := proxy.FindByTicker(context.Background(), "SBER")
		s.Require().NoError(err)
		s.True(readEnvelope(s.T(), dir, "SBER.json").ExpiresAt.IsZero())

		time.Sleep(100 * 365 * 24 * time.Hour)
		got, err := proxy.FindByTicker(context.Background(), "SBER")
		s.Require().NoError(err)
		s.Equal(cls, got)
	})
}

func sberClassification() domaincompany.Classification {
	return domaincompany.Classification{
		Exchange:            domaincompany.ExchangeMOEX,
		Currency:            domaincompany.CurrencyRUB,
		Sector:              "Финансы",
		Industry:            "Банковская деятельность",
		Country:             "Россия",
		PrimaryReportTicker: "SBER",
	}
}

func readEnvelope(t *testing.T, dir, name string) envelopeJSON {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, name)) //nolint:gosec // dir = t.TempDir(), name — литерал теста
	if err != nil {
		t.Fatalf("read envelope file: %v", err)
	}
	var envelope envelopeJSON
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	return envelope
}
