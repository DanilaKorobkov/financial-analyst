// Package company — реализация company.IdentityGateway поверх блока
// description эндпоинта MOEX ISS /iss/securities/{TICKER}.json.
package company

import (
	"context"
	"fmt"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
)

// IdentityGateway ходит в /iss/securities/{TICKER}.json (блок description).
type IdentityGateway struct {
	client *moex.Client
}

// NewIdentityGateway собирает gateway поверх общего MOEX-клиента.
func NewIdentityGateway(client *moex.Client) *IdentityGateway {
	return &IdentityGateway{client: client}
}

// FindByTicker запрашивает description у MOEX и переводит в domaincompany.Identity.
//
// Возвращает domaincompany.ErrNotFound, если ISS вернула пустой блок
// description (тикер не существует).
func (g *IdentityGateway) FindByTicker(ctx context.Context, ticker string) (domaincompany.Identity, error) {
	resp, err := g.client.R().
		SetContext(ctx).
		SetPathParam("ticker", ticker).
		Get("/securities/{ticker}.json")
	if err != nil {
		if resp == nil || resp.StatusCode() == 0 {
			return domaincompany.Identity{}, fmt.Errorf("moex request: %w", err)
		}
		return domaincompany.Identity{}, err //nolint:wrapcheck // err уже сформирован OnAfterResponse в moex.NewClient ("moex http status N")
	}

	fields, err := parseDescription(resp.Body())
	if err != nil {
		return domaincompany.Identity{}, fmt.Errorf("parse description: %w", err)
	}

	return mapDescription(fields)
}
