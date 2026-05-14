// Package stockinfo — источник секции StockInfo поверх блока info
// эндпоинта FinanceMarker /api/fm/v2/stocks/{exchange}:{code}. Делает
// один HTTP-вызов, парсит ответ и возвращает заполненный
// company.StockInfo.
//
// Поддерживается только MOEX: единственная биржа, по которой проект
// возвращает карточки.
package stockinfo

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
	// includeInfo — значение query-параметра include, ограничивающее ответ
	// блоком info (классификация, описание, ссылки).
	includeInfo = "info"

	// codeExchangeMOEX — строковый код Московской биржи в формате FinanceMarker.
	codeExchangeMOEX = "MOEX"

	// cacheTTL — срок жизни HTTP-ответа /stocks в кеше клиента. Блок info
	// (классификация, описание эмитента, ссылки) меняется крайне редко,
	// поэтому раз в месяц достаточно. Источник декларирует свой TTL у
	// каждого исходящего запроса через ctx — фактическое хранение
	// принадлежит httpcache-уровню клиента (см. httpcache.WithTTL).
	cacheTTL = 30 * 24 * time.Hour
)

// Source — реализация company.StockInfoSource для блока info FinanceMarker.
// Кеширование живёт уровнем ниже, на http-transport клиента; источник
// только декларирует TTL у своего запроса через ctx (httpcache.WithTTL).
type Source struct {
	client *client.Client
}

// New собирает источник поверх общего FinanceMarker-клиента.
func New(c *client.Client) *Source {
	return &Source{client: c}
}

// FindByTicker запрашивает карточку эмитента, переводит блок info в
// StockInfo. Сетевые и HTTP-ошибки приходят из общего клиента уже
// классифицированными (см. client.New / classifyError), 404 здесь
// переводится в company.ErrNotFound.
//
// TTL для http-кеша выставляется на ctx через httpcache.WithTTL. Если
// у клиента кеш не подключён (CacheDir пустой) — аннотация остаётся в
// ctx без эффекта, и запрос идёт в сеть как есть.
func (s *Source) FindByTicker(ctx context.Context, ticker string) (company.StockInfo, error) {
	ctx = httpcache.WithTTL(ctx, cacheTTL)
	symbol := codeExchangeMOEX + ":" + ticker

	var dto stockDTO
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("symbol", symbol).
		SetQueryParam("include", includeInfo).
		SetResult(&dto).
		Get("/stocks/{symbol}")
	if err != nil {
		switch {
		case resp == nil || resp.StatusCode() == 0:
			return company.StockInfo{}, fmt.Errorf("financemarker request: %w", err)
		case !resp.IsError():
			return company.StockInfo{}, fmt.Errorf("decode financemarker payload: %w", err)
		case errors.Is(err, client.ErrNotFound):
			return company.StockInfo{}, company.ErrNotFound
		default:
			return company.StockInfo{}, err //nolint:wrapcheck // err уже сформирован classifyError общего клиента
		}
	}

	return translateStockInfo(&dto.Info), nil
}
