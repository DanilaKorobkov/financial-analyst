package companycard_test

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

	domaincard "github.com/DanilaKorobkov/financial-analyst/internal/domain/companycard"
	fccard "github.com/DanilaKorobkov/financial-analyst/internal/infra/filecache/companycard"
	companycard_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/companycard"
)

type repositorySuite struct {
	suite.Suite

	delegate *companycard_mock.Repository
	repo     *fccard.Repository
	dir      string
}

// envelopeOnDisk — формат записи, ожидаемый в файле кеша. Дублирует
// приватный тип репозитория ровно настолько, чтобы тесты могли убедиться
// в раскладке JSON и значении ExpiresAt.
type envelopeOnDisk struct {
	ExpiresAt time.Time       `json:"expires_at"`
	Card      domaincard.Card `json:"card"`
}

func TestRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) SetupTest() {
	s.dir = s.T().TempDir()
	s.delegate = companycard_mock.NewRepository(s.T())
	s.repo = fccard.NewRepository(fccard.ConfigRepository{
		Delegate: s.delegate,
		Dir:      s.dir,
	})
}

func (s *repositorySuite) TestFindByTickerCacheMissPopulatesFile() {
	card := sberCard()
	s.delegate.EXPECT().
		FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER").
		Return(card, nil).
		Once()

	got, err := s.repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER")

	s.Require().NoError(err)
	s.Equal(card, got)
	s.FileExists(filepath.Join(s.dir, "1_SBER.json"))
}

func (s *repositorySuite) TestFindByTickerCacheHitSkipsDelegate() {
	card := sberCard()
	s.delegate.EXPECT().
		FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER").
		Return(card, nil).
		Once()

	_, err := s.repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER")
	s.Require().NoError(err)

	// Второй вызов идёт из кеша — делегат не должен быть позван повторно.
	got, err := s.repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER")
	s.Require().NoError(err)
	s.Equal(card, got)
}

func (s *repositorySuite) TestFindByTickerDelegateErrorNotCached() {
	s.delegate.EXPECT().
		FindByTicker(context.Background(), domaincard.ExchangeMOEX, "GAZP").
		Return(domaincard.Card{}, domaincard.ErrNotFound).
		Twice()

	_, err := s.repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "GAZP")
	s.Require().ErrorIs(err, domaincard.ErrNotFound)

	// Повторный вызов должен снова уйти в делегат — отрицательных записей нет.
	_, err = s.repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "GAZP")
	s.Require().ErrorIs(err, domaincard.ErrNotFound)

	entries, readErr := os.ReadDir(s.dir)
	s.Require().NoError(readErr)
	s.Empty(entries, "ошибки делегата не пишутся в кеш")
}

func (s *repositorySuite) TestFindByTickerCorruptedCacheFallsBackToDelegate() {
	key := filepath.Join(s.dir, "1_SBER.json")
	s.Require().NoError(os.WriteFile(key, []byte("{not-json"), 0o600))

	card := sberCard()
	s.delegate.EXPECT().
		FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER").
		Return(card, nil).
		Once()

	got, err := s.repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER")
	s.Require().NoError(err)
	s.Equal(card, got)
}

// TestFindByTickerUnsafeTickerStillCached проверяет, что тикер с
// разделителями пути и точками всё равно проходит через кеш:
// url.PathEscape делает ключ безопасным для diskv, файл лежит строго
// внутри Dir, а повторный запрос обслуживается без обращения к делегату.
func (s *repositorySuite) TestFindByTickerUnsafeTickerStillCached() {
	const unsafe = "../etc/passwd"
	card := sberCard()
	s.delegate.EXPECT().
		FindByTicker(context.Background(), domaincard.ExchangeMOEX, unsafe).
		Return(card, nil).
		Once()

	_, err := s.repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, unsafe)
	s.Require().NoError(err)

	// Файл материализовался строго внутри Dir — наружу не вылез.
	entries, readErr := os.ReadDir(s.dir)
	s.Require().NoError(readErr)
	s.Len(entries, 1)
	wantName := fmt.Sprintf("%d_%s.json", int(domaincard.ExchangeMOEX), url.PathEscape(unsafe))
	s.Equal(wantName, entries[0].Name())

	// Повторный вызов идёт из кеша — делегат не зовётся повторно.
	got, err := s.repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, unsafe)
	s.Require().NoError(err)
	s.Equal(card, got)
}

// TestFindByTickerExpiredEntryRefreshed: при cache hit с истёкшим
// ExpiresAt запись считается просроченной — Repository идёт в делегат
// и перезаписывает файл новой картой и сдвинутым ExpiresAt. Время
// контролируется testing/synctest: внутри bubble time.Now() и
// time.Sleep работают на синтетических часах bubble.
func (s *repositorySuite) TestFindByTickerExpiredEntryRefreshed() {
	synctest.Test(s.T(), func(_ *testing.T) {
		const ttl = time.Hour
		dir := s.T().TempDir()
		delegate := companycard_mock.NewRepository(s.T())
		repo := fccard.NewRepository(fccard.ConfigRepository{
			Delegate: delegate,
			Dir:      dir,
			TTL:      ttl,
		})

		first := sberCard()
		second := sberCard()
		second.Description = "обновлённое описание"
		delegate.EXPECT().
			FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER").
			Return(first, nil).
			Once()
		delegate.EXPECT().
			FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER").
			Return(second, nil).
			Once()

		writtenAt := time.Now().UTC()
		got, err := repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER")
		s.Require().NoError(err)
		s.Equal(first, got)
		s.Equal(writtenAt.Add(ttl), readEnvelope(s.T(), dir, "1_SBER.json").ExpiresAt)

		// Сдвигаемся за пределы TTL — следующий вызов должен идти в делегат.
		time.Sleep(ttl + time.Second)
		refreshedAt := time.Now().UTC()
		got, err = repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER")
		s.Require().NoError(err)
		s.Equal(second, got)
		envelope := readEnvelope(s.T(), dir, "1_SBER.json")
		s.Equal(second, envelope.Card)
		s.Equal(refreshedAt.Add(ttl), envelope.ExpiresAt)
	})
}

// TestFindByTickerWithinTTLServesFromCache: пока now < ExpiresAt, кеш
// считается живым и делегат не зовётся.
func (s *repositorySuite) TestFindByTickerWithinTTLServesFromCache() {
	synctest.Test(s.T(), func(_ *testing.T) {
		dir := s.T().TempDir()
		delegate := companycard_mock.NewRepository(s.T())
		repo := fccard.NewRepository(fccard.ConfigRepository{
			Delegate: delegate,
			Dir:      dir,
			TTL:      time.Hour,
		})

		card := sberCard()
		delegate.EXPECT().
			FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER").
			Return(card, nil).
			Once()

		_, err := repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER")
		s.Require().NoError(err)

		// Доходим почти до конца TTL — запись должна продолжать жить.
		time.Sleep(time.Hour - time.Second)
		got, err := repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER")
		s.Require().NoError(err)
		s.Equal(card, got)
	})
}

// TestFindByTickerZeroTTLNeverExpires: TTL == 0 означает «без
// экспирации» — ExpiresAt в файле нулевой, и сколько бы часы ни шли,
// запись всегда считается живой.
func (s *repositorySuite) TestFindByTickerZeroTTLNeverExpires() {
	synctest.Test(s.T(), func(_ *testing.T) {
		dir := s.T().TempDir()
		delegate := companycard_mock.NewRepository(s.T())
		repo := fccard.NewRepository(fccard.ConfigRepository{
			Delegate: delegate,
			Dir:      dir,
		})

		card := sberCard()
		delegate.EXPECT().
			FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER").
			Return(card, nil).
			Once()

		_, err := repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER")
		s.Require().NoError(err)
		s.True(readEnvelope(s.T(), dir, "1_SBER.json").ExpiresAt.IsZero())

		time.Sleep(100 * 365 * 24 * time.Hour)
		got, err := repo.FindByTicker(context.Background(), domaincard.ExchangeMOEX, "SBER")
		s.Require().NoError(err)
		s.Equal(card, got)
	})
}

func sberCard() domaincard.Card {
	return domaincard.Card{
		Ticker:                "SBER",
		Exchange:              domaincard.ExchangeMOEX,
		Name:                  "Сбербанк",
		Sector:                "Финансы",
		Industry:              "Банковская деятельность",
		IndustryGroup:         "Банковская деятельность",
		Country:               "Россия",
		Currency:              domaincard.CurrencyRUB,
		PrimaryReportTicker:   "SBER",
		PrimaryReportExchange: domaincard.ExchangeMOEX,
		Description:           "ПАО «Сбербанк» — крупнейший универсальный банк России.",
		Site:                  "https://www.sberbank.com",
		DiscLink:              "https://www.sberbank.com/ru/investor-relations",
		SectorID:              40,
		IndustryID:            401010,
		IndustryGroupID:       4010,
	}
}

func readEnvelope(t *testing.T, dir, name string) envelopeOnDisk {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, name)) //nolint:gosec // dir = t.TempDir(), name — литерал теста
	if err != nil {
		t.Fatalf("read envelope file: %v", err)
	}
	var envelope envelopeOnDisk
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	return envelope
}
