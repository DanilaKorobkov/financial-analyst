package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/companycard"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	companycard_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/companycard"
)

type companyCardServiceSuite struct {
	suite.Suite

	identities      *companycard_mock.IdentityGateway
	classifications *companycard_mock.ClassificationGateway
	service         *services.CompanyCardService
}

func TestCompanyCardServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(companyCardServiceSuite))
}

func (s *companyCardServiceSuite) SetupTest() {
	s.identities = companycard_mock.NewIdentityGateway(s.T())
	s.classifications = companycard_mock.NewClassificationGateway(s.T())
	s.service = services.NewCompanyCardService(s.identities, s.classifications)
}

func (s *companyCardServiceSuite) TestGetCardHappyPath() {
	identity := companycard.Identity{
		Ticker:       "SBER",
		ISIN:         "RU0009029540",
		Name:         "Сбербанк России ПАО ао",
		SecurityType: companycard.SecurityTypeCommonShare,
		ListingLevel: companycard.ListingLevelFirst,
	}
	classification := companycard.Classification{
		Exchange:            companycard.ExchangeMOEX,
		Currency:            companycard.CurrencyRUB,
		Sector:              "Финансы",
		Industry:            "Банковская деятельность",
		Country:             "Россия",
		PrimaryReportTicker: "SBER",
	}
	s.identities.EXPECT().FindByTicker(mock.Anything, "SBER").Return(identity, nil).Once()
	s.classifications.EXPECT().FindByTicker(mock.Anything, "SBER").Return(classification, nil).Once()

	got, err := s.service.GetCard(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal(companycard.Card{Identity: identity, Classification: classification}, got)
}

func (s *companyCardServiceSuite) TestGetCardPassesTickerAsIs() {
	// Источники регистронезависимы, сервис не нормализует ввод.
	s.identities.EXPECT().FindByTicker(mock.Anything, "sBeR").
		Return(companycard.Identity{Ticker: "sBeR"}, nil).Once()
	s.classifications.EXPECT().FindByTicker(mock.Anything, "sBeR").
		Return(companycard.Classification{}, nil).Once()

	_, err := s.service.GetCard(context.Background(), "sBeR")

	s.Require().NoError(err)
}

func (s *companyCardServiceSuite) TestGetCardEmptyTicker() {
	_, err := s.service.GetCard(context.Background(), "")

	s.Require().ErrorIs(err, services.ErrTickerEmpty)
}

func (s *companyCardServiceSuite) TestGetCardIdentityNotFoundPropagates() {
	s.identities.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(companycard.Identity{}, companycard.ErrNotFound).Once()
	s.classifications.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(companycard.Classification{}, nil).Maybe()

	_, err := s.service.GetCard(context.Background(), "ZZZZ")

	s.Require().ErrorIs(err, companycard.ErrNotFound)
}

func (s *companyCardServiceSuite) TestGetCardClassificationNotFoundPropagates() {
	s.identities.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(companycard.Identity{Ticker: "ZZZZ"}, nil).Maybe()
	s.classifications.EXPECT().FindByTicker(mock.Anything, "ZZZZ").
		Return(companycard.Classification{}, companycard.ErrNotFound).Once()

	_, err := s.service.GetCard(context.Background(), "ZZZZ")

	s.Require().ErrorIs(err, companycard.ErrNotFound)
}

func (s *companyCardServiceSuite) TestGetCardArbitraryErrorPropagates() {
	sentinel := errors.New("boom")
	s.identities.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(companycard.Identity{}, sentinel).Once()
	s.classifications.EXPECT().FindByTicker(mock.Anything, "SBER").
		Return(companycard.Classification{}, nil).Maybe()

	_, err := s.service.GetCard(context.Background(), "SBER")

	s.Require().ErrorIs(err, sentinel)
}
