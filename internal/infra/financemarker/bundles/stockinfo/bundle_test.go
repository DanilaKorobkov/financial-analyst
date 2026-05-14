package stockinfo_test

import (
	"context"
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/bundles/stockinfo"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/client"
)

//go:embed testdata/*.json
var testdataFS embed.FS

type bundleSuite struct {
	suite.Suite

	handler func(http.ResponseWriter, *http.Request)
	server  *httptest.Server
	bundle  *stockinfo.Bundle
}

func TestBundleSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(bundleSuite))
}

func (s *bundleSuite) SetupTest() {
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
	s.bundle = stockinfo.New(c)
}

func (s *bundleSuite) TearDownTest() {
	s.server.Close()
}

func (s *bundleSuite) TestMetadata() {
	s.Equal(stockinfo.ID, s.bundle.BundleID())
	s.Require().Len(s.bundle.Fields(), 19)
}

func (s *bundleSuite) TestFetchHappyPath() {
	body := s.readFixture("sber_card.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/api/fm/v2/stocks/MOEX:SBER", r.URL.Path)
		s.Equal("test-token", r.URL.Query().Get("api_token"))
		s.Equal("info", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	values, err := s.bundle.Fetch(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal("Сбербанк", values[domaincompany.FieldIssuerName])
	s.Equal("Россия", values[domaincompany.FieldCountry])
	s.Equal("Финансы", values[domaincompany.FieldSector])
	s.Equal("Банковская деятельность", values[domaincompany.FieldIndustryGroup])
	s.Equal("Банковская деятельность", values[domaincompany.FieldIndustry])
	s.Equal("Диверсифицированные банки", values[domaincompany.FieldSubIndustry])
	s.Equal("ПАО «Сбербанк» — крупнейший универсальный банк России.", values[domaincompany.FieldDescription])
	s.Equal("https://www.sberbank.com", values[domaincompany.FieldSite])
	s.Equal("https://www.sberbank.com/ru/investor-relations", values[domaincompany.FieldDisclosureLink])
	s.Equal("SBER", values[domaincompany.FieldPrimaryReportTicker])
	s.Equal(int64(40), values[domaincompany.FieldSectorID])
	s.Equal(int64(4010), values[domaincompany.FieldIndustryGroupID])
	s.Equal(int64(401010), values[domaincompany.FieldIndustryID])
	s.Equal(int64(40101010), values[domaincompany.FieldSubIndustryID])
	s.Equal(domaincompany.ExchangeMOEX, values[domaincompany.FieldExchange])
	s.Equal(domaincompany.ExchangeMOEX, values[domaincompany.FieldPrimaryReportExchange])
	s.Equal(domaincompany.CurrencyRUB, values[domaincompany.FieldCurrency])
	s.Equal(domaincompany.ReportFrequencyQuarterly, values[domaincompany.FieldReportFrequency])
	s.Equal(false, values[domaincompany.FieldSPB])
}

// TestFetchErrorMapping проходит по таблице ответов FinanceMarker и
// проверяет, что классификация HTTP-ошибок + декодер payload-а возвращают
// ожидаемую ошибку. Только 404 поднимается как domaincompany.ErrNotFound;
// остальные коды (401, 403, 400+token_not_found, 5xx) едут наверх как
// непомеченный internal сбой либо как infra-sentinel.
func (s *bundleSuite) TestFetchErrorMapping() {
	cases := []struct {
		errIs       error // если задан — проверяем errors.Is.
		name        string
		errContains string // иначе — проверяем подстроку в Error().
		body        []byte
		status      int
	}{
		{
			name:   "not found mapped to domain ErrNotFound",
			status: http.StatusNotFound,
			body:   s.readFixture("not_found.json"),
			errIs:  domaincompany.ErrNotFound,
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

			_, err := s.bundle.Fetch(context.Background(), "any")

			s.Require().Error(err)
			if c.errIs != nil {
				s.Require().ErrorIs(err, c.errIs)
				return
			}
			s.ErrorContains(err, c.errContains)
		})
	}
}

func (s *bundleSuite) TestFetchContextCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.bundle.Fetch(ctx, "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "financemarker request: Get \""+s.server.URL+"/api/fm/v2/stocks/MOEX:SBER?api_token=test-token&include=info\": context canceled")
}

func (s *bundleSuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}
