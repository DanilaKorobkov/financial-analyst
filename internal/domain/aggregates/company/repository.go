package company

import "context"

// Repository — порт сборки агрегата Company по тикеру. Реализация
// оркеструет источники секций (SecurityDescriptionSource, StockInfoSource, …)
// и складывает результат в Company.
type Repository interface {
	// FindByTicker возвращает агрегат по тикеру. Тикер передаётся источникам
	// как есть, без нормализации.
	// Возвращает ErrCompanyNotFound, если ни один источник секции не нашёл
	// данных по тикеру.
	FindByTicker(ctx context.Context, ticker string) (*Company, error)
}
