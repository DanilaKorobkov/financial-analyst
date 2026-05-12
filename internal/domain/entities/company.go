// Package entities — плоские domain-сущности (без подпапок на агрегат).
package entities

const (
	// SecurityTypeUnspecified — неизвестный или неподдерживаемый тип.
	SecurityTypeUnspecified SecurityType = iota
	// SecurityTypeCommonShare — обыкновенная акция.
	SecurityTypeCommonShare
	// SecurityTypePreferredShare — привилегированная акция.
	SecurityTypePreferredShare
	// SecurityTypeDepositaryReceipt — депозитарная расписка.
	SecurityTypeDepositaryReceipt
)

//nolint:grouper // отдельный блок для независимого iota-enum (счётчик iota не сбрасывается внутри одного блока).
const (
	// ListingLevelUnspecified — уровень не указан.
	ListingLevelUnspecified ListingLevel = iota
	// ListingLevelFirst — первый котировальный уровень.
	ListingLevelFirst
	// ListingLevelSecond — второй котировальный уровень.
	ListingLevelSecond
	// ListingLevelThird — третий котировальный уровень.
	ListingLevelThird
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
