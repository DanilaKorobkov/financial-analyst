// Package moex — провайдер MOEX ISS: сборка общего HTTP-клиента
// (см. подпакет client) и регистрация всех bundles данного источника
// в data.Registry.
//
// MOEX — публичный источник без авторизации и без жёсткой квоты,
// поэтому bundles этого провайдера в кеш не оборачиваются.
package moex

import (
	"fmt"
	"time"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/bundles/securitydescription"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/client"
)

// ProviderID — стабильный идентификатор MOEX-провайдера в реестре.
const ProviderID = "moex"

// Provider собирает bundles MOEX-источника и регистрирует их в реестре.
// Реализует data.Provider.
type Provider struct {
	client *client.Client
}

// ConfigProvider — параметры Provider. Включает параметры HTTP-клиента;
// сам Client инкапсулирован внутри пакета — composition root о нём не знает.
type ConfigProvider struct {
	// BaseURL — корень MOEX ISS без завершающего слэша.
	BaseURL string

	// Timeout — таймаут на один HTTP-запрос.
	Timeout time.Duration
}

// NewProvider собирает Provider: внутри создаёт HTTP-клиент по конфигу.
func NewProvider(cfg ConfigProvider) *Provider {
	return &Provider{
		client: client.New(client.Config{
			BaseURL: cfg.BaseURL,
			Timeout: cfg.Timeout,
		}),
	}
}

// ID — реализация data.Provider.
func (*Provider) ID() string { return ProviderID }

// Register регистрирует все bundles MOEX в реестре.
func (p *Provider) Register(r data.Registrar) error {
	if err := r.Register(securitydescription.New(p.client)); err != nil {
		return fmt.Errorf("register %s/%s: %w", ProviderID, securitydescription.ID, err)
	}
	return nil
}
