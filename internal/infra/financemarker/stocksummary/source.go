// Package stocksummary — источник секции StockSummary поверх блока summary
// эндпоинта FinanceMarker /api/fm/v2/stocks/{exchange}:{code}. Делает
// один HTTP-вызов, парсит ответ и возвращает заполненный
// company.StockSummary.
//
// Поддерживается только MOEX: единственная биржа, по которой проект
// возвращает карточки.
package stocksummary

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
	// includeSummary — значение query-параметра include, ограничивающее ответ
	// блоком summary (сводные метрики эмитента).
	includeSummary = "summary"

	// codeExchangeMOEX — строковый код Московской биржи в формате FinanceMarker.
	codeExchangeMOEX = "MOEX"

	// cacheTTL — срок жизни HTTP-ответа /stocks в кеше клиента. Блок summary
	// пересчитывается при выходе свежих отчётов и изменении цены: чаще,
	// чем info, реже, чем ratios/reports. Раз в сутки — разумный компромисс.
	cacheTTL = 24 * time.Hour
)

// Source — реализация company.StockSummarySource для блока summary FinanceMarker.
// Кеширование живёт уровнем ниже, на http-transport клиента; источник
// только декларирует TTL у своего запроса через ctx (httpcache.WithTTL).
type Source struct {
	client *client.Client
}

// New собирает источник поверх общего FinanceMarker-клиента.
func New(c *client.Client) *Source {
	return &Source{client: c}
}

// FindByTicker запрашивает сводные метрики эмитента и переводит блок
// summary в StockSummary. Сетевые и HTTP-ошибки приходят из общего
// клиента уже классифицированными (см. client.New / classifyError),
// 404 здесь переводится в company.ErrNotFound.
//
// TTL для http-кеша выставляется на ctx через httpcache.WithTTL. Если
// у клиента кеш не подключён (CacheDir пустой) — аннотация остаётся в
// ctx без эффекта, и запрос идёт в сеть как есть.
func (s *Source) FindByTicker(ctx context.Context, ticker string) (company.StockSummary, error) {
	ctx = httpcache.WithTTL(ctx, cacheTTL)
	symbol := codeExchangeMOEX + ":" + ticker

	var dto stockDTO
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("symbol", symbol).
		SetQueryParam("include", includeSummary).
		SetResult(&dto).
		Get("/stocks/{symbol}")
	if err != nil {
		switch {
		case resp == nil || resp.StatusCode() == 0:
			return company.StockSummary{}, fmt.Errorf("financemarker request: %w", err)
		case !resp.IsError():
			return company.StockSummary{}, fmt.Errorf("decode financemarker payload: %w", err)
		case errors.Is(err, client.ErrNotFound):
			return company.StockSummary{}, company.ErrNotFound
		default:
			return company.StockSummary{}, err //nolint:wrapcheck // err уже сформирован classifyError общего клиента
		}
	}

	return translateStockSummary(&dto.Summary), nil
}
