// Package services — domain-сервисы, оркеструют entities и порты.
package services

import (
	"context"
	"errors"

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
// Тикер передаётся в источник как есть — все источники (MOEX ISS,
// FinanceMarker) регистронезависимы, своей нормализации не делаем.
func (s *CompanyInfo) Lookup(ctx context.Context, ticker string) (entities.Company, error) {
	if ticker == "" {
		return entities.Company{}, ErrTickerEmpty
	}
	return s.companies.FindByTicker(ctx, ticker)
}
