package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	company_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/company"
	data_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/data"
)

type companyServiceSuite struct {
	suite.Suite

	profiles            *company_mock.ProfileRepository
	securityDescription *data_mock.Bundle
	stockInfo           *data_mock.Bundle
	service             *services.CompanyService
}

func TestCompanyServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(companyServiceSuite))
}

func (s *companyServiceSuite) SetupTest() {
	s.profiles = company_mock.NewProfileRepository(s.T())
	s.securityDescription = data_mock.NewBundle(s.T())
	s.stockInfo = data_mock.NewBundle(s.T())

	s.securityDescription.EXPECT().BundleID().Return("security-description").Maybe()
	s.securityDescription.EXPECT().
		Fields().
		Return([]data.FieldDescriptor{
			{ID: company.FieldTicker},
			{ID: company.FieldISIN},
			{ID: company.FieldSecurityType},
			{ID: company.FieldListingLevel},
		}).
		Maybe()

	s.stockInfo.EXPECT().BundleID().Return("stock-info").Maybe()
	s.stockInfo.EXPECT().
		Fields().
		Return([]data.FieldDescriptor{
			{ID: company.FieldIssuerName},
			{ID: company.FieldCountry},
			{ID: company.FieldExchange},
			{ID: company.FieldCurrency},
		}).
		Maybe()

	registry := data.NewRegistry()
	s.Require().NoError(registry.Register("moex", s.securityDescription))
	s.Require().NoError(registry.Register("financemarker", s.stockInfo))
	s.service = services.NewCompanyService(services.ConfigCompanyService{
		Profiles: s.profiles,
		Registry: registry,
	})
}

func (s *companyServiceSuite) expectProfile(ticker string, fieldIDs ...data.Field) {
	s.profiles.EXPECT().
		FindByTicker(mock.Anything, ticker).
		Return(company.Profile{FieldIDs: fieldIDs}, nil).
		Once()
}

func (s *companyServiceSuite) TestGetCompanyHappyPath() {
	s.expectProfile(
		"SBER",
		company.FieldTicker,
		company.FieldISIN,
		company.FieldSecurityType,
		company.FieldListingLevel,
		company.FieldIssuerName,
		company.FieldCountry,
		company.FieldExchange,
		company.FieldCurrency,
	)
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "SBER").
		Return(data.FieldValues{
			company.FieldTicker:       "SBER",
			company.FieldISIN:         "RU0009029540",
			company.FieldSecurityType: company.SecurityTypeCommonShare,
			company.FieldListingLevel: company.ListingLevelFirst,
		}, nil).
		Once()
	s.stockInfo.EXPECT().
		Fetch(mock.Anything, "SBER").
		Return(data.FieldValues{
			company.FieldIssuerName: "Сбербанк",
			company.FieldCountry:    "Россия",
			company.FieldExchange:   company.ExchangeMOEX,
			company.FieldCurrency:   company.CurrencyRUB,
		}, nil).
		Once()

	got, err := s.service.GetCompany(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(data.FieldValues{
		company.FieldTicker:       "SBER",
		company.FieldISIN:         "RU0009029540",
		company.FieldSecurityType: company.SecurityTypeCommonShare,
		company.FieldListingLevel: company.ListingLevelFirst,
		company.FieldIssuerName:   "Сбербанк",
		company.FieldCountry:      "Россия",
		company.FieldExchange:     company.ExchangeMOEX,
		company.FieldCurrency:     company.CurrencyRUB,
	}, got)
}

func (s *companyServiceSuite) TestGetCompanyPassesTickerAsIs() {
	s.expectProfile("sBeR", company.FieldTicker)
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "sBeR").
		Return(data.FieldValues{company.FieldTicker: "sBeR"}, nil).
		Once()

	_, err := s.service.GetCompany(context.Background(), "sBeR")
	s.Require().NoError(err)
}

func (s *companyServiceSuite) TestGetCompanyEmptyTickerSkipsProfileLookup() {
	_, err := s.service.GetCompany(context.Background(), "")
	s.Require().ErrorIs(err, services.ErrTickerEmpty)
}

func (s *companyServiceSuite) TestGetCompanyProfileNotFound() {
	s.profiles.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(company.Profile{}, company.ErrProfileNotFound).
		Once()

	_, err := s.service.GetCompany(context.Background(), "missing")
	s.Require().ErrorIs(err, company.ErrProfileNotFound)
}

func (s *companyServiceSuite) TestGetCompanyProfileArbitraryError() {
	boom := errors.New("profile store down")
	s.profiles.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(company.Profile{}, boom).
		Once()

	_, err := s.service.GetCompany(context.Background(), "any")
	s.Require().ErrorIs(err, boom)
}

func (s *companyServiceSuite) TestGetCompanyUnknownFieldPropagates() {
	// Профиль ссылается на поле, которого нет в реестре — реестр сам
	// обнаружит и вернёт ErrFieldNotFound, сервис пробрасывает с пометкой.
	s.expectProfile("SBER", "unknown")

	_, err := s.service.GetCompany(context.Background(), "SBER")
	s.Require().ErrorIs(err, data.ErrFieldNotFound)
}

func (s *companyServiceSuite) TestGetCompanyBundleNotFoundPropagates() {
	s.expectProfile("missing", company.FieldTicker)
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "missing").
		Return(nil, company.ErrNotFound).
		Once()

	_, err := s.service.GetCompany(context.Background(), "missing")
	s.Require().ErrorIs(err, company.ErrNotFound)
}

func (s *companyServiceSuite) TestGetCompanyArbitraryBundleErrorPropagates() {
	sentinel := errors.New("boom")
	s.expectProfile("any", company.FieldTicker)
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "any").
		Return(nil, sentinel).
		Once()

	_, err := s.service.GetCompany(context.Background(), "any")
	s.Require().ErrorIs(err, sentinel)
}

func (s *companyServiceSuite) TestGetCompanyEmptyProfileSkipsFetch() {
	s.expectProfile("SBER")

	got, err := s.service.GetCompany(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Empty(got)
}
