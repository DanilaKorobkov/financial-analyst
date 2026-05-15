package company

import "time"

// StockDividend — одна дивидендная выплата эмитента (историческая или
// будущая). Нулевое значение поля = «не отдано источником».
type StockDividend struct {
	// LastBuyDate — последний день покупки бумаги с правом на дивиденд.
	LastBuyDate time.Time

	// ReestrCloseDate — дата закрытия реестра акционеров (отсечка).
	ReestrCloseDate time.Time

	// ChangedAt — момент последнего обновления записи.
	ChangedAt time.Time

	// Link — ссылка на первоисточник раскрытия.
	Link string

	// LastBuyPrice — цена закрытия бумаги на LastBuyDate. Нулевое значение —
	// будущая выплата без зафиксированной цены.
	LastBuyPrice float64

	// DivAmount — размер дивиденда на одну акцию в Currency.
	DivAmount float64

	// DivPercent — дивидендная доходность к цене на дату фиксации (%).
	DivPercent float64

	// Year — финансовый год, по которому платится дивиденд.
	Year int64

	// Currency — валюта выплаты.
	Currency Currency

	// Type — тип выплаты (годовая, полугодовая, квартальная, особая).
	Type DividendType
}
