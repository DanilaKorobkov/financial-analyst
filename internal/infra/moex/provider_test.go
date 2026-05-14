package moex_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/bundles/securitydescription"
)

type providerSuite struct {
	suite.Suite
}

func TestProviderSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(providerSuite))
}

func (s *providerSuite) newProvider() *moex.Provider {
	return moex.NewProvider(moex.ConfigProvider{
		BaseURL: "http://example.invalid",
		Timeout: time.Second,
	})
}

func (s *providerSuite) TestID() {
	p := s.newProvider()
	s.Equal(moex.ProviderID, p.ID())
}

// TestRegisterAddsSecurityDescription: после регистрации в реестре
// появляется единственный bundle MOEX с правильной парой (provider, bundle).
func (s *providerSuite) TestRegisterAddsSecurityDescription() {
	p := s.newProvider()
	reg := data.NewRegistry()
	s.Require().NoError(reg.RegisterProvider(p))

	b, err := reg.Bundle(moex.ProviderID, securitydescription.ID)
	s.Require().NoError(err)
	s.Equal(securitydescription.ID, b.BundleID())
}

// TestRegisterPropagatesError: если реестр отказывается принять bundle
// (например, тот же ключ уже занят), Provider возвращает ошибку с именем
// своей пары — composition root по сообщению поймёт, кто упал.
func (s *providerSuite) TestRegisterPropagatesError() {
	p := s.newProvider()
	reg := data.NewRegistry()
	s.Require().NoError(reg.RegisterProvider(p))

	err := reg.RegisterProvider(p)
	s.Require().Error(err)
	s.ErrorContains(err, moex.ProviderID)
	s.ErrorContains(err, securitydescription.ID)
}
