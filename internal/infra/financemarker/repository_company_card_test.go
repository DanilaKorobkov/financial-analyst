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
	s.repo = financemarker.NewCompanyCardRepository(financemarker.ConfigCompanyCardRepository{
		BaseURL: s.server.URL + "/api/fm/v2",
		Token:   "test-token",
		Timeout: 5 * time.Second,
	})
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

	card, err := s.repo.FindByTicker(context.Background(), entities.ExchangeMOEX, "SBER")

	s.Require().NoError(err)
	expected := entities.CompanyCard{
		Ticker:                "SBER",
		Exchange:              entities.ExchangeMOEX,
		Name:                  "Сбербанк",
		Sector:                "Финансы",
		Industry:              "Банковская деятельность",
		IndustryGroup:         "Банковская деятельность",
		Country:               "Россия",
		Currency:              entities.CurrencyRUB,
		PrimaryReportTicker:   "SBER",
		PrimaryReportExchange: entities.ExchangeMOEX,
		Description:           "ПАО «Сбербанк» — крупнейший универсальный банк России.",
		Site:                  "https://www.sberbank.com",
		DiscLink:              "https://www.sberbank.com/ru/investor-relations",
		SectorID:              40,
		IndustryID:            401010,
		IndustryGroupID:       4010,
	}
	s.Equal(expected, card)
}

// TestFindByTickerErrorMapping проходит по таблице ответов FinanceMarker и
// проверяет, что mapHTTPError + декодер payload-а возвращают ожидаемую ошибку.
// Только 404 поднимается в domain как entities.ErrMissingCompany; остальные коды
// (401, 403, 400+token_not_found, 5xx) едут наверх как непомеченный internal
// сбой с указанием причины.
func (s *repositorySuite) TestFindByTickerErrorMapping() {
	cases := []struct {
		errIs       error // если задан — проверяем errors.Is.
		name        string
		errContains string // иначе — проверяем подстроку в Error().
		body        []byte
		status      int
	}{
		{
			name:   "not found mapped to domain ErrMissingCompany",
			status: http.StatusNotFound,
			body:   s.readFixture("not_found.json"),
			errIs:  entities.ErrMissingCompany,
		},
		{
			name:        "unauthorized by status reported as internal",
			status:      http.StatusUnauthorized,
			errContains: "financemarker unauthorized: http status 401",
		},
		{
			name:        "token_not_found message reported as internal",
			status:      http.StatusBadRequest,
			body:        s.readFixture("unauthorized.json"),
			errContains: "financemarker unauthorized: token_not_found",
		},
		{
			name:        "quota exceeded reported as internal",
			status:      http.StatusForbidden,
			body:        s.readFixture("quota_exceeded.json"),
			errContains: "financemarker quota exceeded: http status 403",
		},
		{
			name:        "server error reported with code",
			status:      http.StatusInternalServerError,
			errContains: "financemarker http status 500",
		},
		{
			name:        "invalid JSON payload reports decode error",
			status:      http.StatusOK,
			body:        []byte("not json"),
			errContains: "decode financemarker payload",
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			s.handler = func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(c.status)
				if c.body != nil {
					_, _ = w.Write(c.body)
				}
			}

			_, err := s.repo.FindByTicker(context.Background(), entities.ExchangeMOEX, "SBER")

			s.Require().Error(err)
			if c.errIs != nil {
				s.Require().ErrorIs(err, c.errIs)
				return
			}
			s.ErrorContains(err, c.errContains)
		})
	}
}

func (s *repositorySuite) TestFindByTickerExchangeValidation() {
	cases := []struct {
		name        string
		errContains string
		exchange    entities.Exchange
	}{
		{
			name:        "unspecified exchange rejected",
			exchange:    entities.ExchangeUnspecified,
			errContains: "financemarker: exchange is unspecified",
		},
		{
			name:        "unsupported exchange rejected",
			exchange:    entities.Exchange(99),
			errContains: "financemarker: unsupported exchange 99",
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			_, err := s.repo.FindByTicker(context.Background(), c.exchange, "SBER")
			s.Require().Error(err)
			s.ErrorContains(err, c.errContains)
		})
	}
}

func (s *repositorySuite) TestFindByTickerContextCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.repo.FindByTicker(ctx, entities.ExchangeMOEX, "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "financemarker request: Get \""+s.server.URL+"/api/fm/v2/stocks/MOEX:SBER?api_token=test-token&include=info\": context canceled")
}

func (s *repositorySuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}
