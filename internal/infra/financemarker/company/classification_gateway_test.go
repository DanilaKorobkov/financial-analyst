package company_test

import (
	"context"
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker"
	fmcompany "github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/company"
)

//go:embed testdata/*.json
var testdataFS embed.FS

type classificationGatewaySuite struct {
	suite.Suite

	handler func(http.ResponseWriter, *http.Request)
	server  *httptest.Server
	gateway *fmcompany.ClassificationGateway
}

func TestClassificationGatewaySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(classificationGatewaySuite))
}

func (s *classificationGatewaySuite) SetupTest() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handler(w, r)
	}))
	client := financemarker.NewClient(financemarker.ConfigClient{
		BaseURL: s.server.URL + "/api/fm/v2",
		Token:   "test-token",
		Timeout: 5 * time.Second,
	})
	s.gateway = fmcompany.NewClassificationGateway(client)
}

func (s *classificationGatewaySuite) TearDownTest() {
	s.server.Close()
}

func (s *classificationGatewaySuite) TestFindByTickerHappyPath() {
	body := s.readFixture("sber_card.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/api/fm/v2/stocks/MOEX:SBER", r.URL.Path)
		s.Equal("test-token", r.URL.Query().Get("api_token"))
		s.Equal("info", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	got, err := s.gateway.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	expected := domaincompany.Classification{
		Exchange:            domaincompany.ExchangeMOEX,
		Currency:            domaincompany.CurrencyRUB,
		Sector:              "Финансы",
		Industry:            "Банковская деятельность",
		Country:             "Россия",
		PrimaryReportTicker: "SBER",
	}
	s.Equal(expected, got)
}

// TestFindByTickerErrorMapping проходит по таблице ответов FinanceMarker и
// проверяет, что классификация HTTP-ошибок + декодер payload-а возвращают
// ожидаемую ошибку. Только 404 поднимается как domaincompany.ErrNotFound;
// остальные коды (401, 403, 400+token_not_found, 5xx) едут наверх как
// непомеченный internal сбой либо как infra-sentinel.
func (s *classificationGatewaySuite) TestFindByTickerErrorMapping() {
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
			errIs:  financemarker.ErrUnauthorized,
		},
		{
			name:   "token_not_found message reported as infra sentinel",
			status: http.StatusBadRequest,
			body:   s.readFixture("unauthorized.json"),
			errIs:  financemarker.ErrUnauthorized,
		},
		{
			name:   "quota exceeded reported as infra sentinel",
			status: http.StatusForbidden,
			body:   s.readFixture("quota_exceeded.json"),
			errIs:  financemarker.ErrQuotaExceeded,
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

			_, err := s.gateway.FindByTicker(context.Background(), "SBER")

			s.Require().Error(err)
			if c.errIs != nil {
				s.Require().ErrorIs(err, c.errIs)
				return
			}
			s.ErrorContains(err, c.errContains)
		})
	}
}

func (s *classificationGatewaySuite) TestFindByTickerContextCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.gateway.FindByTicker(ctx, "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "financemarker request: Get \""+s.server.URL+"/api/fm/v2/stocks/MOEX:SBER?api_token=test-token&include=info\": context canceled")
}

func (s *classificationGatewaySuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}
