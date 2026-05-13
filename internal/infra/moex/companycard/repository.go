// Package companycard — реализация companycard.IdentityGateway поверх
// блока description эндпоинта MOEX ISS /iss/securities/{TICKER}.json.
package companycard

import (
	"context"
	"fmt"

	domaincard "github.com/DanilaKorobkov/financial-analyst/internal/domain/companycard"
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

// FindByTicker запрашивает description у MOEX и переводит в domaincard.Identity.
//
// Возвращает domaincard.ErrNotFound, если ISS вернула пустой блок
// description (тикер не существует).
func (g *IdentityGateway) FindByTicker(ctx context.Context, ticker string) (domaincard.Identity, error) {
	resp, err := g.client.R().
		SetContext(ctx).
		SetPathParam("ticker", ticker).
		Get("/securities/{ticker}.json")
	if err != nil {
		if resp == nil || resp.StatusCode() == 0 {
			return domaincard.Identity{}, fmt.Errorf("moex request: %w", err)
		}
		return domaincard.Identity{}, err //nolint:wrapcheck // err уже сформирован OnAfterResponse в moex.NewClient ("moex http status N")
	}

	fields, err := parseDescription(resp.Body())
	if err != nil {
		return domaincard.Identity{}, fmt.Errorf("parse description: %w", err)
	}

	return mapDescription(fields)
}
