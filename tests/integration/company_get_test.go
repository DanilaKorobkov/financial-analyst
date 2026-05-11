//go:build integration

package integration

import (
	"context"
	"net/http"

	connectrpc "connectrpc.com/connect"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
)

func (s *IntegrationSuite) get(ticker string) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	return s.client.GetCompany(
		context.Background(),
		connectrpc.NewRequest(&companyv1.GetCompanyRequest{Ticker: ticker}),
	)
}

func (s *IntegrationSuite) TestGetCompanyHappyPath() {
	s.setMoexFixture("SBER", "sber.json")

	resp, err := s.get("SBER")

	s.Require().NoError(err)
	company := resp.Msg.GetCompany()
	s.Require().NotNil(company)
	s.Equal("SBER", company.GetTicker())
	s.Equal("RU0009029540", company.GetIsin())
	s.Equal("stock_shares", company.GetGroup())
	s.Equal(int64(21586948000), company.GetIssueSize())
	s.Require().NotNil(company.ListingLevel)
	s.Equal(int32(1), company.GetListingLevel())
}

func (s *IntegrationSuite) TestGetCompanyNotFound() {
	s.setMoexFixture("ZZZZ", "not_found.json")

	_, err := s.get("ZZZZ")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeNotFound, connectErr.Code())
}

func (s *IntegrationSuite) TestGetCompanyInvalidArgument() {
	_, err := s.get("")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInvalidArgument, connectErr.Code())
}

func (s *IntegrationSuite) TestGetCompanyInternalOnUpstream5xx() {
	s.moexHandler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err := s.get("SBER")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInternal, connectErr.Code())
}
