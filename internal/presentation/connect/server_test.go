package connect_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	connectrpc "connectrpc.com/connect"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
	company_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/company"
)

type serverSuite struct {
	suite.Suite

	identities      *company_mock.IdentityGateway
	classifications *company_mock.ClassificationGateway
	server          *httptest.Server
	client          companyv1connect.CompanyServiceClient
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(serverSuite))
}

func (s *serverSuite) SetupTest() {
	s.identities = company_mock.NewIdentityGateway(s.T())
	s.classifications = company_mock.NewClassificationGateway(s.T())
	srv := pconnect.NewServer(services.NewCompanyService(s.identities, s.classifications))

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
	s.identities.EXPECT().FindByTicker(mock.Anything, "SBER").Return(company.Identity{
		Ticker:       "SBER",
		ISIN:         "RU0009029540",
		Name:         "Сбербанк России ПАО ао",
		SecurityType: company.SecurityTypeCommonShare,
		ListingLevel: company.ListingLevelFirst,
	}, nil).Once()
	s.classifications.EXPECT().FindByTicker(mock.Anything, "SBER").Return(company.Classification{
		Exchange:            company.ExchangeMOEX,
		Currency:            company.CurrencyRUB,
		Sector:              "Финансы",
		Industry:            "Банковская деятельность",
		Country:             "Россия",
		PrimaryReportTicker: "SBER",
	}, nil).Once()

	resp, err := s.call("SBER")

	s.Require().NoError(err)
	got := resp.Msg.GetCompany()
	s.Require().NotNil(got)
	s.Equal("SBER", got.GetTicker())
	s.Equal("RU0009029540", got.GetIsin())
	s.Equal("Сбербанк России ПАО ао", got.GetName())
	s.Equal(companyv1.SecurityType_SECURITY_TYPE_COMMON_SHARE, got.GetSecurityType())
	s.Equal(companyv1.ListingLevel_LISTING_LEVEL_FIRST, got.GetListingLevel())
	s.Equal(companyv1.Exchange_EXCHANGE_MOEX, got.GetExchange())
	s.Equal(companyv1.Currency_CURRENCY_RUB, got.GetCurrency())
	s.Equal("Финансы", got.GetSector())
	s.Equal("Банковская деятельность", got.GetIndustry())
	s.Equal("Россия", got.GetCountry())
	s.Equal("SBER", got.GetPrimaryReportTicker())
}

func (s *serverSuite) TestGetCompanySecurityTypeMapping() {
	cases := []struct {
		name string
		in   company.SecurityType
		want companyv1.SecurityType
	}{
		{"common", company.SecurityTypeCommonShare, companyv1.SecurityType_SECURITY_TYPE_COMMON_SHARE},
		{"preferred", company.SecurityTypePreferredShare, companyv1.SecurityType_SECURITY_TYPE_PREFERRED_SHARE},
		{"depositary", company.SecurityTypeDepositaryReceipt, companyv1.SecurityType_SECURITY_TYPE_DEPOSITARY_RECEIPT},
		{"unspecified", company.SecurityTypeUnspecified, companyv1.SecurityType_SECURITY_TYPE_UNSPECIFIED},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			s.identities.EXPECT().FindByTicker(mock.Anything, "X").
				Return(company.Identity{Ticker: "X", SecurityType: c.in}, nil).Once()
			s.classifications.EXPECT().FindByTicker(mock.Anything, "X").
				Return(company.Classification{}, nil).Once()

			resp, err := s.call("X")

			s.Require().NoError(err)
			s.Equal(c.want, resp.Msg.GetCompany().GetSecurityType())
		})
	}
}

func (s *serverSuite) TestGetCompanyListingLevelMapping() {
	cases := []struct {
		name string
		in   company.ListingLevel
		want companyv1.ListingLevel
	}{
		{"first", company.ListingLevelFirst, companyv1.ListingLevel_LISTING_LEVEL_FIRST},
		{"second", company.ListingLevelSecond, companyv1.ListingLevel_LISTING_LEVEL_SECOND},
		{"third", company.ListingLevelThird, companyv1.ListingLevel_LISTING_LEVEL_THIRD},
		{"unspecified", company.ListingLevelUnspecified, companyv1.ListingLevel_LISTING_LEVEL_UNSPECIFIED},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			s.identities.EXPECT().FindByTicker(mock.Anything, "X").
				Return(company.Identity{Ticker: "X", ListingLevel: c.in}, nil).Once()
			s.classifications.EXPECT().FindByTicker(mock.Anything, "X").
				Return(company.Classification{}, nil).Once()

			resp, err := s.call("X")

			s.Require().NoError(err)
			s.Equal(c.want, resp.Msg.GetCompany().GetListingLevel())
		})
	}
}

func (s *serverSuite) TestGetCompanyExchangeMapping() {
	cases := []struct {
		name string
		in   company.Exchange
		want companyv1.Exchange
	}{
		{"moex", company.ExchangeMOEX, companyv1.Exchange_EXCHANGE_MOEX},
		{"unspecified", company.ExchangeUnspecified, companyv1.Exchange_EXCHANGE_UNSPECIFIED},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			s.identities.EXPECT().FindByTicker(mock.Anything, "X").
				Return(company.Identity{Ticker: "X"}, nil).Once()
			s.classifications.EXPECT().FindByTicker(mock.Anything, "X").
				Return(company.Classification{Exchange: c.in}, nil).Once()

			resp, err := s.call("X")

			s.Require().NoError(err)
			s.Equal(c.want, resp.Msg.GetCompany().GetExchange())
		})
	}
}

func (s *serverSuite) TestGetCompanyCurrencyMapping() {
	cases := []struct {
		name string
		in   company.Currency
		want companyv1.Currency
	}{
		{"rub", company.CurrencyRUB, companyv1.Currency_CURRENCY_RUB},
		{"usd", company.CurrencyUSD, companyv1.Currency_CURRENCY_USD},
		{"eur", company.CurrencyEUR, companyv1.Currency_CURRENCY_EUR},
		{"unspecified", company.CurrencyUnspecified, companyv1.Currency_CURRENCY_UNSPECIFIED},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			s.identities.EXPECT().FindByTicker(mock.Anything, "X").
				Return(company.Identity{Ticker: "X"}, nil).Once()
			s.classifications.EXPECT().FindByTicker(mock.Anything, "X").
				Return(company.Classification{Currency: c.in}, nil).Once()

			resp, err := s.call("X")

			s.Require().NoError(err)
			s.Equal(c.want, resp.Msg.GetCompany().GetCurrency())
		})
	}
}

func (s *serverSuite) TestGetCompanyNotFoundFromIdentity() {
	s.identities.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(company.Identity{}, company.ErrNotFound).Once()
	s.classifications.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(company.Classification{}, nil).Maybe()

	_, err := s.call("ZZZZ")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeNotFound, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyNotFoundFromClassification() {
	s.identities.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(company.Identity{Ticker: "ZZZZ"}, nil).Maybe()
	s.classifications.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(company.Classification{}, company.ErrNotFound).Once()

	_, err := s.call("ZZZZ")

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
	s.identities.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(company.Identity{}, errors.New("downstream boom")).Once()
	s.classifications.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(company.Classification{}, nil).Maybe()

	_, err := s.call("SBER")

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
