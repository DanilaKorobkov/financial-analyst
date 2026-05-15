// Package stock — единый источник секций карточки эмитента поверх
// эндпоинта FinanceMarker /api/fm/v2/stocks/{exchange}:{code}.
//
// Источник принимает набор запрашиваемых блоков, собирает канонический
// query-параметр include (фиксированный порядок секций, без дубликатов)
// и делает один HTTP-вызов на все запрошенные секции. Это:
//
//   - снижает HTTP overhead и шанс расхождения cache-ключа кеша
//     FinanceMarker (один use case → одна строка include, см. справочник
//     api-financemarker/references/stock.md);
//   - оставляет тарификацию day_limit равной фактическому числу блоков
//     (FM считает по блокам внутри запроса, см. references/billing.md).
//
// Поддерживается только MOEX: единственная биржа, по которой проект
// возвращает карточки.
package stock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/cache/httpcache"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/client"
)

const (
	// codeExchangeMOEX — строковый код Московской биржи в формате FinanceMarker.
	codeExchangeMOEX = "MOEX"

	// ttlInfo — TTL HTTP-кеша блока info: классификация и описание эмитента
	// меняются крайне редко.
	ttlInfo = 30 * 24 * time.Hour

	// ttlSummary — TTL HTTP-кеша блока summary: пересчитывается при выходе
	// свежих отчётов и изменении цены.
	ttlSummary = 24 * time.Hour

	// ttlRatios — TTL HTTP-кеша блока ratios: пересчёт привязан к новой
	// отчётности (как summary).
	ttlRatios = 24 * time.Hour

	// ttlReports — TTL HTTP-кеша блока reports: отчёты публикуются раз
	// в квартал; за неделю гарантированно подхватим новый.
	ttlReports = 7 * 24 * time.Hour

	// ttlDividends — TTL HTTP-кеша блока dividends: раскрытие новой
	// рекомендации может появиться в любой день.
	ttlDividends = 24 * time.Hour

	// ttlIdeas — TTL HTTP-кеша блока ideas: лента инвест-идей живая,
	// статусы и цены меняются в течение дня.
	ttlIdeas = 6 * time.Hour

	// ttlInsiderTransactions — TTL HTTP-кеша блока insiderTransactions.
	ttlInsiderTransactions = 24 * time.Hour

	// ttlOperations — TTL HTTP-кеша блока operations: операционные метрики
	// обновляются вместе с отчётами эмитента.
	ttlOperations = 24 * time.Hour

	// ttlOwners — TTL HTTP-кеша блока owners: структура акционеров
	// меняется редко.
	ttlOwners = 7 * 24 * time.Hour

	// ttlShares — TTL HTTP-кеша блока shares: количество акций — почти
	// статичный ряд, пересчёт ежемесячно с большим запасом.
	ttlShares = 30 * 24 * time.Hour
)

// Source — реализация company.StockSource поверх FinanceMarker.
type Source struct {
	client *client.Client
}

// New собирает источник поверх общего FinanceMarker-клиента.
func New(c *client.Client) *Source {
	return &Source{client: c}
}

// FindByTicker делает один запрос за указанными секциями и заполняет
// только их в company.Stock. Порядок секций в include всегда канонический
// (info, summary). 404 переводится в company.ErrNotFound; прочие
// HTTP-ошибки приходят из общего клиента уже классифицированными.
func (s *Source) FindByTicker(ctx context.Context, ticker string, opts company.StockOptions) (company.Stock, error) {
	include, err := buildInclude(opts)
	if err != nil {
		return company.Stock{}, fmt.Errorf("financemarker stock: %w", err)
	}

	ctx = httpcache.WithTTL(ctx, pickTTL(opts))
	symbol := codeExchangeMOEX + ":" + ticker

	var dto stockDTO
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("symbol", symbol).
		SetQueryParam("include", include).
		SetResult(&dto).
		Get("/stocks/{symbol}")
	if err != nil {
		switch {
		case resp == nil || resp.StatusCode() == 0:
			return company.Stock{}, fmt.Errorf("financemarker request: %w", err)
		case !resp.IsError():
			return company.Stock{}, fmt.Errorf("decode financemarker payload: %w", err)
		case errors.Is(err, client.ErrNotFound):
			return company.Stock{}, company.ErrNotFound
		default:
			return company.Stock{}, err //nolint:wrapcheck // err уже сформирован classifyError общего клиента
		}
	}

	return assemble(&dto, opts), nil
}

// pickTTL берёт минимальный TTL по запрошенным секциям. Гарантирует,
// что комбинированный ответ не пере-кешируется дольше самой «свежей»
// секции — одна строка include = один кеш-ключ = один TTL.
func pickTTL(opts company.StockOptions) time.Duration {
	var ttl time.Duration
	consider := func(enabled bool, sectionTTL time.Duration) {
		if !enabled {
			return
		}
		if ttl == 0 || sectionTTL < ttl {
			ttl = sectionTTL
		}
	}
	consider(opts.WithInfo, ttlInfo)
	consider(opts.WithSummary, ttlSummary)
	consider(opts.WithRatios, ttlRatios)
	consider(opts.WithReports, ttlReports)
	consider(opts.WithDividends, ttlDividends)
	consider(opts.WithIdeas, ttlIdeas)
	consider(opts.WithInsiderTransactions, ttlInsiderTransactions)
	consider(opts.WithOperations, ttlOperations)
	consider(opts.WithOwners, ttlOwners)
	consider(opts.WithShares, ttlShares)
	return ttl
}

// assemble заполняет запрошенные секции в Stock — порядок сборки
// совпадает с порядком полей company.Stock.
func assemble(dto *stockDTO, opts company.StockOptions) company.Stock {
	var out company.Stock
	if opts.WithInfo {
		out.Info = translateStockInfo(&dto.Info)
	}
	if opts.WithSummary {
		out.Summary = translateStockSummary(&dto.Summary)
	}
	if opts.WithRatios {
		out.Ratios = translateStockRatios(dto.Ratios)
	}
	if opts.WithReports {
		out.Reports = translateStockReports(dto.Reports)
	}
	if opts.WithDividends {
		out.Dividends = translateStockDividends(dto.Dividends)
	}
	if opts.WithIdeas {
		out.Ideas = translateStockIdeas(dto.Ideas)
	}
	if opts.WithInsiderTransactions {
		out.InsiderTransactions = translateStockInsiderTransactions(dto.InsiderTransactions)
	}
	if opts.WithOperations {
		out.Operations = translateStockOperations(dto.Operations)
	}
	if opts.WithOwners {
		out.Owners = translateStockOwners(dto.Owners)
	}
	if opts.WithShares {
		out.Shares = translateStockShares(dto.Shares)
	}
	return out
}
