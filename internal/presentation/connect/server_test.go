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
	cards     *entities_mock.CompanyCardRepository
	server    *httptest.Server
	client    companyv1connect.CompanyServiceClient
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(serverSuite))
}

func (s *serverSuite) SetupTest() {
	s.companies = entities_mock.NewCompanyRepository(s.T())
	s.cards = entities_mock.NewCompanyCardRepository(s.T())
	srv := pconnect.NewServer(
		services.NewCompanyInfo(s.companies),
		services.NewCompanyCard(s.cards),
	)

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

	resp, err := s.callCompany("SBER")

	s.Require().NoError(err)
	company := resp.Msg.GetCompany()
	s.Require().NotNil(company)
	s.Equal("SBER", company.GetTicker())
	s.Equal("RU0009029540", company.GetIsin())
	s.Equal("Сбербанк", company.GetName())
	s.Equal(companyv1.SecurityType_SECURITY_TYPE_COMMON_SHARE, company.GetSecurityType())
	s.Equal(companyv1.ListingLevel_LISTING_LEVEL_FIRST, company.GetListingLevel())
}

func (s *serverSuite) TestGetCompanyNotFound() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(entities.Company{}, entities.ErrCompanyNotFound).Once()

	_, err := s.callCompany("ZZZZ")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeNotFound, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyInvalidArgument() {
	_, err := s.callCompany("")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInvalidArgument, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyInternal() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(entities.Company{}, errors.New("downstream boom")).Once()

	_, err := s.callCompany("SBER")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInternal, connectErr.Code())
}

func (s *serverSuite) TestGetCompanySecurityTypeMatrix() {
	cases := []struct {
		in   entities.SecurityType
		want companyv1.SecurityType
	}{
		{entities.SecurityTypeCommonShare, companyv1.SecurityType_SECURITY_TYPE_COMMON_SHARE},
		{entities.SecurityTypePreferredShare, companyv1.SecurityType_SECURITY_TYPE_PREFERRED_SHARE},
		{entities.SecurityTypeDepositaryReceipt, companyv1.SecurityType_SECURITY_TYPE_DEPOSITARY_RECEIPT},
		{entities.SecurityTypeUnspecified, companyv1.SecurityType_SECURITY_TYPE_UNSPECIFIED},
	}
	for _, c := range cases {
		s.Run(c.want.String(), func() {
			s.companies.EXPECT().FindByTicker(mock.Anything, "X").Return(entities.Company{
				Ticker:       "X",
				SecurityType: c.in,
				ListingLevel: entities.ListingLevelSecond,
			}, nil).Once()
			resp, err := s.callCompany("X")
			s.Require().NoError(err)
			s.Equal(c.want, resp.Msg.GetCompany().GetSecurityType())
			s.Equal(companyv1.ListingLevel_LISTING_LEVEL_SECOND, resp.Msg.GetCompany().GetListingLevel())
		})
	}
}

func (s *serverSuite) TestGetCompanyListingLevelThird() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "X").Return(entities.Company{
		Ticker:       "X",
		ListingLevel: entities.ListingLevelThird,
	}, nil).Once()
	resp, err := s.callCompany("X")
	s.Require().NoError(err)
	s.Equal(companyv1.ListingLevel_LISTING_LEVEL_THIRD, resp.Msg.GetCompany().GetListingLevel())
}

func (s *serverSuite) TestGetCompanyCardHappyPath() {
	s.cards.EXPECT().FindByTicker(mock.Anything, "SBER").Return(entities.CompanyCard{
		Ticker:                "SBER",
		Exchange:              "MOEX",
		Name:                  "Сбербанк",
		Sector:                "Финансы",
		SectorID:              40,
		Industry:              "Банковская деятельность",
		IndustryID:            401010,
		IndustryGroup:         "Банковская деятельность",
		IndustryGroupID:       4010,
		Country:               "Россия",
		Currency:              "RUB",
		PrimaryReportTicker:   "SBER",
		PrimaryReportExchange: "MOEX",
		Description:           "ПАО «Сбербанк»",
		Site:                  "https://www.sberbank.com",
		DiscLink:              "https://www.sberbank.com/ru/investor-relations",
	}, nil).Once()

	resp, err := s.callCard("SBER")

	s.Require().NoError(err)
	c := resp.Msg.GetCard()
	s.Require().NotNil(c)
	s.Equal("SBER", c.GetTicker())
	s.Equal("MOEX", c.GetExchange())
	s.Equal("Финансы", c.GetSector())
	s.Equal(int32(40), c.GetSectorId())
	s.Equal("Банковская деятельность", c.GetIndustry())
	s.Equal(int32(401010), c.GetIndustryId())
	s.Equal("Россия", c.GetCountry())
	s.Equal("RUB", c.GetCurrency())
	s.Equal("https://www.sberbank.com", c.GetSite())
}

func (s *serverSuite) TestGetCompanyCardInvalidArgument() {
	_, err := s.callCard("")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInvalidArgument, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyCardNotFound() {
	s.cards.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(entities.CompanyCard{}, entities.ErrNotFound).Once()

	_, err := s.callCard("ZZZZ")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeNotFound, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyCardUnauthenticated() {
	s.cards.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(entities.CompanyCard{}, entities.ErrUnauthorized).Once()

	_, err := s.callCard("SBER")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeUnauthenticated, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyCardResourceExhausted() {
	s.cards.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(entities.CompanyCard{}, entities.ErrQuotaExceeded).Once()

	_, err := s.callCard("SBER")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeResourceExhausted, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyCardInternal() {
	s.cards.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(entities.CompanyCard{}, errors.New("boom")).Once()

	_, err := s.callCard("SBER")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInternal, connectErr.Code())
}

func (s *serverSuite) callCompany(ticker string) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	return s.client.GetCompany(
		context.Background(),
		connectrpc.NewRequest(&companyv1.GetCompanyRequest{Ticker: ticker}),
	)
}

func (s *serverSuite) callCard(ticker string) (*connectrpc.Response[companyv1.GetCompanyCardResponse], error) {
	return s.client.GetCompanyCard(
		context.Background(),
		connectrpc.NewRequest(&companyv1.GetCompanyCardRequest{Ticker: ticker}),
	)
}
