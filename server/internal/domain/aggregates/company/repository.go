package company

import "context"

// Repository — порт сборки агрегата Company по тикеру. Реализация
// собирает агрегат из источников секций (SecurityDescriptionSource,
// StockInfoSource, …) и возвращает его целиком.
type Repository interface {
	// FindByTicker возвращает агрегат по тикеру. Тикер передаётся источникам
	// как есть, без нормализации.
	// Возвращает ErrNotFound, если хотя бы один источник секции не нашёл
	// данных по тикеру — агрегат отдаётся только целиком.
	FindByTicker(ctx context.Context, ticker string) (Company, error)
}
