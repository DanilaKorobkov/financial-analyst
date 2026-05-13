package company

// Численные значения зафиксированы явно: enum-ы выходят наружу через
// proto-контракт presentation-слоя, и `iota` сделал бы порядок строк
// load-bearing — вставка члена в середину сдвинула бы коды.
const (
	// SecurityTypeUnspecified — неизвестный или неподдерживаемый тип.
	SecurityTypeUnspecified SecurityType = 0
	// SecurityTypeCommonShare — обыкновенная акция.
	SecurityTypeCommonShare SecurityType = 1
	// SecurityTypePreferredShare — привилегированная акция.
	SecurityTypePreferredShare SecurityType = 2
	// SecurityTypeDepositaryReceipt — депозитарная расписка.
	SecurityTypeDepositaryReceipt SecurityType = 3

	// ListingLevelUnspecified — уровень не указан.
	ListingLevelUnspecified ListingLevel = 0
	// ListingLevelFirst — первый котировальный уровень.
	ListingLevelFirst ListingLevel = 1
	// ListingLevelSecond — второй котировальный уровень.
	ListingLevelSecond ListingLevel = 2
	// ListingLevelThird — третий котировальный уровень.
	ListingLevelThird ListingLevel = 3

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

// SecurityType — тип бумаги.
type SecurityType int

// ListingLevel — котировальный уровень бумаги.
type ListingLevel int

// Exchange — биржа листинга бумаги.
type Exchange int

// Currency — валюта торгов бумагой.
type Currency int
