package stockinfo_test

import (
	"context"
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/client"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/stockinfo"
)

//go:embed testdata/*.json
var testdataFS embed.FS

type sourceSuite struct {
	suite.Suite

	handler func(http.ResponseWriter, *http.Request)
	server  *httptest.Server
	source  *stockinfo.Source
}

func TestSourceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(sourceSuite))
}

func (s *sourceSuite) SetupTest() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handler(w, r)
	}))
	c := client.New(client.Config{
		BaseURL: s.server.URL + "/api/fm/v2",
		Token:   "test-token",
		Timeout: 5 * time.Second,
	})
	s.source = stockinfo.New(c)
}

func (s *sourceSuite) TearDownTest() {
	s.server.Close()
}

func (s *sourceSuite) TestFindByTickerHappyPath() {
	body := s.readFixture("sber_card.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/api/fm/v2/stocks/MOEX:SBER", r.URL.Path)
		s.Equal("test-token", r.URL.Query().Get("api_token"))
		s.Equal("info", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	got, err := s.source.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal("Сбербанк", got.IssuerName)
	s.Equal("Россия", got.Country)
	s.Equal("Финансы", got.Sector)
	s.Equal("Банковская деятельность", got.IndustryGroup)
	s.Equal("Банковская деятельность", got.Industry)
	s.Equal("Диверсифицированные банки", got.SubIndustry)
	s.Equal("ПАО «Сбербанк» — крупнейший универсальный банк России.", got.Description)
	s.Equal("https://www.sberbank.com", got.Site)
	s.Equal("https://www.sberbank.com/ru/investor-relations", got.DisclosureLink)
	s.Equal("SBER", got.PrimaryReportTicker)
	s.Equal(int64(40), got.SectorID)
	s.Equal(int64(4010), got.IndustryGroupID)
	s.Equal(int64(401010), got.IndustryID)
	s.Equal(int64(40101010), got.SubIndustryID)
	s.Equal(company.ExchangeMOEX, got.Exchange)
	s.Equal(company.ExchangeMOEX, got.PrimaryReportExchange)
	s.Equal(company.CurrencyRUB, got.Currency)
	s.Equal(company.ReportFrequencyQuarterly, got.ReportFrequency)
	s.False(got.SPB)
}

// TestFindByTickerErrorMapping проходит по таблице ответов FinanceMarker.
// Только 404 поднимается как company.ErrNotFound; остальные коды (401, 403,
// 400+token_not_found, 5xx) едут наверх как непомеченный internal сбой
// либо как infra-sentinel.
func (s *sourceSuite) TestFindByTickerErrorMapping() {
	cases := []struct {
		errIs       error
		name        string
		errContains string
		body        []byte
		status      int
	}{
		{
			name:   "not found mapped to domain ErrNotFound",
			status: http.StatusNotFound,
			body:   s.readFixture("not_found.json"),
			errIs:  company.ErrNotFound,
		},
		{
			name:   "unauthorized by status reported as infra sentinel",
			status: http.StatusUnauthorized,
			errIs:  client.ErrUnauthorized,
		},
		{
			name:   "token_not_found message reported as infra sentinel",
			status: http.StatusBadRequest,
			body:   s.readFixture("unauthorized.json"),
			errIs:  client.ErrUnauthorized,
		},
		{
			name:   "quota exceeded reported as infra sentinel",
			status: http.StatusForbidden,
			body:   s.readFixture("quota_exceeded.json"),
			errIs:  client.ErrQuotaExceeded,
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
				if c.body != nil {
					w.Header().Set("Content-Type", "application/json")
				}
				w.WriteHeader(c.status)
				if c.body != nil {
					_, _ = w.Write(c.body)
				}
			}

			_, err := s.source.FindByTicker(context.Background(), "any")

			s.Require().Error(err)
			if c.errIs != nil {
				s.Require().ErrorIs(err, c.errIs)
				return
			}
			s.ErrorContains(err, c.errContains)
		})
	}
}

func (s *sourceSuite) TestFindByTickerContextCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.source.FindByTicker(ctx, "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "financemarker request: Get \""+s.server.URL+"/api/fm/v2/stocks/MOEX:SBER?api_token=test-token&include=info\": context canceled")
}

func (s *sourceSuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}
