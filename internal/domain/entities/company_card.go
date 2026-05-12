package entities

// CompanyCard — лёгкая карточка эмитента: идентификация и классификация
// (страна, валюта, сектор/отрасль по GICS-подобной иерархии).
//
// Используется как самостоятельный элемент списочных ответов и как блок-шапка
// расширенной карточки CompanyMetrics.
type CompanyCard struct {
	// Ticker — биржевой код бумаги. Пример: "SBER".
	Ticker string

	// Exchange — код биржи листинга. Пример: "MOEX".
	Exchange string

	// Name — короткое название бумаги. Пример: "Сбербанк".
	Name string

	// Sector — название сектора эмитента (русское). Пример: "Финансы".
	Sector string

	// Industry — название отрасли (русское).
	Industry string

	// IndustryGroup — название группы отраслей (русское).
	IndustryGroup string

	// Country — страна регистрации эмитента (русское название).
	Country string

	// Currency — валюта торгов бумагой в формате ISO 4217. Пример: "RUB".
	Currency string

	// PrimaryReportTicker — тикер «основной» бумаги эмитента, к которой
	// привязана отчётность. Для привилегированных акций совпадает с обыкновенными.
	PrimaryReportTicker string

	// PrimaryReportExchange — биржа основной бумаги (см. PrimaryReportTicker).
	PrimaryReportExchange string

	// SectorID — числовой код сектора в классификаторе источника.
	SectorID int

	// IndustryID — числовой код отрасли.
	IndustryID int

	// IndustryGroupID — числовой код группы отраслей.
	IndustryGroupID int
}
