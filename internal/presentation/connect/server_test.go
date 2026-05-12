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
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
	entities_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/entities"
)

type serverSuite struct {
	suite.Suite

	companies *entities_mock.CompanyRepository
	server    *httptest.Server
	client    companyv1connect.CompanyServiceClient
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(serverSuite))
}

func (s *serverSuite) SetupTest() {
	s.companies = entities_mock.NewCompanyRepository(s.T())
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
	s.companies.EXPECT().FindByTicker(mock.Anything, "SBER").Return(entities.Company{
		Ticker:       "SBER",
		ISIN:         "RU0009029540",
		Name:         "Сбербанк",
		SecurityType: entities.SecurityTypeCommonShare,
		ListingLevel: entities.ListingLevelFirst,
	}, nil).Once()

	resp, err := s.call("SBER")

	s.Require().NoError(err)
	company := resp.Msg.GetCompany()
	s.Require().NotNil(company)
	s.Equal("SBER", company.GetTicker())
	s.Equal("RU0009029540", company.GetIsin())
	s.Equal("Сбербанк", company.GetName())
	s.Equal(companyv1.SecurityType_SECURITY_TYPE_COMMON_SHARE, company.GetSecurityType())
	s.Equal(companyv1.ListingLevel_LISTING_LEVEL_FIRST, company.GetListingLevel())
}

func (s *serverSuite) TestGetCompanyUnspecifiedListingLevel() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "X").Return(entities.Company{Ticker: "X"}, nil).Once()

	resp, err := s.call("X")

	s.Require().NoError(err)
	s.Equal(companyv1.ListingLevel_LISTING_LEVEL_UNSPECIFIED, resp.Msg.GetCompany().GetListingLevel())
}

func (s *serverSuite) TestGetCompanyNotFound() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "ZZZZ").Return(entities.Company{}, entities.ErrCompanyNotFound).Once()

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
	s.companies.EXPECT().FindByTicker(mock.Anything, "SBER").Return(entities.Company{}, errors.New("downstream boom")).Once()

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
