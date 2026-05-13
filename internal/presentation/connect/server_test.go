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
	srv := pconnect.NewServer(services.NewCompanyInfo(s.companies))

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
	s.companies.EXPECT().FindByTicker(mock.Anything, "SBER").Return(company.Company{
		Ticker:       "SBER",
		ISIN:         "RU0009029540",
		Name:         "Сбербанк",
		SecurityType: company.SecurityTypeCommonShare,
		ListingLevel: company.ListingLevelFirst,
	}, nil).Once()

	resp, err := s.call("SBER")

	s.Require().NoError(err)
	got := resp.Msg.GetCompany()
	s.Require().NotNil(got)
	s.Equal("SBER", got.GetTicker())
	s.Equal("RU0009029540", got.GetIsin())
	s.Equal("Сбербанк", got.GetName())
	s.Equal(companyv1.SecurityType_SECURITY_TYPE_COMMON_SHARE, got.GetSecurityType())
	s.Equal(companyv1.ListingLevel_LISTING_LEVEL_FIRST, got.GetListingLevel())
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
			s.companies.EXPECT().FindByTicker(mock.Anything, "X").
				Return(company.Company{Ticker: "X", SecurityType: c.in}, nil).Once()

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
			s.companies.EXPECT().FindByTicker(mock.Anything, "X").
				Return(company.Company{Ticker: "X", ListingLevel: c.in}, nil).Once()

			resp, err := s.call("X")

			s.Require().NoError(err)
			s.Equal(c.want, resp.Msg.GetCompany().GetListingLevel())
		})
	}
}

func (s *serverSuite) TestGetCompanyUnspecifiedListingLevel() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "X").Return(company.Company{Ticker: "X"}, nil).Once()

	resp, err := s.call("X")

	s.Require().NoError(err)
	s.Equal(companyv1.ListingLevel_LISTING_LEVEL_UNSPECIFIED, resp.Msg.GetCompany().GetListingLevel())
}

func (s *serverSuite) TestGetCompanyNotFound() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "ZZZZ").Return(company.Company{}, company.ErrNotFound).Once()

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
	s.companies.EXPECT().FindByTicker(mock.Anything, "SBER").Return(company.Company{}, errors.New("downstream boom")).Once()

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
