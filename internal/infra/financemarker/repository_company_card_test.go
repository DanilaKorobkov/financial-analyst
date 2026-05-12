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
	repo    *financemarker.CompanyCardRepository
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
	s.repo = financemarker.NewCompanyCardRepository(client)
}

func (s *repositorySuite) TearDownTest() {
	s.server.Close()
}

func (s *repositorySuite) TestFindByTickerHappyPath() {
	body := s.readFixture("sber_card.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/api/fm/v2/stocks/MOEX:SBER", r.URL.Path)
		s.Equal("test-token", r.URL.Query().Get("api_token"))
		s.Equal("info", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	card, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	expected := entities.CompanyCard{
		Ticker:                "SBER",
		Exchange:              "MOEX",
		Name:                  "Сбербанк",
		Sector:                "Финансы",
		Industry:              "Банковская деятельность",
		IndustryGroup:         "Банковская деятельность",
		Country:               "Россия",
		Currency:              "RUB",
		PrimaryReportTicker:   "SBER",
		PrimaryReportExchange: "MOEX",
		Description:           "ПАО «Сбербанк» — крупнейший универсальный банк России.",
		Site:                  "https://www.sberbank.com",
		DiscLink:              "https://www.sberbank.com/ru/investor-relations",
		SectorID:              40,
		IndustryID:            401010,
		IndustryGroupID:       4010,
	}
	s.Equal(expected, card)
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
