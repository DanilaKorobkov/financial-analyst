package connect_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	connectrpc "connectrpc.com/connect"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
	company_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/aggregates/company"
)

type serverSuite struct {
	suite.Suite

	companies *company_mock.Repository
	server    *httptest.Server
	client    companyv1connect.CompanyServiceClient
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(serverSuite))
}

func (s *serverSuite) SetupTest() {
	s.companies = company_mock.NewRepository(s.T())
	srv := pconnect.NewServer(pconnect.ConfigServer{
		Companies: services.NewCompanyService(services.ConfigCompanyService{
			Companies: s.companies,
		}),
	})

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	s.server = httptest.NewServer(mux)
	s.client = companyv1connect.NewCompanyServiceClient(s.server.Client(), s.server.URL)
}

func (s *serverSuite) TearDownTest() {
	s.server.Close()
}

func (s *serverSuite) TestGetCompanyHappyPath() {
	issueDate := time.Date(2007, 7, 20, 0, 0, 0, 0, time.UTC)
	changedAt := time.Date(2026, 5, 11, 3, 32, 6, 0, time.UTC)
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(company.Company{
			SecurityDescription: company.SecurityDescription{
				Ticker:       "SBER",
				ISIN:         "RU0009029540",
				Name:         "Сбербанк России ПАО ао",
				SecurityType: company.SecurityTypeCommonShare,
				ListingLevel: company.ListingLevelFirst,
				IssueDate:    issueDate,
				IssueSize:    21586948000,
			},
			Stock: company.Stock{
				Info: company.StockInfo{
					IssuerName:          "Сбербанк",
					Sector:              "Финансы",
					Industry:            "Банковская деятельность",
					Country:             "Россия",
					PrimaryReportTicker: "SBER",
					Exchange:            company.ExchangeMOEX,
					Currency:            company.CurrencyRUB,
				},
				Summary: company.StockSummary{
					Capital:           97627.3,
					EPS:               78.8,
					PEG:               0.56,
					GrahamTarget:      160.46,
					DividendFrequency: 1,
					DividendYield12M:  10.65,
					GrowthRevenue3Y:   10.59,
					IdeaBuy:           9,
					IdeaTarget:        387.392,
					IdeaConsensus:     company.IdeaConsensusBuy,
					InsiderConsensus:  company.InsiderConsensusBuys,
					ChangedAt:         changedAt,
				},
			},
		}, nil).
		Once()

	resp, err := s.call("SBER")
	s.Require().NoError(err)

	desc := resp.Msg.GetCompany().GetSecurityDescription()
	s.Require().NotNil(desc)
	s.Equal("SBER", desc.GetTicker())
	s.Equal("RU0009029540", desc.GetIsin())
	s.Equal("Сбербанк России ПАО ао", desc.GetName())
	s.Equal(companyv1.SecurityType_COMMON_SHARE, desc.GetSecurityType())
	s.Equal(companyv1.ListingLevel_FIRST, desc.GetListingLevel())
	s.Equal(int64(21586948000), desc.GetIssueSize())
	s.Equal(issueDate.Unix(), desc.GetIssueDate().GetSeconds())

	stock := resp.Msg.GetCompany().GetStock()
	s.Require().NotNil(stock)

	info := stock.GetInfo()
	s.Require().NotNil(info)
	s.Equal("Сбербанк", info.GetIssuerName())
	s.Equal(companyv1.Exchange_MOEX, info.GetExchange())
	s.Equal(companyv1.Currency_RUB, info.GetCurrency())

	summary := stock.GetSummary()
	s.Require().NotNil(summary)
	s.InDelta(97627.3, summary.GetCapital(), 1e-6)
	s.InDelta(78.8, summary.GetEps(), 1e-6)
	s.InDelta(0.56, summary.GetPeg(), 1e-6)
	s.InDelta(160.46, summary.GetGrahamTarget(), 1e-6)
	s.Equal(int64(1), summary.GetDividendFrequency())
	s.InDelta(10.65, summary.GetDividendYield_12M(), 1e-6)
	s.InDelta(10.59, summary.GetGrowthRevenue_3Y(), 1e-6)
	s.Equal(int64(9), summary.GetIdeaBuy())
	s.InDelta(387.392, summary.GetIdeaTarget(), 1e-6)
	s.Equal(companyv1.IdeaConsensus_BUY, summary.GetIdeaConsensus())
	s.Equal(companyv1.InsiderConsensus_BUYS, summary.GetInsiderConsensus())
	s.Equal(changedAt.Unix(), summary.GetChangedAt().GetSeconds())
}

func (s *serverSuite) TestGetCompanyEnumCodes() {
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(company.Company{
			SecurityDescription: company.SecurityDescription{
				SecurityType: company.SecurityTypePreferredShare,
				ListingLevel: company.ListingLevelSecond,
			},
			Stock: company.Stock{
				Info: company.StockInfo{
					Exchange:        company.ExchangeMOEX,
					Currency:        company.CurrencyUSD,
					ReportFrequency: company.ReportFrequencyQuarterly,
				},
			},
		}, nil).
		Once()

	resp, err := s.call("SBER")
	s.Require().NoError(err)

	desc := resp.Msg.GetCompany().GetSecurityDescription()
	s.Equal(companyv1.SecurityType_PREFERRED_SHARE, desc.GetSecurityType())
	s.Equal(companyv1.ListingLevel_SECOND, desc.GetListingLevel())

	info := resp.Msg.GetCompany().GetStock().GetInfo()
	s.Equal(companyv1.Exchange_MOEX, info.GetExchange())
	s.Equal(companyv1.Currency_USD, info.GetCurrency())
	s.Equal(companyv1.ReportFrequency_QUARTERLY, info.GetReportFrequency())
}

// TestGetCompanyEnumUnspecifiedEncodedExplicitly: *Unspecified
// переводится в явный *_UNSPECIFIED proto-enum.
func (s *serverSuite) TestGetCompanyEnumUnspecifiedEncodedExplicitly() {
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(company.Company{
			SecurityDescription: company.SecurityDescription{
				SecurityType: company.SecurityTypeUnspecified,
			},
		}, nil).
		Once()

	resp, err := s.call("any")
	s.Require().NoError(err)

	desc := resp.Msg.GetCompany().GetSecurityDescription()
	s.Equal(companyv1.SecurityType_SECURITY_TYPE_UNSPECIFIED, desc.GetSecurityType())
}

// TestGetCompanyZeroDateIsNilTimestamp: нулевой time.Time переводится
// в nil-Timestamp (поле отсутствует в proto-сообщении).
func (s *serverSuite) TestGetCompanyZeroDateIsNilTimestamp() {
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(company.Company{}, nil).
		Once()

	resp, err := s.call("any")
	s.Require().NoError(err)

	desc := resp.Msg.GetCompany().GetSecurityDescription()
	s.Nil(desc.GetIssueDate())
	s.Nil(desc.GetRegistryDate())
}

func (s *serverSuite) TestGetCompanyNotFoundIsCodeNotFound() {
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(company.Company{}, company.ErrNotFound).
		Once()

	_, err := s.call("missing")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeNotFound, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyInvalidArgument() {
	_, err := s.call("")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInvalidArgument, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyInternal() {
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(company.Company{}, errors.New("downstream boom")).
		Once()

	_, err := s.call("any")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInternal, connectErr.Code())
}

// TestGetCompanyAllStockSections — happy-path по всем секциям карточки
// эмитента (Ratios/Reports/Dividends/Ideas/InsiderTransactions/Operations/
// Owners/Shares). Проверяет, что каждая секция доходит до proto-ответа
// и что enum-ы и timestamp-поля переведены корректно.
func (s *serverSuite) TestGetCompanyAllStockSections() {
	changedAt := time.Date(2026, 5, 15, 3, 32, 22, 0, time.UTC)
	lastBuyDate := time.Date(2026, 7, 17, 0, 0, 0, 0, time.UTC)
	reestrCloseDate := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	dateIn := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)
	transactionDate := time.Date(2025, 12, 11, 0, 0, 0, 0, time.UTC)

	period := company.StockPeriod{
		Year:      2026,
		Month:     3,
		Frequency: company.StockPeriodFrequencyYearToMonth,
		Standard:  company.ReportStandardIFRS,
	}
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(company.Company{
			Stock: company.Stock{
				Ratios: []company.StockRatio{
					{
						ChangedAt: changedAt,
						Period:    period,
						Capital:   99974.2,
						PE:        4.11,
						PBV:       0.83,
						Active:    true,
					},
				},
				Reports: []company.StockReport{
					{
						ChangedAt: changedAt,
						Period:    period,
						Currency:  company.CurrencyRUB,
						Amount:    1_000_000,
						Revenue:   4988300,
						Earnings:  1777700,
					},
				},
				Dividends: []company.StockDividend{
					{
						LastBuyDate:     lastBuyDate,
						ReestrCloseDate: reestrCloseDate,
						ChangedAt:       changedAt,
						DivAmount:       37.64,
						DivPercent:      11.6101,
						Year:            2025,
						Currency:        company.CurrencyRUB,
						Type:            company.DividendTypeYearly,
					},
				},
				Ideas: []company.StockIdea{
					{
						DateIn:    dateIn,
						ChangedAt: changedAt,
						Community: "РСХБ Инвестиции",
						ID:        6237,
						PriceIn:   319.0,
						PriceOut:  360.0,
						Status:    company.IdeaStatusActive,
					},
				},
				InsiderTransactions: []company.StockInsiderTransaction{
					{
						TransactionDate: transactionDate,
						Insider:         "Информация не раскрывается",
						Type:            company.InsiderTransactionTypePurchase,
					},
				},
				Operations: []company.StockOperation{
					{
						MetricID: "car_loans",
						Unit:     "₽",
						Period:   period,
						Amount:   1_000_000_000,
						Value:    170.4,
						Curs:     1.0,
					},
				},
				Owners: []company.StockOwner{
					{
						ChangedAt: changedAt,
						Owner:     "Прочие",
						Period:    period,
						Own:       50.0,
					},
				},
				Shares: []company.StockShare{
					{
						Ticker: "SBERP",
						Period: period,
						Num:    1_000_000_000,
					},
				},
			},
		}, nil).
		Once()

	resp, err := s.call("SBER")
	s.Require().NoError(err)

	stock := resp.Msg.GetCompany().GetStock()
	s.Require().NotNil(stock)

	s.Require().Len(stock.GetRatios(), 1)
	ratio := stock.GetRatios()[0]
	s.InDelta(99974.2, ratio.GetCapital(), 1e-6)
	s.True(ratio.GetActive())
	s.Equal(companyv1.StockPeriodFrequency_STOCK_PERIOD_FREQUENCY_YEAR_TO_MONTH, ratio.GetPeriod().GetFrequency())
	s.Equal(companyv1.ReportStandard_IFRS, ratio.GetPeriod().GetStandard())

	s.Require().Len(stock.GetReports(), 1)
	report := stock.GetReports()[0]
	s.InDelta(4988300, report.GetRevenue(), 1e-6)
	s.InDelta(1777700, report.GetEarnings(), 1e-6)
	s.Equal(int64(1_000_000), report.GetAmount())
	s.Equal(companyv1.Currency_RUB, report.GetCurrency())

	s.Require().Len(stock.GetDividends(), 1)
	div := stock.GetDividends()[0]
	s.InDelta(37.64, div.GetDivAmount(), 1e-6)
	s.InDelta(11.6101, div.GetDivPercent(), 1e-6)
	s.Equal(int64(2025), div.GetYear())
	s.Equal(companyv1.DividendType_DIVIDEND_TYPE_YEARLY, div.GetType())
	s.Equal(lastBuyDate.Unix(), div.GetLastBuyDate().GetSeconds())

	s.Require().Len(stock.GetIdeas(), 1)
	idea := stock.GetIdeas()[0]
	s.Equal(int64(6237), idea.GetId())
	s.Equal("РСХБ Инвестиции", idea.GetCommunity())
	s.Equal(companyv1.IdeaStatus_ACTIVE, idea.GetStatus())
	s.Equal(dateIn.Unix(), idea.GetDateIn().GetSeconds())

	s.Require().Len(stock.GetInsiderTransactions(), 1)
	tx := stock.GetInsiderTransactions()[0]
	s.Equal("Информация не раскрывается", tx.GetInsider())
	s.Equal(companyv1.InsiderTransactionType_PURCHASE, tx.GetType())
	s.Equal(transactionDate.Unix(), tx.GetTransactionDate().GetSeconds())

	s.Require().Len(stock.GetOperations(), 1)
	op := stock.GetOperations()[0]
	s.Equal("car_loans", op.GetMetricId())
	s.Equal("₽", op.GetUnit())
	s.InDelta(170.4, op.GetValue(), 1e-6)
	s.Equal(int64(1_000_000_000), op.GetAmount())

	s.Require().Len(stock.GetOwners(), 1)
	owner := stock.GetOwners()[0]
	s.Equal("Прочие", owner.GetOwner())
	s.InDelta(50.0, owner.GetOwn(), 1e-6)

	s.Require().Len(stock.GetShares(), 1)
	share := stock.GetShares()[0]
	s.Equal("SBERP", share.GetTicker())
	s.Equal(int64(1_000_000_000), share.GetNum())
}

// TestGetCompanyEmptySectionsAreNil — секции, не отданные источником,
// доезжают до proto-ответа как пустые/nil срезы.
func (s *serverSuite) TestGetCompanyEmptySectionsAreNil() {
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(company.Company{}, nil).
		Once()

	resp, err := s.call("any")
	s.Require().NoError(err)

	stock := resp.Msg.GetCompany().GetStock()
	s.Require().NotNil(stock)
	s.Empty(stock.GetRatios())
	s.Empty(stock.GetReports())
	s.Empty(stock.GetDividends())
	s.Empty(stock.GetIdeas())
	s.Empty(stock.GetInsiderTransactions())
	s.Empty(stock.GetOperations())
	s.Empty(stock.GetOwners())
	s.Empty(stock.GetShares())
}

func (s *serverSuite) call(ticker string) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	return s.client.GetCompany(
		context.Background(),
		connectrpc.NewRequest(&companyv1.GetCompanyRequest{Ticker: ticker}),
	)
}
