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

// wantSberInfo — ожидаемая Info-секция для SBER из фикстуры sber_all_sections.json.
func wantSberInfo() company.StockInfo {
	return company.StockInfo{
		IssuerName:            "Сбербанк",
		Sector:                "Финансы",
		IndustryGroup:         "Банковская деятельность",
		Industry:              "Банковская деятельность",
		SubIndustry:           "Диверсифицированные банки",
		Country:               "Россия",
		Description:           "ПАО «Сбербанк» — российский финансовый конгломерат,  крупнейший банк в России, Центральной и Восточной Европе, один из ведущих международных финансовых институтов.\n\n\nПредоставляет широкий спектр банковских услуг. В рамках стратегии трансформации «Сбербанка» в технологическую компанию начинает расти доля небанковских услуг, таких как онлайн-магазины электронной торговли, телекомы, страхование, медицина и прочее. \n\n\n",
		Site:                  "https://www.sberbank.com/ru",
		DisclosureLink:        "https://www.sberbank.com/ru/investor-relations/groupresults",
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

// wantSberSummary — ожидаемая Summary-секция для SBER из фикстуры
// sber_all_sections.json.
func wantSberSummary() company.StockSummary {
	return company.StockSummary{
		ChangedAt:          time.Date(2026, 5, 15, 3, 32, 22, 0, time.UTC),
		Capital:            99974.2,
		EPS:                78.8,
		PEG:                0.56,
		PeterLynchTarget:   94.45,
		GrahamTarget:       157.4,
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
		IdeaPotential:      19.4916,
		IdeaConsensus:      company.IdeaConsensusBuy,
		InsiderConsensus:   company.InsiderConsensusBuys,
	}
}

func (s *sourceSuite) serveAllSections(wantInclude string) {
	body := s.readFixture("sber_all_sections.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/api/fm/v2/stocks/MOEX:SBER", r.URL.Path)
		s.Equal("test-token", r.URL.Query().Get("api_token"))
		s.Equal(wantInclude, r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}

func (s *sourceSuite) TestFindByTickerInfoOnly() {
	s.serveAllSections("info")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithInfo: true},
	)

	s.Require().NoError(err)
	s.Equal(wantSberInfo(), got.Info)
	s.Zero(got.Summary)
	s.Nil(got.Ratios)
}

func (s *sourceSuite) TestFindByTickerSummaryOnly() {
	s.serveAllSections("summary")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithSummary: true},
	)

	s.Require().NoError(err)
	s.Zero(got.Info)
	s.Equal(wantSberSummary(), got.Summary)
}

func (s *sourceSuite) TestFindByTickerInfoAndSummary() {
	s.serveAllSections("info,summary")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithInfo: true, WithSummary: true},
	)

	s.Require().NoError(err)
	s.Equal(wantSberInfo(), got.Info)
	s.Equal(wantSberSummary(), got.Summary)
}

func (s *sourceSuite) TestFindByTickerRatios() {
	s.serveAllSections("ratios")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithRatios: true},
	)

	s.Require().NoError(err)
	s.Len(got.Ratios, 21)
	first := got.Ratios[0]
	s.Equal(2011, first.Period.Year)
	s.Equal(12, first.Period.Month)
	s.Equal(company.StockPeriodFrequencyYearly, first.Period.Frequency)
	s.Equal(company.ReportStandardIFRS, first.Period.Standard)
	s.InDelta(54822.8, first.Capital, 1e-6)
	s.InDelta(5.58, first.PE, 1e-6)
	s.False(first.Active)

	last := got.Ratios[20]
	s.True(last.Active)
	s.Equal(company.StockPeriodFrequencyYearToMonth, last.Period.Frequency)
}

func (s *sourceSuite) TestFindByTickerReports() {
	s.serveAllSections("reports")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithReports: true},
	)

	s.Require().NoError(err)
	s.Len(got.Reports, 92)
	first := got.Reports[0]
	s.Equal(2011, first.Period.Year)
	s.Equal(company.ReportStandardIFRS, first.Period.Standard)
	s.Equal(company.CurrencyRUB, first.Currency)
	s.Equal(int64(1000000), first.Amount)
	s.InDelta(741947, first.Revenue, 1e-6)
	s.InDelta(315900, first.Earnings, 1e-6)
}

func (s *sourceSuite) TestFindByTickerDividends() {
	s.serveAllSections("dividends")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithDividends: true},
	)

	s.Require().NoError(err)
	s.Len(got.Dividends, 18)
	first := got.Dividends[0]
	s.Equal(time.Date(2026, 7, 17, 0, 0, 0, 0, time.UTC), first.LastBuyDate)
	s.Equal(time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC), first.ReestrCloseDate)
	s.InDelta(37.64, first.DivAmount, 1e-6)
	s.InDelta(11.6101, first.DivPercent, 1e-6)
	s.Equal(int64(2025), first.Year)
	s.Equal(company.CurrencyRUB, first.Currency)
	s.Equal(company.DividendTypeYearly, first.Type)
	s.Zero(first.LastBuyPrice) // будущая выплата — цена ещё не зафиксирована
}

func (s *sourceSuite) TestFindByTickerIdeas() {
	s.serveAllSections("ideas")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithIdeas: true},
	)

	s.Require().NoError(err)
	s.Len(got.Ideas, 12)
	first := got.Ideas[0]
	s.Equal(int64(6237), first.ID)
	s.Equal("РСХБ Инвестиции", first.Community)
	s.Equal("Сбер открыл сезон охоты на прибыль", first.Idea)
	s.InDelta(319.0, first.PriceIn, 1e-6)
	s.InDelta(360.0, first.PriceOut, 1e-6)
	s.InDelta(13.0, first.ProfitPotential, 1e-6)
}

func (s *sourceSuite) TestFindByTickerInsiderTransactions() {
	s.serveAllSections("insiderTransactions")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithInsiderTransactions: true},
	)

	s.Require().NoError(err)
	s.Len(got.InsiderTransactions, 10)
	first := got.InsiderTransactions[0]
	s.Equal(time.Date(2025, 12, 11, 0, 0, 0, 0, time.UTC), first.TransactionDate)
	s.Equal("Информация не раскрывается", first.Insider)
	s.Equal(company.InsiderTransactionTypePurchase, first.Type)
}

func (s *sourceSuite) TestFindByTickerOperations() {
	s.serveAllSections("operations")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithOperations: true},
	)

	s.Require().NoError(err)
	s.Len(got.Operations, 267)
	first := got.Operations[0]
	s.Equal("car_loans", first.MetricID)
	s.InDelta(170.4, first.Value, 1e-6)
	s.Equal("₽", first.Unit)
	s.Equal(int64(1000000000), first.Amount)
	s.Equal(2014, first.Period.Year)
	s.Equal(company.StockPeriodFrequencyYearly, first.Period.Frequency)
	s.InDelta(1.0, first.Curs, 1e-6)
}

func (s *sourceSuite) TestFindByTickerOwners() {
	s.serveAllSections("owners")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithOwners: true},
	)

	s.Require().NoError(err)
	s.Len(got.Owners, 2)
	first := got.Owners[0]
	s.Equal("Прочие", first.Owner)
	s.InDelta(50.0, first.Own, 1e-6)
	s.Equal(2022, first.Period.Year)
	s.Equal(6, first.Period.Month)
}

func (s *sourceSuite) TestFindByTickerShares() {
	s.serveAllSections("shares")

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{WithShares: true},
	)

	s.Require().NoError(err)
	s.Len(got.Shares, 116)
	first := got.Shares[0]
	s.Equal("SBERP", first.Ticker)
	s.Equal(int64(1000000000), first.Num)
	s.Equal(2011, first.Period.Year)
}

func (s *sourceSuite) TestFindByTickerAllSectionsCanonicalInclude() {
	wantInclude := "info,summary,ratios,reports,dividends,ideas,insiderTransactions,operations,owners,shares"
	s.serveAllSections(wantInclude)

	got, err := s.source.FindByTicker(
		context.Background(),
		"SBER",
		company.StockOptions{
			WithInfo:                true,
			WithSummary:             true,
			WithRatios:              true,
			WithReports:             true,
			WithDividends:           true,
			WithIdeas:               true,
			WithInsiderTransactions: true,
			WithOperations:          true,
			WithOwners:              true,
			WithShares:              true,
		},
	)

	s.Require().NoError(err)
	s.Equal(wantSberInfo(), got.Info)
	s.Equal(wantSberSummary(), got.Summary)
	s.Len(got.Ratios, 21)
	s.Len(got.Reports, 92)
	s.Len(got.Dividends, 18)
	s.Len(got.Ideas, 12)
	s.Len(got.InsiderTransactions, 10)
	s.Len(got.Operations, 267)
	s.Len(got.Owners, 2)
	s.Len(got.Shares, 116)
}

func (s *sourceSuite) TestFindByTickerEmptyOptions() {
	called := false
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	_, err := s.source.FindByTicker(context.Background(), "SBER", company.StockOptions{})

	s.Require().Error(err)
	s.Require().ErrorContains(err, "no sections requested")
	s.False(called, "HTTP не должен вызываться при пустом наборе секций")
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
				company.StockOptions{WithInfo: true, WithSummary: true},
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
		company.StockOptions{WithInfo: true, WithSummary: true},
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
		company.StockOptions{WithSummary: true},
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
