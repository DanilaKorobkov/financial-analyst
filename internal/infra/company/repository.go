// Package company — infra-реализация company.Repository: собирает
// агрегат Company из источников секций (SecurityDescriptionSource —
// описание ценной бумаги, StockSource — карточка эмитента и сводные
// метрики одним вызовом). Источники вызываются параллельно через
// conc/pool с отменой контекста по первой ошибке.
package company

import (
	"context"
	"fmt"

	"github.com/sourcegraph/conc/pool"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

// stockOptions — секции карточки эмитента, которые входят в агрегат Company.
var stockOptions = domaincompany.StockOptions{
	WithInfo:    true,
	WithSummary: true,
}

// Repository — infra-реализация domaincompany.Repository.
type Repository struct {
	securityDescription domaincompany.SecurityDescriptionSource
	stock               domaincompany.StockSource
}

// ConfigRepository — параметры Repository.
type ConfigRepository struct {
	// SecurityDescription — источник описания ценной бумаги.
	SecurityDescription domaincompany.SecurityDescriptionSource

	// Stock — источник секций карточки эмитента (StockInfo + StockSummary)
	// одним вызовом.
	Stock domaincompany.StockSource
}

// NewRepository собирает Repository поверх источников секций.
func NewRepository(cfg ConfigRepository) *Repository {
	return &Repository{
		securityDescription: cfg.SecurityDescription,
		stock:               cfg.Stock,
	}
}

// FindByTicker качает все секции параллельно. Любая ошибка источника
// (включая domaincompany.ErrNotFound) fail-fast отменяет пул и
// возвращается наверх — агрегат отдаётся только целиком.
func (r *Repository) FindByTicker(ctx context.Context, ticker string) (domaincompany.Company, error) {
	var agg domaincompany.Company

	p := pool.New().WithErrors().WithContext(ctx).WithCancelOnError().WithFirstError()
	p.Go(func(ctx context.Context) error {
		section, err := r.securityDescription.FindByTicker(ctx, ticker)
		if err != nil {
			return fmt.Errorf("security description: %w", err)
		}
		agg.SecurityDescription = section
		return nil
	})
	p.Go(func(ctx context.Context) error {
		s, err := r.stock.FindByTicker(ctx, ticker, stockOptions)
		if err != nil {
			return fmt.Errorf("stock: %w", err)
		}
		agg.Stock = s
		return nil
	})
	if err := p.Wait(); err != nil {
		return domaincompany.Company{}, fmt.Errorf("fetch sections: %w", err)
	}

	return agg, nil
}
