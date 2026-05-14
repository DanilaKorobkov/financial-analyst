// Package financemarker — провайдер FinanceMarker: сборка общего
// HTTP-клиента (см. подпакет client) и регистрация всех bundles
// данного источника в data.Registry.
//
// Bundles данного провайдера лежат в подпакете `bundles/` — по одному
// файлу на endpoint FinanceMarker. Файловый кеш над каждым bundle —
// внутренняя деталь провайдера: квота на API-токен жёсткая, поэтому
// все bundles этого источника по умолчанию кешируются.
package financemarker

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	fcbundle "github.com/DanilaKorobkov/financial-analyst/internal/infra/filecache/bundle"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/bundles/stockinfo"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/client"
)

const (
	// ProviderID — стабильный идентификатор FinanceMarker-провайдера в реестре.
	ProviderID = "financemarker"

	// stockInfoCacheTTL — TTL файлового кеша stock-info. Хардкод вместо env:
	// операционный параметр одного bundle, меняется крайне редко. Когда
	// понадобится тонкая настройка под разные bundles — вынесем в
	// ConfigProvider или в derived freshness (см. VISION).
	stockInfoCacheTTL = 720 * time.Hour
)

// Provider собирает bundles FinanceMarker-источника и регистрирует их
// в реестре. Реализует data.Provider.
//
// FinanceMarker — авторизованный источник с квотой на токен, поэтому
// все bundles этого провайдера оборачиваются в файловый кеш. Кеш —
// внутренняя деталь провайдера: composition root о ней не знает,
// просто отдаёт корневой каталог под файлы.
type Provider struct {
	client    *client.Client
	cacheRoot string
}

// ConfigProvider — параметры Provider. Включает параметры HTTP-клиента
// и корень файлового кеша; сам Client инкапсулирован внутри пакета —
// composition root о нём не знает.
type ConfigProvider struct {
	// BaseURL — корень FinanceMarker REST API без завершающего слэша.
	BaseURL string

	// Token — API-токен из профиля FinanceMarker.
	Token string

	// CacheRoot — корневой каталог файлового кеша. Под каждый bundle
	// провайдер создаёт собственный подкаталог `<CacheRoot>/<provider>/<bundle>`.
	CacheRoot string

	// Timeout — таймаут на один HTTP-запрос.
	Timeout time.Duration
}

// NewProvider собирает Provider: внутри создаёт HTTP-клиент по конфигу.
func NewProvider(cfg ConfigProvider) *Provider {
	return &Provider{
		client: client.New(client.Config{
			BaseURL: cfg.BaseURL,
			Token:   cfg.Token,
			Timeout: cfg.Timeout,
		}),
		cacheRoot: cfg.CacheRoot,
	}
}

// ID — реализация data.Provider.
func (*Provider) ID() string { return ProviderID }

// Register регистрирует все bundles FinanceMarker в реестре, оборачивая
// каждый в файловый кеш.
func (p *Provider) Register(r data.Registrar) error {
	cached := fcbundle.NewProxy(fcbundle.ConfigProxy{
		Delegate: stockinfo.New(p.client),
		Dir:      filepath.Join(p.cacheRoot, ProviderID, stockinfo.ID),
		TTL:      stockInfoCacheTTL,
	})
	if err := r.Register(cached); err != nil {
		return fmt.Errorf("register %s/%s: %w", ProviderID, stockinfo.ID, err)
	}
	return nil
}
