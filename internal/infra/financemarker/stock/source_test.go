package stock_test

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
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/stock"
)

//go:embed testdata/*.json
var testdataFS embed.FS

type sourceSuite struct {
	suite.Suite

	handler func(http.ResponseWriter, *http.Request)
	server  *httptest.Server
	source  *stock.Source
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
	s.source = stock.New(c)
}

func (s *sourceSuite) TearDownTest() {
	s.server.Close()
}

// wantSberInfo — ожидаемая Info-секция для SBER. Используется в тестах
// info-only и в комбинированном — данные 1:1.
func wantSberInfo() company.StockInfo {
	return company.StockInfo{
		IssuerName:            "Сбербанк",
		Sector:                "Финансы",
		IndustryGroup:         "Банковская деятельность",
		Industry:              "Банковская деятельность",
		SubIndustry:           "Диверсифицированные банки",
		Country:               "Россия",
		Description:           "ПАО «Сбербанк» — крупнейший универсальный банк России.",
		Site:                  "https://www.sberbank.com",
		DisclosureLink:        "https://www.sberbank.com/ru/investor-relations",
		PrimaryReportTicker:   "SBER",
		SectorID:              40,
		IndustryGroupID:       4010,
		IndustryID:            401010,
		SubIndustryID:         40101010,
		PrimaryReportExchange: company.ExchangeMOEX,
		Exchange:              company.ExchangeMOEX,
		Currency:              company.CurrencyRUB,
		ReportFrequency:       company.ReportFrequencyQuarterly,
		SPB:                   false,
	}
}

// wantSberSummary — ожидаемая Summary-секция для SBER. Используется в
// тестах summary-only и в комбинированном — данные 1:1.
func wantSberSummary() company.StockSummary {
	return company.StockSummary{
		ChangedAt:          time.Date(2026, 5, 11, 3, 32, 6, 0, time.UTC),
		Capital:            97627.3,
		EPS:                78.8,
		PEG:                0.56,
		PeterLynchTarget:   96.77,
		GrahamTarget:       160.46,
		DividendFrequency:  1,
		DividendStrike:     4,
		DividendGrowth:     3,
		DividendIndex:      0.7,
		DividendYield12M:   10.65,
		DividendYield3Y:    10.55,
		DividendYield5Y:    7.5,
		DividendGapLast:    280,
		DividendGapAverage: 482,
		GrowthRevenue3Y:    10.59,
		GrowthRevenue5Y:    13.42,
		GrowthEarnings3Y:   5.62,
		GrowthEarnings5Y:   7.37,
		GrowthAssets3Y:     9.69,
		GrowthAssets5Y:     10.89,
		GrowthEquity3Y:     10.45,
		GrowthEquity5Y:     9.72,
		GrowthFCF5Y:        7.57,
		IdeaBuy:            9,
		IdeaHold:           3,
		IdeaSell:           0,
		IdeaTarget:         387.392,
		IdeaPotential:      20.9125,
		IdeaConsensus:      company.IdeaConsensusBuy,
		InsiderConsensus:   company.InsiderConsensusBuys,
	}
}

func (s *sourceSuite) TestFindByTickerInfoAndSummary() {
	body := s.readFixture("sber_all_sections.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/api/fm/v2/stocks/MOEX:SBER", r.URL.Path)
		s.Equal("test-token", r.URL.Query().Get("api_token"))
		s.Equal("info,summary", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		[]company.StockSection{company.StockSectionInfo, company.StockSectionSummary},
	)

	s.Require().NoError(err)
	// Комбинированный ответ должен 1:1 сложиться из тех же кусочков,
	// что приходят в info-only и summary-only ответах.
	s.Equal(wantSberInfo(), got.Info)
	s.Equal(wantSberSummary(), got.Summary)
}

// TestFindByTickerCanonicalIncludeOrder: даже если вызывающий передаст
// секции в произвольном порядке с дубликатами, include всегда собирается
// в каноническом порядке (info,summary) — это требование справочника для
// устойчивости кеш-ключа.
func (s *sourceSuite) TestFindByTickerCanonicalIncludeOrder() {
	body := s.readFixture("sber_all_sections.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("info,summary", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	_, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		[]company.StockSection{
			company.StockSectionSummary,
			company.StockSectionInfo,
			company.StockSectionInfo,
		},
	)
	s.Require().NoError(err)
}

func (s *sourceSuite) TestFindByTickerInfoOnly() {
	body := s.readFixture("sber_info.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("info", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		[]company.StockSection{company.StockSectionInfo},
	)

	s.Require().NoError(err)
	s.Equal(wantSberInfo(), got.Info)
	s.Zero(got.Summary)
}

func (s *sourceSuite) TestFindByTickerSummaryOnly() {
	body := s.readFixture("sber_summary.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("summary", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		[]company.StockSection{company.StockSectionSummary},
	)

	s.Require().NoError(err)
	s.Zero(got.Info)
	s.Equal(wantSberSummary(), got.Summary)
}

func (s *sourceSuite) TestFindByTickerEmptySections() {
	called := false
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	_, err := s.source.FindByTicker(context.Background(), "SBER", nil)

	s.Require().Error(err)
	s.Require().ErrorContains(err, "no sections requested")
	s.False(called, "HTTP не должен вызываться при пустом наборе секций")
}

func (s *sourceSuite) TestFindByTickerUnspecifiedOnly() {
	called := false
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	_, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		[]company.StockSection{company.StockSectionUnspecified},
	)

	s.Require().Error(err)
	s.Require().ErrorContains(err, "no sections requested")
	s.False(called)
}

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

			_, err := s.source.FindByTicker(
				context.Background(),
				"any",
				[]company.StockSection{company.StockSectionInfo, company.StockSectionSummary},
			)

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

	_, err := s.source.FindByTicker(
		ctx,
		"SBER",
		[]company.StockSection{company.StockSectionInfo, company.StockSectionSummary},
	)

	s.Require().Error(err)
	s.ErrorContains(err, "context canceled")
}

func (s *sourceSuite) TestFindByTickerInvalidChangedAt() {
	body := []byte(`{"summary":{"changed_at":"not-a-date","idea_consensus":"WAT"}}`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	got, err := s.source.FindByTicker(
		context.Background(),
		"any",
		[]company.StockSection{company.StockSectionSummary},
	)

	s.Require().NoError(err)
	s.True(got.Summary.ChangedAt.IsZero())
	s.Equal(company.IdeaConsensusUnspecified, got.Summary.IdeaConsensus)
}

func (s *sourceSuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}
