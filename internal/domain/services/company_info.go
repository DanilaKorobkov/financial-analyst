// Package services — domain-сервисы, оркеструют entities и порты.
package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// ErrTickerEmpty — клиент передал пустой тикер.
var ErrTickerEmpty = errors.New("ticker is empty")

// CompanyInfo — сервис поиска справочной карточки компании по тикеру.
type CompanyInfo struct {
	companies entities.CompanyRepository
}

// NewCompanyInfo собирает сервис вокруг репозитория компаний.
func NewCompanyInfo(companies entities.CompanyRepository) *CompanyInfo {
	return &CompanyInfo{companies: companies}
}

// Lookup проверяет непустоту тикера и делегирует поиск репозиторию.
// Тикер передаётся как есть, без нормализации.
func (s *CompanyInfo) Lookup(ctx context.Context, ticker string) (entities.Company, error) {
	if ticker == "" {
		return entities.Company{}, ErrTickerEmpty
	}
	company, err := s.companies.FindByTicker(ctx, ticker)
	if err != nil {
		return entities.Company{}, fmt.Errorf("lookup ticker %q: %w", ticker, err)
	}
	return company, nil
}
