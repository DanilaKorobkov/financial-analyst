package company

import "time"

// StockInsiderTransaction — сделка инсайдера эмитента в узком формате,
// возвращаемом блоком `insiderTransactions` ответа `stock/`. Полный
// формат сделок с объёмом и ценой публикуется отдельной лентой
// инсайдерских сделок и в этот агрегат не входит.
type StockInsiderTransaction struct {
	// TransactionDate — дата фактической сделки.
	TransactionDate time.Time

	// Insider — имя/название инсайдера. Может содержать строку
	// «Информация не раскрывается», если данные скрыты по требованию ЦБ.
	Insider string

	// InsiderTitle — должность инсайдера (для физлиц-менеджеров).
	InsiderTitle string

	// Type — направление сделки (покупка/продажа).
	Type InsiderTransactionType
}
