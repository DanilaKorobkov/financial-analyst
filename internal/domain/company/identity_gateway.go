package company

import "context"

// IdentityGateway — порт доступа к идентификационной секции компании
// (Identity) во внешнем источнике-«реестре бумаг».
type IdentityGateway interface {
	// FindByTicker возвращает идентификационную секцию по тикеру.
	// Тикер передаётся как есть, без нормализации.
	// Если по тикеру нет данных — возвращает ErrNotFound.
	FindByTicker(ctx context.Context, ticker string) (Identity, error)
}
