package entities

import "context"

// CompanyRepository — порт доступа к справочнику компаний.
type CompanyRepository interface {
	// FindByTicker возвращает справочную карточку компании по тикеру.
	// Если по тикеру нет данных — возвращает ErrCompanyNotFound.
	FindByTicker(ctx context.Context, ticker string) (Company, error)
}

// CompanyMetricsRepository — порт доступа к расширенной карточке эмитента
// (классификация, сводные метрики, дивиденды, темпы роста, консенсусы).
type CompanyMetricsRepository interface {
	// FindByTicker возвращает расширенную карточку эмитента по тикеру.
	//
	// Возвращает ErrNotFound, если эмитента с таким тикером нет.
	// Возвращает ErrUnauthorized, если источник отверг запрос по доступу
	// (токен отсутствует/невалиден).
	// Возвращает ErrQuotaExceeded, если у источника исчерпана квота.
	FindByTicker(ctx context.Context, ticker string) (CompanyMetrics, error)
}
