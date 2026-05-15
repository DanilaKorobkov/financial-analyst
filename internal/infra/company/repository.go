// Package company — infra-реализация company.Repository: собирает
// агрегат Company из источников секций (SecurityDescriptionSource,
// StockInfoSource, StockSummarySource). Источники вызываются параллельно
// через conc/pool с отменой контекста по первой ошибке.
package company

import (
	"context"
	"fmt"

	"github.com/sourcegraph/conc/pool"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

// Repository — infra-реализация domaincompany.Repository.
type Repository struct {
	securityDescription domaincompany.SecurityDescriptionSource
	stockInfo           domaincompany.StockInfoSource
	stockSummary        domaincompany.StockSummarySource
}

// ConfigRepository — параметры Repository.
type ConfigRepository struct {
	// SecurityDescription — источник описания ценной бумаги.
	SecurityDescription domaincompany.SecurityDescriptionSource

	// StockInfo — источник карточки эмитента.
	StockInfo domaincompany.StockInfoSource

	// StockSummary — источник сводных метрик эмитента.
	StockSummary domaincompany.StockSummarySource
}

// NewRepository собирает Repository поверх источников секций.
func NewRepository(cfg ConfigRepository) *Repository {
	return &Repository{
		securityDescription: cfg.SecurityDescription,
		stockInfo:           cfg.StockInfo,
		stockSummary:        cfg.StockSummary,
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
		section, err := r.stockInfo.FindByTicker(ctx, ticker)
		if err != nil {
			return fmt.Errorf("stock info: %w", err)
		}
		agg.StockInfo = section
		return nil
	})
	p.Go(func(ctx context.Context) error {
		section, err := r.stockSummary.FindByTicker(ctx, ticker)
		if err != nil {
			return fmt.Errorf("stock summary: %w", err)
		}
		agg.StockSummary = section
		return nil
	})
	if err := p.Wait(); err != nil {
		return domaincompany.Company{}, fmt.Errorf("fetch sections: %w", err)
	}

	return agg, nil
}
