// Package entities — плоские domain-сущности (без подпапок на агрегат).
package entities

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

// Company — справочная карточка эмитента.
type Company struct {
	// Ticker — биржевой код бумаги. Пример: "SBER".
	Ticker string

	// ISIN — международный идентификатор. Пример: "RU0009029540".
	ISIN string

	// Name — полное название бумаги. Пример: "Сбербанк России ПАО ао".
	Name string

	// SecurityType — тип бумаги. SecurityTypeUnspecified — неизвестный
	// или неподдерживаемый тип.
	SecurityType SecurityType

	// ListingLevel — котировальный уровень. ListingLevelUnspecified —
	// уровень не указан.
	ListingLevel ListingLevel
}

// SecurityType — тип бумаги.
type SecurityType int

// ListingLevel — котировальный уровень бумаги.
type ListingLevel int
