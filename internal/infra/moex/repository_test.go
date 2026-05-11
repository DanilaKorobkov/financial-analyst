package moex_test

import (
	"context"
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
)

//go:embed testdata/*.json
var testdataFS embed.FS

type RepositorySuite struct {
	suite.Suite

	handler func(http.ResponseWriter, *http.Request)
	server  *httptest.Server
	repo    *moex.CompanyRepository
}

func (s *RepositorySuite) SetupTest() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handler(w, r)
	}))
	s.repo = moex.NewCompanyRepository(s.server.URL+"/iss", 5*time.Second)
}

func (s *RepositorySuite) TearDownTest() {
	s.server.Close()
}

func (s *RepositorySuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}

func (s *RepositorySuite) TestFindByTickerHappyPath() {
	body := s.readFixture("sber.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/iss/securities/SBER.json", r.URL.Path)
		s.Equal("extended", r.URL.Query().Get("iss.json"))
		s.Equal("off", r.URL.Query().Get("iss.meta"))
		s.Equal("description", r.URL.Query().Get("iss.only"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	company, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	expected := entities.Company{
		IssueDate:    time.Date(2007, 7, 20, 0, 0, 0, 0, time.UTC),
		Ticker:       "SBER",
		ISIN:         "RU0009029540",
		Name:         "Сбербанк России ПАО ао",
		ShortName:    "Сбербанк",
		RegNumber:    "10301481B",
		SecurityType: "common_share",
		Group:        "stock_shares",
		FaceUnit:     "SUR",
		IssueSize:    21586948000,
		EmitterID:    484,
		FaceValue:    3.0,
		ListingLevel: 1,
		Sessions: entities.Sessions{
			Morning: true,
			Evening: true,
			Weekend: true,
		},
	}
	s.Equal(expected, company)
}

func (s *RepositorySuite) TestFindByTickerInvalidJSON() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}

	_, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().Error(err)
	s.NotErrorIs(err, entities.ErrCompanyNotFound)
}

func (s *RepositorySuite) TestFindByTickerNotFound() {
	body := s.readFixture("not_found.json")
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.repo.FindByTicker(context.Background(), "ZZZZ")

	s.Require().ErrorIs(err, entities.ErrCompanyNotFound)
}

func (s *RepositorySuite) TestFindByTickerServerError() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().Error(err)
	s.NotErrorIs(err, entities.ErrCompanyNotFound)
}

func (s *RepositorySuite) TestFindByTickerContextCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.repo.FindByTicker(ctx, "SBER")

	s.Require().Error(err)
}

func TestRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositorySuite))
}
