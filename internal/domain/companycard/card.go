// Package companycard — агрегат «карточка эмитента»: сущность,
// её секции (Identity, Classification), порты-gateways к внешним
// источникам, доменные ошибки и enum-ы.
//
// Card — единственный агрегат пакета. Identity и Classification —
// value objects, секции карточки; самостоятельной жизни без эмитента
// не имеют. Card собирается domain-сервисом из ответов двух gateway,
// каждый из которых отвечает за свою секцию.
package companycard

// Card — карточка эмитента. Identity и Classification встроены, чтобы
// сохранить плоский доступ снаружи (card.Ticker, card.Sector) при
// явном разделении на секции, каждую из которых отдаёт свой gateway.
type Card struct {
	Identity
	Classification
}

// Identity — идентификационная секция карточки: то, что однозначно
// определяет бумагу.
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

// Classification — классификационная секция карточки эмитента:
// биржа, валюта, отраслевая принадлежность, страна и привязка
// отчётности.
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
