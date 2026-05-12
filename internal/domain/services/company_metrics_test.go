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

type companyMetricsSuite struct {
	suite.Suite

	metrics *entities_mock.CompanyMetricsRepository
	service *services.CompanyMetrics
}

func TestCompanyMetricsSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(companyMetricsSuite))
}

func (s *companyMetricsSuite) SetupTest() {
	s.metrics = entities_mock.NewCompanyMetricsRepository(s.T())
	s.service = services.NewCompanyMetrics(s.metrics)
}

func (s *companyMetricsSuite) TestFindByTickerHappyPath() {
	expected := entities.CompanyMetrics{
		Card: entities.CompanyCard{Ticker: "SBER", Name: "Сбербанк"},
		EPS:  78.8,
	}
	s.metrics.EXPECT().FindByTicker(context.Background(), "SBER").Return(expected, nil).Once()

	got, err := s.service.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(expected, got)
}

func (s *companyMetricsSuite) TestFindByTickerPassesTickerAsIs() {
	expected := entities.CompanyMetrics{Card: entities.CompanyCard{Ticker: "SBER"}}
	s.metrics.EXPECT().FindByTicker(context.Background(), "sBeR").Return(expected, nil).Once()

	_, err := s.service.FindByTicker(context.Background(), "sBeR")

	s.Require().NoError(err)
}

func (s *companyMetricsSuite) TestFindByTickerEmpty() {
	_, err := s.service.FindByTicker(context.Background(), "")

	s.Require().ErrorIs(err, services.ErrTickerEmpty)
}

func (s *companyMetricsSuite) TestFindByTickerNotFoundPropagates() {
	s.metrics.EXPECT().FindByTicker(context.Background(), "ZZZZ").
		Return(entities.CompanyMetrics{}, entities.ErrNotFound).Once()

	_, err := s.service.FindByTicker(context.Background(), "ZZZZ")

	s.Require().ErrorIs(err, entities.ErrNotFound)
}

func (s *companyMetricsSuite) TestFindByTickerUnauthorizedPropagates() {
	s.metrics.EXPECT().FindByTicker(context.Background(), "SBER").
		Return(entities.CompanyMetrics{}, entities.ErrUnauthorized).Once()

	_, err := s.service.FindByTicker(context.Background(), "SBER")

	s.Require().ErrorIs(err, entities.ErrUnauthorized)
}

func (s *companyMetricsSuite) TestFindByTickerQuotaExceededPropagates() {
	s.metrics.EXPECT().FindByTicker(context.Background(), "SBER").
		Return(entities.CompanyMetrics{}, entities.ErrQuotaExceeded).Once()

	_, err := s.service.FindByTicker(context.Background(), "SBER")

	s.Require().ErrorIs(err, entities.ErrQuotaExceeded)
}

func (s *companyMetricsSuite) TestFindByTickerArbitraryErrorPropagates() {
	sentinel := errors.New("boom")
	s.metrics.EXPECT().FindByTicker(context.Background(), "SBER").
		Return(entities.CompanyMetrics{}, sentinel).Once()

	_, err := s.service.FindByTicker(context.Background(), "SBER")

	s.Require().ErrorIs(err, sentinel)
}
