// Package services — domain-сервисы, оркеструют агрегаты и порты.
package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/sourcegraph/conc/pool"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/companycard"
)

// ErrTickerEmpty — клиент передал пустой тикер.
var ErrTickerEmpty = errors.New("ticker is empty")

// CompanyCardService — сервис сборки карточки эмитента из двух секций.
// Каждая секция приходит из своего внешнего источника-gateway; обе
// обязательные, частичная карточка не возвращается.
type CompanyCardService struct {
	identities      companycard.IdentityGateway
	classifications companycard.ClassificationGateway
}

// NewCompanyCardService собирает сервис вокруг двух gateway.
func NewCompanyCardService(
	identities companycard.IdentityGateway,
	classifications companycard.ClassificationGateway,
) *CompanyCardService {
	return &CompanyCardService{
		identities:      identities,
		classifications: classifications,
	}
}

// GetCard проверяет непустоту тикера и параллельно тянет обе секции.
// При ошибке любого gateway возвращает её, не дожидаясь второго.
// Тикер передаётся как есть, без нормализации.
func (s *CompanyCardService) GetCard(ctx context.Context, ticker string) (companycard.Card, error) {
	if ticker == "" {
		return companycard.Card{}, ErrTickerEmpty
	}

	var (
		identity       companycard.Identity
		classification companycard.Classification
	)

	p := pool.New().WithErrors().WithContext(ctx)
	p.Go(func(ctx context.Context) error {
		found, err := s.identities.FindByTicker(ctx, ticker)
		if err != nil {
			return fmt.Errorf("identity: %w", err)
		}
		identity = found
		return nil
	})
	p.Go(func(ctx context.Context) error {
		found, err := s.classifications.FindByTicker(ctx, ticker)
		if err != nil {
			return fmt.Errorf("classification: %w", err)
		}
		classification = found
		return nil
	})

	if err := p.Wait(); err != nil {
		return companycard.Card{}, fmt.Errorf("get card %q: %w", ticker, err)
	}

	return companycard.Card{Identity: identity, Classification: classification}, nil
}
