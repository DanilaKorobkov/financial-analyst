// Package stock — единый источник секций карточки эмитента поверх
// эндпоинта FinanceMarker /api/fm/v2/stocks/{exchange}:{code}.
//
// Источник принимает набор запрашиваемых блоков, собирает канонический
// query-параметр include (фиксированный порядок секций, без дубликатов) и делает
// один HTTP-вызов на все запрошенные секции. Это:
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
)

// ttlBySection — срок жизни HTTP-ответа в кеше клиента по секциям.
// При комбинированном запросе берётся минимум по запрошенным секциям —
// одна строка include соответствует одному кеш-ключу и одному TTL.
// StockSectionUnspecified присутствует с нулевым TTL и отсеивается в
// pickTTL.
var ttlBySection = map[company.StockSection]time.Duration{
	company.StockSectionUnspecified: 0,
	company.StockSectionInfo:        ttlInfo,
	company.StockSectionSummary:     ttlSummary,
}

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
// (см. canonicalIncludeOrder). 404 переводится в company.ErrNotFound;
// прочие HTTP-ошибки приходят из общего клиента уже классифицированными.
func (s *Source) FindByTicker(ctx context.Context, ticker string, sections []company.StockSection) (company.Stock, error) {
	include, err := buildInclude(sections)
	if err != nil {
		return company.Stock{}, fmt.Errorf("financemarker stock: %w", err)
	}

	ctx = httpcache.WithTTL(ctx, pickTTL(sections))
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

	return assemble(&dto, sections), nil
}

// pickTTL берёт минимальный TTL по запрошенным секциям. Гарантирует,
// что комбинированный ответ не пере-кешируется дольше самой «свежей»
// секции. Для пустого/неизвестного набора TTL получится нулевым — это
// безопасно, потому что в этом случае запрос отклоняется ещё в buildInclude.
func pickTTL(sections []company.StockSection) time.Duration {
	var ttl time.Duration
	for _, s := range sections {
		v, ok := ttlBySection[s]
		if !ok || v == 0 {
			continue
		}
		if ttl == 0 || v < ttl {
			ttl = v
		}
	}
	return ttl
}

// assemble обходит canonicalIncludeOrder и заполняет запрошенные секции
// в Stock — порядок сборки совпадает с порядком в include и с порядком
// полей company.Stock.
func assemble(dto *stockDTO, sections []company.StockSection) company.Stock {
	requested := make(map[company.StockSection]struct{}, len(sections))
	for _, s := range sections {
		requested[s] = struct{}{}
	}

	var out company.Stock
	for _, s := range canonicalIncludeOrder {
		if _, ok := requested[s]; !ok {
			continue
		}
		switch s {
		case company.StockSectionInfo:
			out.Info = translateStockInfo(&dto.Info)
		case company.StockSectionSummary:
			out.Summary = translateStockSummary(&dto.Summary)
		case company.StockSectionUnspecified:
			// не используется как запрос — отсеяно в buildInclude.
		}
	}
	return out
}
