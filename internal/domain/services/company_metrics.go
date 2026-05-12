package services

import (
	"context"
	"fmt"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// CompanyMetrics — сервис поиска расширенной карточки эмитента.
type CompanyMetrics struct {
	metrics entities.CompanyMetricsRepository
}

// NewCompanyMetrics собирает сервис вокруг репозитория расширенных карточек.
func NewCompanyMetrics(metrics entities.CompanyMetricsRepository) *CompanyMetrics {
	return &CompanyMetrics{metrics: metrics}
}

// FindByTicker проверяет непустоту тикера и делегирует поиск репозиторию.
// Тикер передаётся как есть, без нормализации.
func (s *CompanyMetrics) FindByTicker(
	ctx context.Context,
	ticker string,
) (entities.CompanyMetrics, error) {
	if ticker == "" {
		return entities.CompanyMetrics{}, ErrTickerEmpty
	}
	metrics, err := s.metrics.FindByTicker(ctx, ticker)
	if err != nil {
		return entities.CompanyMetrics{}, fmt.Errorf("find metrics %q: %w", ticker, err)
	}
	return metrics, nil
}
