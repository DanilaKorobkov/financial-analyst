package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	company_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/company"
)

type companyServiceSuite struct {
	suite.Suite

	identities      *company_mock.IdentityGateway
	classifications *company_mock.ClassificationGateway
	service         *services.CompanyService
}

func TestCompanyServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(companyServiceSuite))
}

func (s *companyServiceSuite) SetupTest() {
	s.identities = company_mock.NewIdentityGateway(s.T())
	s.classifications = company_mock.NewClassificationGateway(s.T())
	s.service = services.NewCompanyService(s.identities, s.classifications)
}

func (s *companyServiceSuite) TestGetCompanyHappyPath() {
	identity := sberIdentity()
	classification := sberClassification()
	s.identities.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(identity, nil).
		Once()
	s.classifications.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(classification, nil).
		Once()

	got, err := s.service.GetCompany(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(company.Company{Identity: identity, Classification: classification}, got)
}

func (s *companyServiceSuite) TestGetCompanyPassesTickerAsIs() {
	// Источники регистронезависимы, сервис не нормализует ввод.
	s.identities.EXPECT().
		FindByTicker(mock.Anything, "sBeR").
		Return(company.Identity{Ticker: "sBeR"}, nil).
		Once()
	s.classifications.EXPECT().
		FindByTicker(mock.Anything, "sBeR").
		Return(company.Classification{}, nil).
		Once()

	_, err := s.service.GetCompany(context.Background(), "sBeR")

	s.Require().NoError(err)
}

func (s *companyServiceSuite) TestGetCompanyEmptyTicker() {
	_, err := s.service.GetCompany(context.Background(), "")

	s.Require().ErrorIs(err, services.ErrTickerEmpty)
}

func (s *companyServiceSuite) TestGetCompanyIdentityNotFoundPropagates() {
	s.identities.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(company.Identity{}, company.ErrNotFound).
		Once()
	s.classifications.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(company.Classification{}, nil).
		Maybe()

	_, err := s.service.GetCompany(context.Background(), "missing")

	s.Require().ErrorIs(err, company.ErrNotFound)
}

func (s *companyServiceSuite) TestGetCompanyClassificationNotFoundPropagates() {
	s.identities.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(company.Identity{Ticker: "missing"}, nil).
		Maybe()
	s.classifications.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(company.Classification{}, company.ErrNotFound).
		Once()

	_, err := s.service.GetCompany(context.Background(), "missing")

	s.Require().ErrorIs(err, company.ErrNotFound)
}

func (s *companyServiceSuite) TestGetCompanyArbitraryErrorPropagates() {
	sentinel := errors.New("boom")
	s.identities.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(company.Identity{}, sentinel).
		Once()
	s.classifications.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(company.Classification{}, nil).
		Maybe()

	_, err := s.service.GetCompany(context.Background(), "any")

	s.Require().ErrorIs(err, sentinel)
}

func sberIdentity() company.Identity {
	return company.Identity{
		Ticker:       "SBER",
		ISIN:         "RU0009029540",
		Name:         "Сбербанк России ПАО ао",
		SecurityType: company.SecurityTypeCommonShare,
		ListingLevel: company.ListingLevelFirst,
	}
}

func sberClassification() company.Classification {
	return company.Classification{
		Exchange:            company.ExchangeMOEX,
		Currency:            company.CurrencyRUB,
		Sector:              "Финансы",
		Industry:            "Банковская деятельность",
		Country:             "Россия",
		PrimaryReportTicker: "SBER",
	}
}
