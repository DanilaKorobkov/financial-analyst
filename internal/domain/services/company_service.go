// Package services — domain-сервисы, оркеструют агрегаты и порты.
package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/sourcegraph/conc/pool"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
)

// ErrTickerEmpty — клиент передал пустой тикер.
var ErrTickerEmpty = errors.New("ticker is empty")

// CompanyService — сервис сборки компании по тикеру. Каждая секция
// агрегата приходит из своего внешнего источника-gateway; все
// источники, на которые сервис подписан, обязательные — частичная
// карточка не возвращается.
type CompanyService struct {
	identities      company.IdentityGateway
	classifications company.ClassificationGateway
}

// NewCompanyService собирает сервис вокруг двух gateway.
func NewCompanyService(
	identities company.IdentityGateway,
	classifications company.ClassificationGateway,
) *CompanyService {
	return &CompanyService{
		identities:      identities,
		classifications: classifications,
	}
}

// GetCompany проверяет непустоту тикера и параллельно тянет обе
// секции. При ошибке любого gateway возвращает её, не дожидаясь
// второго: WithFirstError отдаёт первую ошибку как есть,
// WithCancelOnError отменяет ctx второй горутины. Тикер передаётся
// как есть, без нормализации.
//
// Горутины пишут в разные поля одного экземпляра Company — гонок нет:
// поля непересекающиеся, чтение результата идёт после Wait.
func (s *CompanyService) GetCompany(ctx context.Context, ticker string) (company.Company, error) {
	if ticker == "" {
		return company.Company{}, ErrTickerEmpty
	}

	var result company.Company

	p := pool.New().
		WithErrors().
		WithFirstError().
		WithContext(ctx).
		WithCancelOnError()
	p.Go(func(ctx context.Context) error {
		identity, err := s.identities.FindByTicker(ctx, ticker)
		if err != nil {
			return fmt.Errorf("identity: %w", err)
		}
		result.Identity = identity
		return nil
	})
	p.Go(func(ctx context.Context) error {
		classification, err := s.classifications.FindByTicker(ctx, ticker)
		if err != nil {
			return fmt.Errorf("classification: %w", err)
		}
		result.Classification = classification
		return nil
	})

	if err := p.Wait(); err != nil {
		return company.Company{}, fmt.Errorf("get company %q: %w", ticker, err)
	}

	return result, nil
}
