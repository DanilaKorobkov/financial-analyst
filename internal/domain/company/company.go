// Package company — агрегат «компания» (эмитент бумаги): сущность,
// её секции (Identity, Classification, ...), порты-gateways к внешним
// источникам, доменные ошибки и enum-ы.
//
// Company — единственный агрегат пакета. Секции — Identity,
// Classification и далее метрики, мультипликаторы, отчётность — это
// value objects: самостоятельной жизни без эмитента не имеют. Каждый
// gateway отвечает за свою секцию; собирает их вместе domain-сервис.
package company

// Company — карточка эмитента. Identity и Classification встроены,
// чтобы сохранить плоский доступ снаружи (c.Ticker, c.Sector) при
// явном разделении на секции, каждую из которых отдаёт свой gateway.
type Company struct {
	Identity
	Classification
}

// Identity — идентификационная секция: то, что однозначно определяет
// бумагу.
type Identity struct {
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

// Classification — классификационная секция: биржа, валюта, отраслевая
// принадлежность, страна и привязка отчётности.
//
// Порядок полей продиктован fieldalignment: pointer-содержащие
// (string) идут до non-pointer (enum-int), чтобы GC сканировал
// меньший диапазон.
type Classification struct {
	// Sector — название сектора эмитента. Пример: "Финансы".
	Sector string

	// Industry — название отрасли. Пример: "Банковская деятельность".
	Industry string

	// Country — страна регистрации эмитента. Пример: "Россия".
	Country string

	// PrimaryReportTicker — тикер «основной» бумаги эмитента, к которой
	// привязана отчётность. Для привилегированных акций совпадает с
	// обыкновенными. Пример: "SBER".
	PrimaryReportTicker string

	// Exchange — биржа листинга. Пример: ExchangeMOEX.
	Exchange Exchange

	// Currency — валюта торгов бумагой. Пример: CurrencyRUB.
	Currency Currency
}
