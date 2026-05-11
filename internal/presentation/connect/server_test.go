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
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
	entities_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/entities"
)

type ServerSuite struct {
	suite.Suite

	companies *entities_mock.CompanyRepository
	server    *httptest.Server
	client    companyv1connect.CompanyServiceClient
}

func (s *ServerSuite) SetupTest() {
	s.companies = entities_mock.NewCompanyRepository(s.T())
	srv := pconnect.NewServer(services.NewCompanyInfo(s.companies))

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	s.server = httptest.NewServer(mux)
	s.client = companyv1connect.NewCompanyServiceClient(s.server.Client(), s.server.URL)
}

func (s *ServerSuite) TearDownTest() {
	s.server.Close()
}

func (s *ServerSuite) call(ticker string) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	return s.client.GetCompany(
		context.Background(),
		connectrpc.NewRequest(&companyv1.GetCompanyRequest{Ticker: ticker}),
	)
}

func (s *ServerSuite) TestGetCompanyHappyPath() {
	issueDate := time.Date(2007, 7, 20, 0, 0, 0, 0, time.UTC)
	s.companies.EXPECT().FindByTicker(mock.Anything, "SBER").Return(entities.Company{
		Ticker:       "SBER",
		Name:         "Сбербанк",
		SecurityType: "common_share",
		IssueDate:    issueDate,
		ListingLevel: 1,
		Sessions:     entities.Sessions{Morning: true},
	}, nil).Once()

	resp, err := s.call("SBER")

	s.Require().NoError(err)
	company := resp.Msg.GetCompany()
	s.Require().NotNil(company)
	s.Equal("SBER", company.GetTicker())
	s.Equal("Сбербанк", company.GetName())
	s.Equal(companyv1.SecurityType_SECURITY_TYPE_COMMON_SHARE, company.GetSecurityType())
	s.Require().NotNil(company.ListingLevel)
	s.Equal(int32(1), company.GetListingLevel())
	s.Equal(issueDate.Unix(), company.GetIssueDate().AsTime().Unix())
	s.Require().NotNil(company.GetSessions())
	s.True(company.GetSessions().GetMorning())
}

func (s *ServerSuite) TestGetCompanyOmitsListingLevelWhenZero() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "X").Return(entities.Company{Ticker: "X"}, nil).Once()

	resp, err := s.call("X")

	s.Require().NoError(err)
	s.Nil(resp.Msg.GetCompany().ListingLevel)
	s.Nil(resp.Msg.GetCompany().IssueDate)
}

func (s *ServerSuite) TestGetCompanyNotFound() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "ZZZZ").Return(entities.Company{}, entities.ErrCompanyNotFound).Once()

	_, err := s.call("ZZZZ")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeNotFound, connectErr.Code())
}

func (s *ServerSuite) TestGetCompanyInvalidArgument() {
	_, err := s.call("")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInvalidArgument, connectErr.Code())
}

func (s *ServerSuite) TestGetCompanyInternal() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "SBER").Return(entities.Company{}, errors.New("downstream boom")).Once()

	_, err := s.call("SBER")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInternal, connectErr.Code())
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServerSuite))
}
