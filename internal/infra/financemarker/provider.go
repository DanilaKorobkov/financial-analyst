// Package financemarker — провайдер FinanceMarker: сборка общего
// HTTP-клиента (см. подпакет client) и регистрация всех bundles
// данного источника в data.Registry.
//
// Bundles данного провайдера лежат в подпакете `bundles/` — по одному
// файлу на endpoint FinanceMarker. Файловый кеш живёт на http-transport
// клиента: каждый bundle декларирует у своего запроса TTL через
// httpcache.WithTTL, transport сам решает, отдать ответ из файла или
// сходить в сеть. Bundle про хранилище и формат записи ничего не знает.
package financemarker

import (
	"path/filepath"
	"time"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/bundles/stockinfo"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/client"
)

// ProviderID — стабильный идентификатор FinanceMarker-провайдера в реестре.
const ProviderID = "financemarker"

// Provider собирает bundles FinanceMarker-источника и регистрирует их
// в реестре. Реализует data.Provider.
//
// FinanceMarker — авторизованный источник с квотой на токен; кеширование
// HTTP-ответов выполняется единым transport-слоем клиента, а per-endpoint
// TTL задают сами bundles в своём Fetch.
type Provider struct {
	client *client.Client
}

// ConfigProvider — параметры Provider. Включает параметры HTTP-клиента
// и корень файлового кеша; сам Client скрыт внутри пакета —
// composition root о нём не знает.
type ConfigProvider struct {
	// BaseURL — корень FinanceMarker REST API без завершающего слэша.
	BaseURL string

	// Token — API-токен из профиля FinanceMarker.
	Token string

	// CacheRoot — корневой каталог файлового кеша HTTP-ответов. Пустая
	// строка отключает кеш: клиент идёт в сеть на каждый запрос. Под
	// каталог провайдера используется поддиректория `<CacheRoot>/<provider>`.
	CacheRoot string

	// Timeout — таймаут на один HTTP-запрос.
	Timeout time.Duration
}

// NewProvider собирает Provider: внутри создаёт HTTP-клиент по конфигу.
func NewProvider(cfg ConfigProvider) *Provider {
	clientCfg := client.Config{
		BaseURL: cfg.BaseURL,
		Token:   cfg.Token,
		Timeout: cfg.Timeout,
	}
	if cfg.CacheRoot != "" {
		clientCfg.CacheDir = filepath.Join(cfg.CacheRoot, ProviderID)
	}
	return &Provider{client: client.New(clientCfg)}
}

// ID — реализация data.Provider.
func (*Provider) ID() string { return ProviderID }

// Bundles — все bundles FinanceMarker-источника. Provider раздаёт
// им общий HTTP-клиент; кеширование bundles друг с другом и с
// composition root не координируют — оно происходит уровнем ниже,
// внутри клиента.
func (p *Provider) Bundles() []data.Bundle {
	return []data.Bundle{
		stockinfo.New(p.client),
	}
}
