package company

import "context"

// Repository — порт доступа к коллекции справочных карточек эмитентов.
type Repository interface {
	// FindByTicker возвращает справочную карточку компании по тикеру.
	// Если по тикеру нет данных — возвращает ErrNotFound.
	FindByTicker(ctx context.Context, ticker string) (Company, error)
}
