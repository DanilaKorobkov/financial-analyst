package company_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
	infracompany "github.com/DanilaKorobkov/financial-analyst/internal/infra/company"
	company_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/aggregates/company"
)

// wantSections — набор секций, который Repository запрашивает у StockSource.
// Должен совпадать с stockSections в реализации.
var wantSections = []domaincompany.StockSection{
	domaincompany.StockSectionInfo,
	domaincompany.StockSectionSummary,
}

type repositorySuite struct {
	suite.Suite

	securityDescription *company_mock.SecurityDescriptionSource
	stock               *company_mock.StockSource
	repo                *infracompany.Repository
}

func TestRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) SetupTest() {
	s.securityDescription = company_mock.NewSecurityDescriptionSource(s.T())
	s.stock = company_mock.NewStockSource(s.T())
	s.repo = infracompany.NewRepository(infracompany.ConfigRepository{
		SecurityDescription: s.securityDescription,
		Stock:               s.stock,
	})
}

func (s *repositorySuite) TestFindByTickerHappyPath() {
	wantDesc := domaincompany.SecurityDescription{Ticker: "SBER", ISIN: "RU0009029540"}
	wantStock := domaincompany.Stock{
		Info:    domaincompany.StockInfo{IssuerName: "Сбербанк", Country: "Россия"},
		Summary: domaincompany.StockSummary{Capital: 97627.3, EPS: 78.8, IdeaConsensus: domaincompany.IdeaConsensusBuy},
	}
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(wantDesc, nil).
		Once()
	s.stock.EXPECT().
		FindByTicker(mock.Anything, "SBER", wantSections).
		Return(wantStock, nil).
		Once()

	got, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(wantDesc, got.SecurityDescription)
	s.Equal(wantStock, got.Stock)
}

func (s *repositorySuite) TestFindByTickerSecurityDescriptionMissing() {
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.SecurityDescription{}, domaincompany.ErrNotFound).
		Once()
	s.stock.EXPECT().
		FindByTicker(mock.Anything, "missing", wantSections).
		Return(domaincompany.Stock{}, nil).
		Maybe()

	_, err := s.repo.FindByTicker(context.Background(), "missing")
	s.Require().ErrorIs(err, domaincompany.ErrNotFound)
}

func (s *repositorySuite) TestFindByTickerStockMissing() {
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.SecurityDescription{}, nil).
		Maybe()
	s.stock.EXPECT().
		FindByTicker(mock.Anything, "missing", wantSections).
		Return(domaincompany.Stock{}, domaincompany.ErrNotFound).
		Once()

	_, err := s.repo.FindByTicker(context.Background(), "missing")
	s.Require().ErrorIs(err, domaincompany.ErrNotFound)
}

func (s *repositorySuite) TestFindByTickerSecurityDescriptionError() {
	boom := errors.New("moex boom")
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(domaincompany.SecurityDescription{}, boom).
		Once()
	s.stock.EXPECT().
		FindByTicker(mock.Anything, "any", wantSections).
		Return(domaincompany.Stock{}, nil).
		Maybe()

	_, err := s.repo.FindByTicker(context.Background(), "any")
	s.Require().ErrorIs(err, boom)
}

func (s *repositorySuite) TestFindByTickerStockError() {
	boom := errors.New("fm boom")
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(domaincompany.SecurityDescription{}, nil).
		Maybe()
	s.stock.EXPECT().
		FindByTicker(mock.Anything, "any", wantSections).
		Return(domaincompany.Stock{}, boom).
		Once()

	_, err := s.repo.FindByTicker(context.Background(), "any")
	s.Require().ErrorIs(err, boom)
}
