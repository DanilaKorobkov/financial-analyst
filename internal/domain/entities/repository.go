package entities

import "context"

// CompanyRepository — порт доступа к справочнику компаний.
type CompanyRepository interface {
	// FindByTicker возвращает справочную карточку компании по тикеру.
	// Если по тикеру нет данных — возвращает ErrCompanyNotFound.
	FindByTicker(ctx context.Context, ticker string) (Company, error)
}
