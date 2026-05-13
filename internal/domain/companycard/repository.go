package companycard

import "context"

// Repository — порт доступа к коллекции карточек эмитентов.
type Repository interface {
	// FindByTicker возвращает карточку эмитента по бирже и тикеру.
	// Если по паре нет данных — возвращает ErrNotFound.
	FindByTicker(ctx context.Context, exchange Exchange, ticker string) (Card, error)
}
