package financemarker_test

import (
	"context"
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker"
)

//go:embed testdata/*.json
var testdataFS embed.FS

type repositorySuite struct {
	suite.Suite

	handler func(http.ResponseWriter, *http.Request)
	server  *httptest.Server
	repo    *financemarker.CompanyMetricsRepository
}

func TestRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) SetupTest() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handler(w, r)
	}))
	client := financemarker.NewClient(s.server.URL+"/api/fm/v2", "test-token", 5*time.Second)
	s.repo = financemarker.NewCompanyMetricsRepository(client)
}

func (s *repositorySuite) TearDownTest() {
	s.server.Close()
}

func (s *repositorySuite) TestFindByTickerHappyPath() {
	body := s.readFixture("sber_metrics.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/api/fm/v2/stocks/MOEX:SBER", r.URL.Path)
		s.Equal("test-token", r.URL.Query().Get("api_token"))
		s.Equal("info,summary", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	metrics, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	expected := entities.CompanyMetrics{
		Card: entities.CompanyCard{
			Ticker:                "SBER",
			Exchange:              "MOEX",
			Name:                  "Сбербанк",
			Sector:                "Финансы",
			SectorID:              40,
			Industry:              "Банковская деятельность",
			IndustryID:            401010,
			IndustryGroup:         "Банковская деятельность",
			IndustryGroupID:       4010,
			Country:               "Россия",
			Currency:              "RUB",
			PrimaryReportTicker:   "SBER",
			PrimaryReportExchange: "MOEX",
		},
		Description:        "ПАО «Сбербанк» — крупнейший универсальный банк России.",
		Site:               "https://www.sberbank.com",
		DiscLink:           "https://www.sberbank.com/ru/investor-relations",
		Capital:            97627.3,
		EPS:                78.8,
		PEG:                0.56,
		PeterLynchTarget:   96.77,
		GrahamTarget:       160.46,
		DividendFrequency:  1,
		DividendStrike:     4,
		DividendGrowth:     3,
		DividendIndex:      0.7,
		DividendYield12m:   10.65,
		DividendYield3y:    10.55,
		DividendYield5y:    7.5,
		DividendGapLast:    280,
		DividendGapAverage: 482,
		GrowthRevenue3y:    10.59,
		GrowthRevenue5y:    13.42,
		GrowthEarnings3y:   5.62,
		GrowthEarnings5y:   7.37,
		GrowthAssets3y:     9.69,
		GrowthAssets5y:     10.89,
		GrowthEquity3y:     10.45,
		GrowthEquity5y:     9.72,
		GrowthFCF5y:        7.57,
		IdeaBuy:            9,
		IdeaHold:           3,
		IdeaSell:           0,
		IdeaConsensus:      entities.IdeaConsensusBuy,
		IdeaTarget:         387.392,
		IdeaPotential:      20.9125,
		InsiderConsensus:   entities.InsiderConsensusBuys,
		ChangedAt:          time.Date(2026, 5, 11, 3, 32, 6, 0, time.UTC),
	}
	s.Equal(expected, metrics)
}

func (s *repositorySuite) TestFindByTickerNotFound() {
	body := s.readFixture("not_found.json")
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(body)
	}

	_, err := s.repo.FindByTicker(context.Background(), "ZZZZ")

	s.Require().ErrorIs(err, entities.ErrNotFound)
}

func (s *repositorySuite) TestFindByTickerUnauthorizedFromStatus() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}

	_, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().ErrorIs(err, entities.ErrUnauthorized)
}

func (s *repositorySuite) TestFindByTickerUnauthorizedFromTokenMessage() {
	body := s.readFixture("unauthorized.json")
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(body)
	}

	_, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().ErrorIs(err, entities.ErrUnauthorized)
}

func (s *repositorySuite) TestFindByTickerQuotaExceeded() {
	body := s.readFixture("quota_exceeded.json")
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write(body)
	}

	_, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().ErrorIs(err, entities.ErrQuotaExceeded)
}

func (s *repositorySuite) TestFindByTickerServerError() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "financemarker http status 500")
}

func (s *repositorySuite) TestFindByTickerInvalidJSON() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}

	_, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "decode financemarker payload")
}

func (s *repositorySuite) TestFindByTickerContextCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.repo.FindByTicker(ctx, "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "context canceled")
	s.ErrorContains(err, "financemarker request")
}

func (s *repositorySuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}
