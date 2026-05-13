package companycard

// Численные значения зафиксированы явно: enum-ы выходят наружу через
// proto-контракт presentation-слоя, и `iota` сделал бы порядок строк
// load-bearing — вставка члена в середину сдвинула бы коды.
const (
	// ExchangeUnspecified — биржа не определена или неподдерживаемая.
	ExchangeUnspecified Exchange = 0
	// ExchangeMOEX — Московская биржа.
	ExchangeMOEX Exchange = 1
)

// Exchange — биржа листинга бумаги.
type Exchange int
