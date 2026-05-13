package company

import "context"

// ClassificationGateway — порт доступа к классификационной секции
// компании (Classification) во внешнем источнике-справочнике.
type ClassificationGateway interface {
	// FindByTicker возвращает классификационную секцию по тикеру.
	// Тикер передаётся как есть, без нормализации.
	// Если по тикеру нет данных — возвращает ErrNotFound.
	FindByTicker(ctx context.Context, ticker string) (Classification, error)
}
