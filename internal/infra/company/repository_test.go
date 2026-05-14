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

type repositorySuite struct {
	suite.Suite

	securityDescription *company_mock.SecurityDescriptionSource
	stockInfo           *company_mock.StockInfoSource
	stockSummary        *company_mock.StockSummarySource
	repo                *infracompany.Repository
}

func TestRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) SetupTest() {
	s.securityDescription = company_mock.NewSecurityDescriptionSource(s.T())
	s.stockInfo = company_mock.NewStockInfoSource(s.T())
	s.stockSummary = company_mock.NewStockSummarySource(s.T())
	s.repo = infracompany.NewRepository(infracompany.ConfigRepository{
		SecurityDescription: s.securityDescription,
		StockInfo:           s.stockInfo,
		StockSummary:        s.stockSummary,
	})
}

func (s *repositorySuite) TestFindByTickerHappyPath() {
	wantDesc := domaincompany.SecurityDescription{Ticker: "SBER", ISIN: "RU0009029540"}
	wantInfo := domaincompany.StockInfo{IssuerName: "Сбербанк", Country: "Россия"}
	wantSummary := domaincompany.StockSummary{Capital: 97627.3, EPS: 78.8, IdeaConsensus: domaincompany.IdeaConsensusBuy}
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(wantDesc, nil).
		Once()
	s.stockInfo.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(wantInfo, nil).
		Once()
	s.stockSummary.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(wantSummary, nil).
		Once()

	got, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(wantDesc, got.SecurityDescription)
	s.Equal(wantInfo, got.StockInfo)
	s.Equal(wantSummary, got.StockSummary)
}

func (s *repositorySuite) TestFindByTickerSecurityDescriptionMissing() {
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.SecurityDescription{}, domaincompany.ErrNotFound).
		Once()
	s.stockInfo.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.StockInfo{}, nil).
		Maybe()
	s.stockSummary.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.StockSummary{}, nil).
		Maybe()

	_, err := s.repo.FindByTicker(context.Background(), "missing")
	s.Require().ErrorIs(err, domaincompany.ErrNotFound)
}

func (s *repositorySuite) TestFindByTickerStockInfoMissing() {
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.SecurityDescription{}, nil).
		Maybe()
	s.stockInfo.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.StockInfo{}, domaincompany.ErrNotFound).
		Once()
	s.stockSummary.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.StockSummary{}, nil).
		Maybe()

	_, err := s.repo.FindByTicker(context.Background(), "missing")
	s.Require().ErrorIs(err, domaincompany.ErrNotFound)
}

func (s *repositorySuite) TestFindByTickerStockSummaryMissing() {
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.SecurityDescription{}, nil).
		Maybe()
	s.stockInfo.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.StockInfo{}, nil).
		Maybe()
	s.stockSummary.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(domaincompany.StockSummary{}, domaincompany.ErrNotFound).
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
	s.stockInfo.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(domaincompany.StockInfo{}, nil).
		Maybe()
	s.stockSummary.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(domaincompany.StockSummary{}, nil).
		Maybe()

	_, err := s.repo.FindByTicker(context.Background(), "any")
	s.Require().ErrorIs(err, boom)
}

func (s *repositorySuite) TestFindByTickerStockInfoError() {
	boom := errors.New("fm boom")
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(domaincompany.SecurityDescription{}, nil).
		Maybe()
	s.stockInfo.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(domaincompany.StockInfo{}, boom).
		Once()
	s.stockSummary.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(domaincompany.StockSummary{}, nil).
		Maybe()

	_, err := s.repo.FindByTicker(context.Background(), "any")
	s.Require().ErrorIs(err, boom)
}

func (s *repositorySuite) TestFindByTickerStockSummaryError() {
	boom := errors.New("fm summary boom")
	s.securityDescription.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(domaincompany.SecurityDescription{}, nil).
		Maybe()
	s.stockInfo.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(domaincompany.StockInfo{}, nil).
		Maybe()
	s.stockSummary.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(domaincompany.StockSummary{}, boom).
		Once()

	_, err := s.repo.FindByTicker(context.Background(), "any")
	s.Require().ErrorIs(err, boom)
}
