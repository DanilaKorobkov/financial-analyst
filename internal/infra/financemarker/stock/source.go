// Package stock — единый источник секций карточки эмитента поверх
// эндпоинта FinanceMarker /api/fm/v2/stocks/{exchange}:{code}.
//
// Источник принимает набор запрашиваемых секций и делает по одному
// HTTP-запросу на каждую секцию параллельно. Каждый запрос несёт свой
// include=<code> и свой TTL HTTP-кеша, поэтому:
//
//   - длинные TTL (info=30d, shares=30d, reports/owners=7d) реально
//     работают — их не сбрасывает протухание короткой секции (ideas=6h);
//   - расход тарифа FinanceMarker не растёт: day_limit тарифицируется
//     по числу блоков в запросе (см. references/billing.md), и 10 блоков
//     в одном запросе и 10 запросов по блоку расходуют одинаково;
//   - cache hit на длинной секции экономит и тариф, и RTT.
//
// Поддерживается только MOEX: единственная биржа, по которой проект
// возвращает карточки.
package stock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sourcegraph/conc/pool"

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

// FindByTicker делает по одному HTTP-запросу на каждую запрошенную секцию
// параллельно и собирает результат в company.Stock. 404 переводится в
// company.ErrNotFound; прочие HTTP-ошибки приходят из общего клиента уже
// классифицированными. Первая ошибка любой параллельной горутины
// отменяет остальные через ctx (fail-fast).
func (s *Source) FindByTicker(ctx context.Context, ticker string, opts company.StockOptions) (company.Stock, error) {
	sections := enabledSections(opts)
	if len(sections) == 0 {
		return company.Stock{}, fmt.Errorf("financemarker stock: %w", errEmptyOptions)
	}

	symbol := codeExchangeMOEX + ":" + ticker
	p := pool.NewWithResults[func(*company.Stock)]().
		WithErrors().
		WithContext(ctx)
	for _, sec := range sections {
		p.Go(func(ctx context.Context) (func(*company.Stock), error) {
			return s.fetchSection(ctx, symbol, sec)
		})
	}

	appliers, err := p.Wait()
	if err != nil {
		return company.Stock{}, err //nolint:wrapcheck // ошибки уже сформированы в fetchSection
	}

	var out company.Stock
	for _, apply := range appliers {
		apply(&out)
	}
	return out, nil
}

// fetchSection делает один HTTP-запрос за одной секцией, аннотирует ctx
// её собственным TTL и возвращает аппликатор, который применяет результат
// к итоговому company.Stock. Маппинг ошибок идентичен исходному
// поведению объединённого запроса.
func (s *Source) fetchSection(
	ctx context.Context,
	symbol string,
	sec sectionFetch,
) (func(*company.Stock), error) {
	ctx = httpcache.WithTTL(ctx, sec.ttl)

	var dto stockDTO
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("symbol", symbol).
		SetQueryParam("include", sec.code).
		SetResult(&dto).
		Get("/stocks/{symbol}")
	if err != nil {
		switch {
		case resp == nil || resp.StatusCode() == 0:
			return nil, fmt.Errorf("financemarker request: %w", err)
		case !resp.IsError():
			return nil, fmt.Errorf("decode financemarker payload: %w", err)
		case errors.Is(err, client.ErrNotFound):
			return nil, company.ErrNotFound
		default:
			return nil, err //nolint:wrapcheck // err уже сформирован classifyError общего клиента
		}
	}

	return sec.apply(&dto), nil
}
