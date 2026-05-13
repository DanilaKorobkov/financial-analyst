// Package services — domain-сервисы, оркеструют агрегаты и порты.
package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
)

// ErrTickerEmpty — клиент передал пустой тикер.
var ErrTickerEmpty = errors.New("ticker is empty")

// CompanyInfo — сервис поиска справочной карточки компании по тикеру.
type CompanyInfo struct {
	companies company.Repository
}

// NewCompanyInfo собирает сервис вокруг репозитория компаний.
func NewCompanyInfo(companies company.Repository) *CompanyInfo {
	return &CompanyInfo{companies: companies}
}

// Lookup проверяет непустоту тикера и делегирует поиск репозиторию.
// Тикер передаётся как есть, без нормализации.
func (s *CompanyInfo) Lookup(ctx context.Context, ticker string) (company.Company, error) {
	if ticker == "" {
		return company.Company{}, ErrTickerEmpty
	}
	found, err := s.companies.FindByTicker(ctx, ticker)
	if err != nil {
		return company.Company{}, fmt.Errorf("lookup ticker %q: %w", ticker, err)
	}
	return found, nil
}
