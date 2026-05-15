package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	company_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/aggregates/company"
)

type companyServiceSuite struct {
	suite.Suite

	companies *company_mock.Repository
	service   *services.CompanyService
}

func TestCompanyServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(companyServiceSuite))
}

func (s *companyServiceSuite) SetupTest() {
	s.companies = company_mock.NewRepository(s.T())
	s.service = services.NewCompanyService(services.ConfigCompanyService{
		Companies: s.companies,
	})
}

func (s *companyServiceSuite) TestGetCompanyHappyPath() {
	want := company.Company{
		SecurityDescription: company.SecurityDescription{
			Ticker:       "SBER",
			ISIN:         "RU0009029540",
			SecurityType: company.SecurityTypeCommonShare,
			ListingLevel: company.ListingLevelFirst,
			IssueDate:    time.Date(2007, 7, 20, 0, 0, 0, 0, time.UTC),
		},
		Stock: company.Stock{
			Info: company.StockInfo{
				IssuerName: "Сбербанк",
				Country:    "Россия",
				Exchange:   company.ExchangeMOEX,
				Currency:   company.CurrencyRUB,
			},
		},
	}
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(want, nil).
		Once()

	got, err := s.service.GetCompany(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(want, got)
}

func (s *companyServiceSuite) TestGetCompanyPassesTickerAsIs() {
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "sBeR").
		Return(company.Company{}, nil).
		Once()

	_, err := s.service.GetCompany(context.Background(), "sBeR")
	s.Require().NoError(err)
}

func (s *companyServiceSuite) TestGetCompanyEmptyTickerSkipsRepository() {
	_, err := s.service.GetCompany(context.Background(), "")
	s.Require().ErrorIs(err, services.ErrTickerEmpty)
}

func (s *companyServiceSuite) TestGetCompanyNotFoundPropagates() {
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(company.Company{}, company.ErrNotFound).
		Once()

	_, err := s.service.GetCompany(context.Background(), "missing")
	s.Require().ErrorIs(err, company.ErrNotFound)
}

func (s *companyServiceSuite) TestGetCompanyArbitraryErrorPropagates() {
	boom := errors.New("repository down")
	s.companies.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(company.Company{}, boom).
		Once()

	_, err := s.service.GetCompany(context.Background(), "any")
	s.Require().ErrorIs(err, boom)
}
