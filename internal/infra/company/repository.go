// Package company — infra-реализация company.Repository: собирает
// агрегат Company из источников секций (SecurityDescriptionSource,
// StockInfoSource). Источники вызываются параллельно через conc/pool
// с отменой контекста по первой не-ErrNotFound ошибке.
package company

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/sourcegraph/conc/pool"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

// Repository — infra-реализация domaincompany.Repository.
type Repository struct {
	securityDescription domaincompany.SecurityDescriptionSource
	stockInfo           domaincompany.StockInfoSource
}

// ConfigRepository — параметры Repository.
type ConfigRepository struct {
	// SecurityDescription — источник описания ценной бумаги.
	SecurityDescription domaincompany.SecurityDescriptionSource

	// StockInfo — источник карточки эмитента.
	StockInfo domaincompany.StockInfoSource
}

// NewRepository собирает Repository поверх источников секций.
func NewRepository(cfg ConfigRepository) *Repository {
	return &Repository{
		securityDescription: cfg.SecurityDescription,
		stockInfo:           cfg.StockInfo,
	}
}

// FindByTicker качает обе секции параллельно. Семантика: ErrNotFound от
// источника превращается в nil-секцию без ошибки; любая другая ошибка
// fail-fast отменяет пул. Если ОБА источника ответили ErrNotFound, возвращает
// domaincompany.ErrCompanyNotFound.
func (r *Repository) FindByTicker(ctx context.Context, ticker string) (*domaincompany.Company, error) {
	agg := &domaincompany.Company{}
	var mu sync.Mutex

	p := pool.New().WithErrors().WithContext(ctx).WithCancelOnError().WithFirstError()
	p.Go(func(ctx context.Context) error {
		section, err := r.securityDescription.FindByTicker(ctx, ticker)
		if err != nil {
			if errors.Is(err, domaincompany.ErrNotFound) {
				return nil
			}
			return fmt.Errorf("security description: %w", err)
		}
		mu.Lock()
		agg.SecurityDescription = section
		mu.Unlock()
		return nil
	})
	p.Go(func(ctx context.Context) error {
		section, err := r.stockInfo.FindByTicker(ctx, ticker)
		if err != nil {
			if errors.Is(err, domaincompany.ErrNotFound) {
				return nil
			}
			return fmt.Errorf("stock info: %w", err)
		}
		mu.Lock()
		agg.StockInfo = section
		mu.Unlock()
		return nil
	})
	if err := p.Wait(); err != nil {
		return nil, fmt.Errorf("fetch sections: %w", err)
	}

	if agg.SecurityDescription == nil && agg.StockInfo == nil {
		return nil, domaincompany.ErrCompanyNotFound
	}
	return agg, nil
}
