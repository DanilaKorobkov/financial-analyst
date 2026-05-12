// Package entities — плоские domain-сущности (без подпапок на агрегат).
package entities

// Company — справочная карточка эмитента MOEX. Поля из блока description
// эндпоинта /iss/securities/{TICKER}.json. Примеры значений — по SBER.
type Company struct {
	// Ticker — биржевой код бумаги (SECID). Пример: "SBER".
	Ticker string

	// ISIN — международный идентификатор. Пример: "RU0009029540".
	ISIN string

	// Name — полное название бумаги. Пример: "Сбербанк России ПАО ао".
	Name string

	// SecurityType — тип бумаги. SecurityTypeUnspecified — биржа вернула
	// неизвестное значение поля TYPE (см. справочник get_securitygroups MOEX).
	SecurityType SecurityType

	// ListingLevel — котировальный уровень MOEX. ListingLevelUnspecified —
	// биржа не указала уровень (поле LISTLEVEL отсутствовало).
	ListingLevel ListingLevel
}

// SecurityType — тип бумаги MOEX (поле TYPE блока description).
type SecurityType int

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

// ListingLevel — котировальный уровень MOEX.
type ListingLevel int

const (
	// ListingLevelUnspecified — биржа не указала уровень.
	ListingLevelUnspecified ListingLevel = iota
	// ListingLevelFirst — первый котировальный уровень.
	ListingLevelFirst
	// ListingLevelSecond — второй котировальный уровень.
	ListingLevelSecond
	// ListingLevelThird — третий котировальный уровень.
	ListingLevelThird
)
