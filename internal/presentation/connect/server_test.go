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

type serverSuite struct {
	suite.Suite

	companies      *entities_mock.CompanyRepository
	companyMetrics *entities_mock.CompanyMetricsRepository
	server         *httptest.Server
	client         companyv1connect.CompanyServiceClient
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(serverSuite))
}

func (s *serverSuite) SetupTest() {
	s.companies = entities_mock.NewCompanyRepository(s.T())
	s.companyMetrics = entities_mock.NewCompanyMetricsRepository(s.T())
	srv := pconnect.NewServer(
		services.NewCompanyInfo(s.companies),
		services.NewCompanyMetrics(s.companyMetrics),
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

func (s *serverSuite) TestGetCompanyUnspecifiedListingLevel() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "X").Return(entities.Company{Ticker: "X"}, nil).Once()

	resp, err := s.callCompany("X")

	s.Require().NoError(err)
	s.Equal(companyv1.ListingLevel_LISTING_LEVEL_UNSPECIFIED, resp.Msg.GetCompany().GetListingLevel())
}

func (s *serverSuite) TestGetCompanyNotFound() {
	s.companies.EXPECT().FindByTicker(mock.Anything, "ZZZZ").Return(entities.Company{}, entities.ErrCompanyNotFound).Once()

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
	s.companies.EXPECT().FindByTicker(mock.Anything, "SBER").Return(entities.Company{}, errors.New("downstream boom")).Once()

	_, err := s.callCompany("SBER")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInternal, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyMetricsHappyPath() {
	changedAt := time.Date(2026, 5, 11, 3, 32, 6, 0, time.UTC)
	s.companyMetrics.EXPECT().FindByTicker(mock.Anything, "SBER").Return(entities.CompanyMetrics{
		Card: entities.CompanyCard{
			Ticker:   "SBER",
			Exchange: "MOEX",
			Name:     "Сбербанк",
			Sector:   "Финансы",
			SectorID: 40,
			Currency: "RUB",
		},
		Description:      "ПАО «Сбербанк»",
		Site:             "https://www.sberbank.com",
		EPS:              78.8,
		PEG:              0.56,
		IdeaConsensus:    entities.IdeaConsensusBuy,
		InsiderConsensus: entities.InsiderConsensusBuys,
		ChangedAt:        changedAt,
	}, nil).Once()

	resp, err := s.callMetrics("SBER")

	s.Require().NoError(err)
	m := resp.Msg.GetMetrics()
	s.Require().NotNil(m)
	s.Equal("SBER", m.GetCard().GetTicker())
	s.Equal("Финансы", m.GetCard().GetSector())
	s.InDelta(78.8, m.GetEps(), 0.0001)
	s.InDelta(0.56, m.GetPeg(), 0.0001)
	s.Equal(companyv1.IdeaConsensus_IDEA_CONSENSUS_BUY, m.GetIdeaConsensus())
	s.Equal(companyv1.InsiderConsensus_INSIDER_CONSENSUS_BUYS, m.GetInsiderConsensus())
	s.Equal(changedAt.Unix(), m.GetChangedAt().AsTime().Unix())
}

func (s *serverSuite) TestGetCompanyMetricsUnspecifiedConsensuses() {
	s.companyMetrics.EXPECT().FindByTicker(mock.Anything, "X").Return(entities.CompanyMetrics{
		Card: entities.CompanyCard{Ticker: "X"},
	}, nil).Once()

	resp, err := s.callMetrics("X")

	s.Require().NoError(err)
	m := resp.Msg.GetMetrics()
	s.Equal(companyv1.IdeaConsensus_IDEA_CONSENSUS_UNSPECIFIED, m.GetIdeaConsensus())
	s.Equal(companyv1.InsiderConsensus_INSIDER_CONSENSUS_UNSPECIFIED, m.GetInsiderConsensus())
	s.Nil(m.GetChangedAt())
}

func (s *serverSuite) TestGetCompanyMetricsInvalidArgument() {
	_, err := s.callMetrics("")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInvalidArgument, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyMetricsNotFound() {
	s.companyMetrics.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(entities.CompanyMetrics{}, entities.ErrNotFound).Once()

	_, err := s.callMetrics("ZZZZ")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeNotFound, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyMetricsUnauthenticated() {
	s.companyMetrics.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(entities.CompanyMetrics{}, entities.ErrUnauthorized).Once()

	_, err := s.callMetrics("SBER")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeUnauthenticated, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyMetricsResourceExhausted() {
	s.companyMetrics.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(entities.CompanyMetrics{}, entities.ErrQuotaExceeded).Once()

	_, err := s.callMetrics("SBER")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeResourceExhausted, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyMetricsInternal() {
	s.companyMetrics.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(entities.CompanyMetrics{}, errors.New("boom")).Once()

	_, err := s.callMetrics("SBER")

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

func (s *serverSuite) TestGetCompanyMetricsConsensusMatrix() {
	type want struct {
		idea    companyv1.IdeaConsensus
		insider companyv1.InsiderConsensus
	}
	cases := []struct {
		idea     entities.IdeaConsensus
		insider  entities.InsiderConsensus
		expected want
	}{
		{entities.IdeaConsensusHold, entities.InsiderConsensusSells, want{
			companyv1.IdeaConsensus_IDEA_CONSENSUS_HOLD,
			companyv1.InsiderConsensus_INSIDER_CONSENSUS_SELLS,
		}},
		{entities.IdeaConsensusSell, entities.InsiderConsensusMixed, want{
			companyv1.IdeaConsensus_IDEA_CONSENSUS_SELL,
			companyv1.InsiderConsensus_INSIDER_CONSENSUS_MIXED,
		}},
	}
	for _, c := range cases {
		s.Run(c.expected.idea.String(), func() {
			s.companyMetrics.EXPECT().FindByTicker(mock.Anything, "X").Return(entities.CompanyMetrics{
				Card:             entities.CompanyCard{Ticker: "X"},
				IdeaConsensus:    c.idea,
				InsiderConsensus: c.insider,
			}, nil).Once()
			resp, err := s.callMetrics("X")
			s.Require().NoError(err)
			s.Equal(c.expected.idea, resp.Msg.GetMetrics().GetIdeaConsensus())
			s.Equal(c.expected.insider, resp.Msg.GetMetrics().GetInsiderConsensus())
		})
	}
}

func (s *serverSuite) callCompany(ticker string) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	return s.client.GetCompany(
		context.Background(),
		connectrpc.NewRequest(&companyv1.GetCompanyRequest{Ticker: ticker}),
	)
}

func (s *serverSuite) callMetrics(ticker string) (*connectrpc.Response[companyv1.GetCompanyMetricsResponse], error) {
	return s.client.GetCompanyMetrics(
		context.Background(),
		connectrpc.NewRequest(&companyv1.GetCompanyMetricsRequest{Ticker: ticker}),
	)
}
