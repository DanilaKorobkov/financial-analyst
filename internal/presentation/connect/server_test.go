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
			StockInfo: company.StockInfo{
				IssuerName:          "Сбербанк",
				Sector:              "Финансы",
				Industry:            "Банковская деятельность",
				Country:             "Россия",
				PrimaryReportTicker: "SBER",
				Exchange:            company.ExchangeMOEX,
				Currency:            company.CurrencyRUB,
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

	info := resp.Msg.GetCompany().GetStockInfo()
	s.Require().NotNil(info)
	s.Equal("Сбербанк", info.GetIssuerName())
	s.Equal(companyv1.Exchange_MOEX, info.GetExchange())
	s.Equal(companyv1.Currency_RUB, info.GetCurrency())
}

func (s *serverSuite) TestGetCompanyEnumCodes() {
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(company.Company{
			SecurityDescription: company.SecurityDescription{
				SecurityType: company.SecurityTypePreferredShare,
				ListingLevel: company.ListingLevelSecond,
			},
			StockInfo: company.StockInfo{
				Exchange:        company.ExchangeMOEX,
				Currency:        company.CurrencyUSD,
				ReportFrequency: company.ReportFrequencyQuarterly,
			},
		}, nil).
		Once()

	resp, err := s.call("SBER")
	s.Require().NoError(err)

	desc := resp.Msg.GetCompany().GetSecurityDescription()
	s.Equal(companyv1.SecurityType_PREFERRED_SHARE, desc.GetSecurityType())
	s.Equal(companyv1.ListingLevel_SECOND, desc.GetListingLevel())

	info := resp.Msg.GetCompany().GetStockInfo()
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

func (s *serverSuite) call(ticker string) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	return s.client.GetCompany(
		context.Background(),
		connectrpc.NewRequest(&companyv1.GetCompanyRequest{Ticker: ticker}),
	)
}
