package stocksummary_test

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
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/stocksummary"
)

//go:embed testdata/*.json
var testdataFS embed.FS

type sourceSuite struct {
	suite.Suite

	handler func(http.ResponseWriter, *http.Request)
	server  *httptest.Server
	source  *stocksummary.Source
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
	s.source = stocksummary.New(c)
}

func (s *sourceSuite) TearDownTest() {
	s.server.Close()
}

func (s *sourceSuite) TestFindByTickerHappyPath() {
	body := s.readFixture("sber_summary.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/api/fm/v2/stocks/MOEX:SBER", r.URL.Path)
		s.Equal("test-token", r.URL.Query().Get("api_token"))
		s.Equal("summary", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	got, err := s.source.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	s.InDelta(97627.3, got.Capital, 1e-6)
	s.InDelta(78.8, got.EPS, 1e-6)
	s.InDelta(0.56, got.PEG, 1e-6)
	s.InDelta(96.77, got.PeterLynchTarget, 1e-6)
	s.InDelta(160.46, got.GrahamTarget, 1e-6)
	s.Equal(1, got.DividendFrequency)
	s.Equal(4, got.DividendStrike)
	s.Equal(3, got.DividendGrowth)
	s.InDelta(0.7, got.DividendIndex, 1e-6)
	s.InDelta(10.65, got.DividendYield12M, 1e-6)
	s.InDelta(10.55, got.DividendYield3Y, 1e-6)
	s.InDelta(7.5, got.DividendYield5Y, 1e-6)
	s.Equal(280, got.DividendGapLast)
	s.Equal(482, got.DividendGapAverage)
	s.InDelta(10.59, got.GrowthRevenue3Y, 1e-6)
	s.InDelta(13.42, got.GrowthRevenue5Y, 1e-6)
	s.InDelta(5.62, got.GrowthEarnings3Y, 1e-6)
	s.InDelta(7.37, got.GrowthEarnings5Y, 1e-6)
	s.InDelta(9.69, got.GrowthAssets3Y, 1e-6)
	s.InDelta(10.89, got.GrowthAssets5Y, 1e-6)
	s.InDelta(10.45, got.GrowthEquity3Y, 1e-6)
	s.InDelta(9.72, got.GrowthEquity5Y, 1e-6)
	s.InDelta(7.57, got.GrowthFCF5Y, 1e-6)
	s.Zero(got.GrowthEBITDA3Y)
	s.Zero(got.GrowthNetDebt5Y)
	s.Equal(9, got.IdeaBuy)
	s.Equal(3, got.IdeaHold)
	s.Equal(0, got.IdeaSell)
	s.InDelta(387.392, got.IdeaTarget, 1e-6)
	s.InDelta(20.9125, got.IdeaPotential, 1e-6)
	s.Equal(company.IdeaConsensusBuy, got.IdeaConsensus)
	s.Equal(company.InsiderConsensusBuys, got.InsiderConsensus)
	s.Equal(time.Date(2026, 5, 11, 3, 32, 6, 0, time.UTC), got.ChangedAt)
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
	s.ErrorContains(err, "financemarker request: Get \""+s.server.URL+"/api/fm/v2/stocks/MOEX:SBER?api_token=test-token&include=summary\": context canceled")
}

func (s *sourceSuite) TestFindByTickerInvalidChangedAt() {
	body := []byte(`{"summary":{"changed_at":"not-a-date","idea_consensus":"WAT"}}`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	got, err := s.source.FindByTicker(context.Background(), "any")

	s.Require().NoError(err)
	s.True(got.ChangedAt.IsZero())
	s.Equal(company.IdeaConsensusUnspecified, got.IdeaConsensus)
}

func (s *sourceSuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}
