package company

import "time"

// StockOwner — доля одного акционера на дату среза. Owner может быть как
// конкретным владельцем, так и категорией (`Прочие`, `Free float`).
type StockOwner struct {
	// ChangedAt — момент обновления записи в FM.
	ChangedAt time.Time

	// Owner — имя/название акционера или категории.
	Owner string

	// Link — ссылка на источник раскрытия структуры владения.
	Link string

	// Period — дата среза владения.
	Period StockPeriod

	// Own — доля владения в процентах от уставного капитала.
	Own float64
}
