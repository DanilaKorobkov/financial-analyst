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
	services_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/services"
)

type companyServiceSuite struct {
	suite.Suite

	profiles *company_mock.ProfileRepository
	fetcher  *services_mock.DataFetcher
	service  *services.CompanyService
}

func TestCompanyServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(companyServiceSuite))
}

func (s *companyServiceSuite) SetupTest() {
	s.profiles = company_mock.NewProfileRepository(s.T())
	s.fetcher = services_mock.NewDataFetcher(s.T())
	s.service = services.NewCompanyService(services.ConfigCompanyService{
		Profiles: s.profiles,
		Fetcher:  s.fetcher,
	})
}

func (s *companyServiceSuite) expectProfile(ticker string, fieldIDs ...data.Field) {
	s.profiles.EXPECT().
		FindByTicker(mock.Anything, ticker).
		Return(company.Profile{FieldIDs: fieldIDs}, nil).
		Once()
}

func (s *companyServiceSuite) TestGetCompanyHappyPath() {
	fields := []data.Field{
		company.FieldTicker,
		company.FieldISIN,
		company.FieldSecurityType,
		company.FieldListingLevel,
		company.FieldIssuerName,
		company.FieldCountry,
		company.FieldExchange,
		company.FieldCurrency,
	}
	want := data.FieldValues{
		company.FieldTicker:       "SBER",
		company.FieldISIN:         "RU0009029540",
		company.FieldSecurityType: company.SecurityTypeCommonShare,
		company.FieldListingLevel: company.ListingLevelFirst,
		company.FieldIssuerName:   "Сбербанк",
		company.FieldCountry:      "Россия",
		company.FieldExchange:     company.ExchangeMOEX,
		company.FieldCurrency:     company.CurrencyRUB,
	}
	s.expectProfile("SBER", fields...)
	s.fetcher.EXPECT().
		Fetch(mock.Anything, "SBER", fields).
		Return(want, nil).
		Once()

	got, err := s.service.GetCompany(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(want, got)
}

func (s *companyServiceSuite) TestGetCompanyPassesTickerAsIs() {
	s.expectProfile("sBeR", company.FieldTicker)
	s.fetcher.EXPECT().
		Fetch(mock.Anything, "sBeR", []data.Field{company.FieldTicker}).
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
	// Профиль ссылается на поле, которого нет у fetcher'а — он сам
	// обнаружит и вернёт ErrFieldNotFound, сервис пробрасывает с пометкой.
	s.expectProfile("SBER", "unknown")
	s.fetcher.EXPECT().
		Fetch(mock.Anything, "SBER", []data.Field{"unknown"}).
		Return(nil, data.ErrFieldNotFound).
		Once()

	_, err := s.service.GetCompany(context.Background(), "SBER")
	s.Require().ErrorIs(err, data.ErrFieldNotFound)
}

func (s *companyServiceSuite) TestGetCompanyFetcherNotFoundPropagates() {
	s.expectProfile("missing", company.FieldTicker)
	s.fetcher.EXPECT().
		Fetch(mock.Anything, "missing", []data.Field{company.FieldTicker}).
		Return(nil, company.ErrNotFound).
		Once()

	_, err := s.service.GetCompany(context.Background(), "missing")
	s.Require().ErrorIs(err, company.ErrNotFound)
}

func (s *companyServiceSuite) TestGetCompanyArbitraryFetcherErrorPropagates() {
	sentinel := errors.New("boom")
	s.expectProfile("any", company.FieldTicker)
	s.fetcher.EXPECT().
		Fetch(mock.Anything, "any", []data.Field{company.FieldTicker}).
		Return(nil, sentinel).
		Once()

	_, err := s.service.GetCompany(context.Background(), "any")
	s.Require().ErrorIs(err, sentinel)
}

func (s *companyServiceSuite) TestGetCompanyEmptyProfilePassesEmptyFields() {
	s.expectProfile("SBER")
	s.fetcher.EXPECT().
		Fetch(mock.Anything, "SBER", []data.Field(nil)).
		Return(data.FieldValues{}, nil).
		Once()

	got, err := s.service.GetCompany(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Empty(got)
}
