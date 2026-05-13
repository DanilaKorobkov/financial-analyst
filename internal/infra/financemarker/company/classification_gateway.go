// Package company — реализация company.ClassificationGateway поверх
// FinanceMarker /api/fm/v2/stocks/{exchange}:{code} (блок info).
package company

import (
	"context"
	"errors"
	"fmt"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker"
)

const (
	// includeInfo — значение query-параметра include, ограничивающее ответ
	// блоком info (классификация, описание, ссылки).
	includeInfo = "info"

	// codeExchangeMOEX — строковый код Московской биржи в формате FinanceMarker.
	codeExchangeMOEX = "MOEX"
)

// ClassificationGateway ходит в /stocks/{exchange}:{ticker} (блок info).
// Поддерживается только MOEX: единственная биржа, по которой проект
// возвращает карточки.
type ClassificationGateway struct {
	client *financemarker.Client
}

// NewClassificationGateway собирает gateway поверх общего FinanceMarker-клиента.
func NewClassificationGateway(client *financemarker.Client) *ClassificationGateway {
	return &ClassificationGateway{client: client}
}

// FindByTicker запрашивает карточку эмитента и переводит классификационный
// блок info в domaincompany.Classification.
//
// Сетевые и HTTP-ошибки приходят из общего клиента уже классифицированными
// (см. financemarker.NewClient / classifyError), 404 здесь переводится в
// domaincompany.ErrNotFound.
func (g *ClassificationGateway) FindByTicker(
	ctx context.Context,
	ticker string,
) (domaincompany.Classification, error) {
	symbol := codeExchangeMOEX + ":" + ticker

	var dto stockDTO
	resp, err := g.client.R().
		SetContext(ctx).
		SetPathParam("symbol", symbol).
		SetQueryParam("include", includeInfo).
		SetResult(&dto).
		Get("/stocks/{symbol}")
	if err != nil {
		switch {
		case resp == nil || resp.StatusCode() == 0:
			return domaincompany.Classification{}, fmt.Errorf("financemarker request: %w", err)
		case !resp.IsError():
			return domaincompany.Classification{}, fmt.Errorf("decode financemarker payload: %w", err)
		case errors.Is(err, financemarker.ErrNotFound):
			return domaincompany.Classification{}, domaincompany.ErrNotFound
		default:
			return domaincompany.Classification{}, err //nolint:wrapcheck // err уже сформирован classifyError общего клиента
		}
	}

	return translateClassification(&dto.Info), nil
}
