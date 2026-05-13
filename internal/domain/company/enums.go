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
)

// SecurityType — тип бумаги.
type SecurityType int

// ListingLevel — котировальный уровень бумаги.
type ListingLevel int
