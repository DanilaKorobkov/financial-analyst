package companycard

const (
	// CurrencyUnspecified — валюта не определена или неподдерживаемая.
	CurrencyUnspecified Currency = 0
	// CurrencyRUB — российский рубль (ISO 4217: RUB).
	CurrencyRUB Currency = 1
	// CurrencyUSD — доллар США (ISO 4217: USD).
	CurrencyUSD Currency = 2
	// CurrencyEUR — евро (ISO 4217: EUR).
	CurrencyEUR Currency = 3
)

// Currency — валюта торгов бумагой.
type Currency int
