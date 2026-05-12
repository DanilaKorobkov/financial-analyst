package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	entities_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/entities"
)

type companyInfoSuite struct {
	suite.Suite

	companies *entities_mock.CompanyRepository
	service   *services.CompanyInfo
}

func TestCompanyInfoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(companyInfoSuite))
}

func (s *companyInfoSuite) SetupTest() {
	s.companies = entities_mock.NewCompanyRepository(s.T())
	s.service = services.NewCompanyInfo(s.companies)
}

func (s *companyInfoSuite) TestLookupHappyPath() {
	expected := entities.Company{Ticker: "SBER", Name: "Сбербанк"}
	s.companies.EXPECT().FindByTicker(context.Background(), "SBER").Return(expected, nil).Once()

	got, err := s.service.Lookup(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(expected, got)
}

func (s *companyInfoSuite) TestLookupPassesTickerAsIs() {
	// Источники регистронезависимы, сервис не нормализует ввод.
	expected := entities.Company{Ticker: "SBER"}
	s.companies.EXPECT().FindByTicker(context.Background(), "sBeR").Return(expected, nil).Once()

	_, err := s.service.Lookup(context.Background(), "sBeR")

	s.Require().NoError(err)
}

func (s *companyInfoSuite) TestLookupEmptyTicker() {
	_, err := s.service.Lookup(context.Background(), "")

	s.Require().ErrorIs(err, services.ErrTickerEmpty)
}

func (s *companyInfoSuite) TestLookupNotFoundPropagates() {
	s.companies.EXPECT().FindByTicker(context.Background(), "ZZZZ").Return(entities.Company{}, entities.ErrCompanyNotFound).Once()

	_, err := s.service.Lookup(context.Background(), "ZZZZ")

	s.Require().ErrorIs(err, entities.ErrCompanyNotFound)
}

func (s *companyInfoSuite) TestLookupArbitraryErrorPropagates() {
	sentinel := errors.New("boom")
	s.companies.EXPECT().FindByTicker(context.Background(), "SBER").Return(entities.Company{}, sentinel).Once()

	_, err := s.service.Lookup(context.Background(), "SBER")

	s.Require().ErrorIs(err, sentinel)
}
