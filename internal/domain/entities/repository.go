package entities

import "context"

// CompanyRepository — порт доступа к справочнику компаний.
type CompanyRepository interface {
	// FindByTicker возвращает справочную карточку компании по тикеру.
	// Если по тикеру нет данных — возвращает ErrCompanyNotFound.
	FindByTicker(ctx context.Context, ticker string) (Company, error)
}

// CompanyCardRepository — порт доступа к карточке эмитента: идентификация,
// классификация (сектор / отрасль / страна / валюта), описание и ссылки.
type CompanyCardRepository interface {
	// FindByTicker возвращает карточку эмитента по тикеру.
	//
	// Возвращает ErrNotFound, если эмитента с таким тикером нет.
	// Возвращает ErrUnauthorized, если источник отверг запрос по доступу
	// (токен отсутствует или невалиден).
	// Возвращает ErrQuotaExceeded, если у источника исчерпана квота.
	FindByTicker(ctx context.Context, ticker string) (CompanyCard, error)
}
