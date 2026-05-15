package company

// StockInfo — карточка эмитента: классификация GICS, страна, описание,
// ссылки, основная отчётная бумага, биржа и частота отчётности.
type StockInfo struct {
	// IssuerName — название эмитента у справочника.
	IssuerName string

	// Sector — сектор эмитента (GICS).
	Sector string

	// IndustryGroup — группа отраслей (GICS).
	IndustryGroup string

	// Industry — отрасль эмитента (GICS).
	Industry string

	// SubIndustry — под-отрасль эмитента (GICS).
	SubIndustry string

	// Country — страна регистрации эмитента.
	Country string

	// Description — текстовое описание эмитента.
	Description string

	// Site — корпоративный сайт эмитента.
	Site string

	// DisclosureLink — ссылка на страницу раскрытия эмитента.
	DisclosureLink string

	// PrimaryReportTicker — тикер основной бумаги, к которой привязана отчётность.
	PrimaryReportTicker string

	// SectorID — числовой код сектора GICS.
	SectorID int64

	// IndustryGroupID — числовой код группы отраслей GICS.
	IndustryGroupID int64

	// IndustryID — числовой код отрасли GICS.
	IndustryID int64

	// SubIndustryID — числовой код под-отрасли GICS.
	SubIndustryID int64

	// PrimaryReportExchange — биржа основной отчётной бумаги.
	PrimaryReportExchange Exchange

	// Exchange — биржа листинга бумаги.
	Exchange Exchange

	// Currency — валюта торгов бумагой.
	Currency Currency

	// ReportFrequency — частота публикации отчётности эмитентом.
	ReportFrequency ReportFrequency

	// SPB — дополнительный листинг на СПБ-бирже.
	SPB bool
}
