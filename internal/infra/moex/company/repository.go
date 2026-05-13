// Package company — реализация company.Repository поверх блока description
// эндпоинта MOEX ISS /iss/securities/{TICKER}.json.
package company

import (
	"context"
	"fmt"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
)

// Repository ходит в /iss/securities/{TICKER}.json (блок description).
type Repository struct {
	client *moex.Client
}

// NewRepository собирает репозиторий поверх общего MOEX-клиента.
func NewRepository(client *moex.Client) *Repository {
	return &Repository{client: client}
}

// FindByTicker запрашивает description у MOEX и переводит в domaincompany.Company.
//
// Возвращает domaincompany.ErrNotFound, если ISS вернула пустой блок
// description (тикер не существует).
func (r *Repository) FindByTicker(ctx context.Context, ticker string) (domaincompany.Company, error) {
	resp, err := r.client.R().
		SetContext(ctx).
		SetPathParam("ticker", ticker).
		Get("/securities/{ticker}.json")
	if err != nil {
		if resp == nil || resp.StatusCode() == 0 {
			return domaincompany.Company{}, fmt.Errorf("moex request: %w", err)
		}
		return domaincompany.Company{}, err //nolint:wrapcheck // err уже сформирован OnAfterResponse в moex.NewClient ("moex http status N")
	}

	fields, err := parseDescription(resp.Body())
	if err != nil {
		return domaincompany.Company{}, fmt.Errorf("parse description: %w", err)
	}

	return mapDescription(fields)
}
