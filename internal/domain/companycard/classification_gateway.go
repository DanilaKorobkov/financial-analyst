package companycard

import "context"

// ClassificationGateway — порт доступа к классификационной секции
// карточки эмитента (Classification) во внешнем источнике-справочнике.
type ClassificationGateway interface {
	// FindByTicker возвращает классификационную секцию карточки по тикеру.
	// Тикер передаётся как есть, без нормализации.
	// Если по тикеру нет данных — возвращает ErrNotFound.
	FindByTicker(ctx context.Context, ticker string) (Classification, error)
}
