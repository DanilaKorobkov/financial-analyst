// Package companycard — реализация companycard.Repository поверх
// FinanceMarker /api/fm/v2/stocks/{exchange}:{code}.
package companycard

import (
	"context"
	"errors"
	"fmt"

	domaincard "github.com/DanilaKorobkov/financial-analyst/internal/domain/companycard"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker"
)

const (
	// includeInfo — значение query-параметра include, ограничивающее ответ
	// блоком info (классификация, описание, ссылки).
	includeInfo = "info"

	// codeExchangeMOEX — строковый код Московской биржи в формате FinanceMarker.
	codeExchangeMOEX = "MOEX"
)

// Repository ходит в /stocks/{exchange}:{ticker} (блок info).
type Repository struct {
	client *financemarker.Client
}

// NewRepository собирает репозиторий поверх общего FinanceMarker-клиента.
func NewRepository(client *financemarker.Client) *Repository {
	return &Repository{client: client}
}

// FindByTicker запрашивает карточку эмитента и переводит её в domaincard.Card.
// Биржу принимаем явно: символ для FM собирается как "{exchange}:{ticker}".
// Сетевые и HTTP-ошибки приходят из общего клиента уже классифицированными
// (см. financemarker.NewClient / classifyError), 404 здесь переводится в
// domaincard.ErrNotFound.
func (r *Repository) FindByTicker(
	ctx context.Context,
	exchange domaincard.Exchange,
	ticker string,
) (domaincard.Card, error) {
	symbol, err := buildSymbol(exchange, ticker)
	if err != nil {
		return domaincard.Card{}, err
	}

	var dto stockDTO
	resp, err := r.client.R().
		SetContext(ctx).
		SetPathParam("symbol", symbol).
		SetQueryParam("include", includeInfo).
		SetResult(&dto).
		Get("/stocks/{symbol}")
	if err != nil {
		switch {
		case resp == nil || resp.StatusCode() == 0:
			return domaincard.Card{}, fmt.Errorf("financemarker request: %w", err)
		case !resp.IsError():
			return domaincard.Card{}, fmt.Errorf("decode financemarker payload: %w", err)
		case errors.Is(err, financemarker.ErrNotFound):
			return domaincard.Card{}, domaincard.ErrNotFound
		default:
			return domaincard.Card{}, err //nolint:wrapcheck // err уже сформирован classifyError общего клиента
		}
	}

	return translateCard(&dto.Info), nil
}

// buildSymbol собирает path-параметр FinanceMarker вида "{exchange}:{ticker}".
// Неподдерживаемая биржа — ошибка: запрос с непустым кодом, который FM не
// признаёт, всё равно вернёт 404 / 400, и лучше отвалиться явно.
func buildSymbol(exchange domaincard.Exchange, ticker string) (string, error) {
	code, err := exchangeCode(exchange)
	if err != nil {
		return "", err
	}
	return code + ":" + ticker, nil
}

// exchangeCode возвращает строковый код биржи в формате FinanceMarker.
func exchangeCode(exchange domaincard.Exchange) (string, error) {
	switch exchange {
	case domaincard.ExchangeMOEX:
		return codeExchangeMOEX, nil
	case domaincard.ExchangeUnspecified:
		return "", fmt.Errorf("financemarker: exchange is unspecified")
	default:
		return "", fmt.Errorf("financemarker: unsupported exchange %d", exchange)
	}
}
