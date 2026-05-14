package companyprofile_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/companyprofile"
)

type staticSuite struct {
	suite.Suite
}

func TestStaticSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(staticSuite))
}

func (s *staticSuite) TestFindByTickerReturnsConfiguredFields() {
	repo := companyprofile.NewStatic([]string{company.FieldTicker, company.FieldISIN})

	got, err := repo.FindByTicker(context.Background(), "any")

	s.Require().NoError(err)
	s.Equal([]string{company.FieldTicker, company.FieldISIN}, got.FieldIDs)
}

// TestFindByTickerSameForAnyTicker: Static игнорирует тикер по контракту —
// один и тот же набор полей на любую бумагу.
func (s *staticSuite) TestFindByTickerSameForAnyTicker() {
	repo := companyprofile.NewStatic([]string{company.FieldTicker})

	first, err := repo.FindByTicker(context.Background(), "SBER")
	s.Require().NoError(err)
	second, err := repo.FindByTicker(context.Background(), "GAZP")
	s.Require().NoError(err)

	s.Equal(first.FieldIDs, second.FieldIDs)
}

// TestFindByTickerIsolatedFromExternalSlice: входной срез не должен
// влиять на репозиторий после конструктора — иначе внешний код может
// поменять профиль из-под нас.
func (s *staticSuite) TestFindByTickerIsolatedFromExternalSlice() {
	input := []string{company.FieldTicker}
	repo := companyprofile.NewStatic(input)

	input[0] = "moex::other"

	got, err := repo.FindByTicker(context.Background(), "SBER")
	s.Require().NoError(err)
	s.Equal([]string{company.FieldTicker}, got.FieldIDs)
}

func (s *staticSuite) TestNewDefaultStaticHasFields() {
	repo := companyprofile.NewDefaultStatic()

	got, err := repo.FindByTicker(context.Background(), "SBER")
	s.Require().NoError(err)
	s.NotEmpty(got.FieldIDs)
	s.Contains(got.FieldIDs, company.FieldTicker)
}
