// Package services — domain-сервисы, оркеструют агрегаты и порты.
package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

// ErrTickerEmpty — клиент передал пустой тикер.
var ErrTickerEmpty = errors.New("ticker is empty")

// CompanyService собирает агрегат «компания» по тикеру. Сам сервис
// не знает ни про источники, ни про их количество — сборкой управляет
// company.Repository.
type CompanyService struct {
	companies company.Repository
}

// ConfigCompanyService — параметры CompanyService.
type ConfigCompanyService struct {
	// Companies — репозиторий, собирающий агрегат Company по тикеру.
	Companies company.Repository
}

// NewCompanyService собирает сервис вокруг репозитория компаний.
func NewCompanyService(cfg ConfigCompanyService) *CompanyService {
	return &CompanyService{companies: cfg.Companies}
}

// GetCompany проверяет непустоту тикера и просит репозиторий собрать
// агрегат. Тикер передаётся как есть, без нормализации.
//
// Возможные ошибки:
//   - ErrTickerEmpty — пустой тикер;
//   - company.ErrNotFound — хотя бы один источник секции не нашёл
//     бумагу по тикеру;
//   - произвольная ошибка репозитория — пробрасывается с пометкой тикера.
func (s *CompanyService) GetCompany(ctx context.Context, ticker string) (company.Company, error) {
	if ticker == "" {
		return company.Company{}, ErrTickerEmpty
	}

	got, err := s.companies.FindByTicker(ctx, ticker)
	if err != nil {
		return company.Company{}, fmt.Errorf("get company %q: %w", ticker, err)
	}
	return got, nil
}
