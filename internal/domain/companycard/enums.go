package companycard

// Численные значения зафиксированы явно: enum-ы выходят наружу через
// proto-контракт presentation-слоя, и `iota` сделал бы порядок строк
// load-bearing — вставка члена в середину сдвинула бы коды.
const (
	// ExchangeUnspecified — биржа не определена или неподдерживаемая.
	ExchangeUnspecified Exchange = 0
	// ExchangeMOEX — Московская биржа.
	ExchangeMOEX Exchange = 1

	// CurrencyUnspecified — валюта не определена или неподдерживаемая.
	CurrencyUnspecified Currency = 0
	// CurrencyRUB — российский рубль (ISO 4217: RUB).
	CurrencyRUB Currency = 1
	// CurrencyUSD — доллар США (ISO 4217: USD).
	CurrencyUSD Currency = 2
	// CurrencyEUR — евро (ISO 4217: EUR).
	CurrencyEUR Currency = 3
)

// Exchange — биржа листинга бумаги.
type Exchange int

// Currency — валюта торгов бумагой.
type Currency int
