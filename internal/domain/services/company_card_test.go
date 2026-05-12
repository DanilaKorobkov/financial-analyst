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

type companyCardSuite struct {
	suite.Suite

	cards   *entities_mock.CompanyCardRepository
	service *services.CompanyCard
}

func TestCompanyCardSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(companyCardSuite))
}

func (s *companyCardSuite) SetupTest() {
	s.cards = entities_mock.NewCompanyCardRepository(s.T())
	s.service = services.NewCompanyCard(s.cards)
}

func (s *companyCardSuite) TestFindByTickerHappyPath() {
	expected := entities.CompanyCard{Ticker: "SBER", Name: "Сбербанк"}
	s.cards.EXPECT().FindByTicker(context.Background(), "SBER").Return(expected, nil).Once()

	got, err := s.service.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(expected, got)
}

func (s *companyCardSuite) TestFindByTickerPassesTickerAsIs() {
	expected := entities.CompanyCard{Ticker: "SBER"}
	s.cards.EXPECT().FindByTicker(context.Background(), "sBeR").Return(expected, nil).Once()

	_, err := s.service.FindByTicker(context.Background(), "sBeR")

	s.Require().NoError(err)
}

func (s *companyCardSuite) TestFindByTickerEmpty() {
	_, err := s.service.FindByTicker(context.Background(), "")

	s.Require().ErrorIs(err, services.ErrTickerEmpty)
}

func (s *companyCardSuite) TestFindByTickerNotFoundPropagates() {
	s.cards.EXPECT().FindByTicker(context.Background(), "ZZZZ").
		Return(entities.CompanyCard{}, entities.ErrNotFound).Once()

	_, err := s.service.FindByTicker(context.Background(), "ZZZZ")

	s.Require().ErrorIs(err, entities.ErrNotFound)
}

func (s *companyCardSuite) TestFindByTickerUnauthorizedPropagates() {
	s.cards.EXPECT().FindByTicker(context.Background(), "SBER").
		Return(entities.CompanyCard{}, entities.ErrUnauthorized).Once()

	_, err := s.service.FindByTicker(context.Background(), "SBER")

	s.Require().ErrorIs(err, entities.ErrUnauthorized)
}

func (s *companyCardSuite) TestFindByTickerQuotaExceededPropagates() {
	s.cards.EXPECT().FindByTicker(context.Background(), "SBER").
		Return(entities.CompanyCard{}, entities.ErrQuotaExceeded).Once()

	_, err := s.service.FindByTicker(context.Background(), "SBER")

	s.Require().ErrorIs(err, entities.ErrQuotaExceeded)
}

func (s *companyCardSuite) TestFindByTickerArbitraryErrorPropagates() {
	sentinel := errors.New("boom")
	s.cards.EXPECT().FindByTicker(context.Background(), "SBER").
		Return(entities.CompanyCard{}, sentinel).Once()

	_, err := s.service.FindByTicker(context.Background(), "SBER")

	s.Require().ErrorIs(err, sentinel)
}
