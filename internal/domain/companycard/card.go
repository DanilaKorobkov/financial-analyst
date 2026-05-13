// Package companycard — агрегат «карточка эмитента»: сущность,
// порт доступа к коллекции и доменные ошибки (вместе с enum-ами
// Exchange и Currency, специфичными для карточки).
//
// Используется как самостоятельный элемент списочных ответов и как блок-шапка
// будущих расширенных карточек (сводные метрики, мультипликаторы, отчётность).
package companycard

// Card — карточка эмитента.
type Card struct {
	// Ticker — биржевой код бумаги. Пример: "SBER".
	Ticker string

	// Name — короткое название бумаги. Пример: "Сбербанк".
	Name string

	// Sector — название сектора эмитента. Пример: "Финансы".
	Sector string

	// Industry — название отрасли. Пример: "Банковская деятельность".
	Industry string

	// IndustryGroup — название группы отраслей. Пример: "Банковская деятельность".
	IndustryGroup string

	// Country — страна регистрации эмитента. Пример: "Россия".
	Country string

	// PrimaryReportTicker — тикер «основной» бумаги эмитента, к которой
	// привязана отчётность. Для привилегированных акций совпадает с обыкновенными.
	// Пример: "SBER".
	PrimaryReportTicker string

	// Description — описание эмитента в свободной форме.
	// Пример: "ПАО «Сбербанк» — крупнейший универсальный банк России.".
	Description string

	// Site — корпоративный сайт эмитента. Пример: "https://www.sberbank.com".
	Site string

	// DiscLink — ссылка на страницу раскрытия информации (investor relations).
	// Пример: "https://www.sberbank.com/ru/investor-relations".
	DiscLink string

	// SectorID — числовой код сектора в классификаторе источника. Пример: 40.
	SectorID int

	// IndustryID — числовой код отрасли. Пример: 401010.
	IndustryID int

	// IndustryGroupID — числовой код группы отраслей. Пример: 4010.
	IndustryGroupID int

	// Exchange — биржа листинга. Пример: ExchangeMOEX.
	Exchange Exchange

	// Currency — валюта торгов бумагой. Пример: CurrencyRUB.
	Currency Currency

	// PrimaryReportExchange — биржа основной бумаги (см. PrimaryReportTicker).
	// Пример: ExchangeMOEX.
	PrimaryReportExchange Exchange
}
