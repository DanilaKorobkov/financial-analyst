// Package securitydescription — источник секции SecurityDescription
// поверх блока description эндпоинта MOEX ISS /iss/securities/{TICKER}.json.
// Делает один HTTP-вызов, парсит ответ и возвращает заполненный
// company.SecurityDescription.
package securitydescription

import (
	"context"
	"fmt"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/client"
)

// Source — реализация company.SecurityDescriptionSource для блока
// description MOEX ISS.
type Source struct {
	client *client.Client
}

// New собирает источник поверх общего MOEX-клиента.
func New(c *client.Client) *Source {
	return &Source{client: c}
}

// FindByTicker запрашивает description у MOEX ISS и возвращает заполненный
// SecurityDescription. Если ISS вернула пустой блок description (тикер
// не существует), возвращает company.ErrNotFound.
func (s *Source) FindByTicker(ctx context.Context, ticker string) (company.SecurityDescription, error) {
	resp, err := s.client.R().
		SetContext(ctx).
		SetPathParam("ticker", ticker).
		Get("/securities/{ticker}.json")
	if err != nil {
		if resp == nil || resp.StatusCode() == 0 {
			return company.SecurityDescription{}, fmt.Errorf("moex request: %w", err)
		}
		return company.SecurityDescription{}, err //nolint:wrapcheck // err уже сформирован OnAfterResponse в client.New ("moex http status N")
	}

	parsed, err := parseDescription(resp.Body())
	if err != nil {
		return company.SecurityDescription{}, fmt.Errorf("parse description: %w", err)
	}

	return mapDescription(parsed)
}
