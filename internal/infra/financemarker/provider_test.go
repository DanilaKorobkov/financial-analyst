package financemarker_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/bundles/stockinfo"
)

type providerSuite struct {
	suite.Suite
}

func TestProviderSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(providerSuite))
}

func (s *providerSuite) newProvider() *financemarker.Provider {
	return financemarker.NewProvider(financemarker.ConfigProvider{
		BaseURL:   "http://example.invalid",
		Token:     "x",
		Timeout:   time.Second,
		CacheRoot: s.T().TempDir(),
	})
}

func (s *providerSuite) TestID() {
	p := s.newProvider()
	s.Equal(financemarker.ProviderID, p.ID())
}

// TestRegisterAddsStockInfo: после регистрации в реестре появляется
// bundle FinanceMarker с правильной парой (provider, bundle).
// Регистрируется именно кеширующая обёртка — она транзитом отдаёт
// BundleID делегата, поэтому проверяем по этим идентификаторам.
func (s *providerSuite) TestRegisterAddsStockInfo() {
	p := s.newProvider()
	reg := data.NewRegistry()
	s.Require().NoError(reg.RegisterProvider(p))

	b, err := reg.Bundle(financemarker.ProviderID, stockinfo.ID)
	s.Require().NoError(err)
	s.Equal(stockinfo.ID, b.BundleID())
}

// TestRegisterPropagatesError: повторная регистрация — ошибка с
// именем пары, чтобы composition root по сообщению поймал, кто упал.
func (s *providerSuite) TestRegisterPropagatesError() {
	p := s.newProvider()
	reg := data.NewRegistry()
	s.Require().NoError(reg.RegisterProvider(p))

	err := reg.RegisterProvider(p)
	s.Require().Error(err)
	s.ErrorContains(err, financemarker.ProviderID)
	s.ErrorContains(err, stockinfo.ID)
}
