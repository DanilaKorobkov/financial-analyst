// Package company — агрегат «справочная карточка эмитента»: сущность,
// порт доступа к коллекции и доменные ошибки.
package company

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
