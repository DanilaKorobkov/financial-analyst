package stockinfo

import (
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

var (
	// exchangeByCode переводит строковый код биржи FinanceMarker в domain-enum.
	// Незнакомые значения возвращают zero (ExchangeUnspecified) — набор бирж
	// со временем расширяется.
	exchangeByCode = map[string]company.Exchange{
		codeExchangeMOEX: company.ExchangeMOEX,
	}

	// currencyByCode переводит ISO 4217-код валюты FinanceMarker в domain-enum.
	// Незнакомые значения возвращают zero (CurrencyUnspecified).
	currencyByCode = map[string]company.Currency{
		"RUB": company.CurrencyRUB,
		"USD": company.CurrencyUSD,
		"EUR": company.CurrencyEUR,
	}

	// reportFrequencyByCode переводит код report_frequency FinanceMarker
	// ("Y" / "Q") в domain-enum. Незнакомые значения возвращают zero
	// (ReportFrequencyUnspecified).
	reportFrequencyByCode = map[string]company.ReportFrequency{
		"Y": company.ReportFrequencyYearly,
		"Q": company.ReportFrequencyQuarterly,
	}
)

// infoDTO — блок `info` ответа /api/fm/v2/stocks/{exchange}:{code}.
// Содержит полный набор полей, отдаваемых источником.
type infoDTO struct {
	Code                  string `json:"code"`
	Name                  string `json:"name"`
	Exchange              string `json:"exchange"`
	Country               string `json:"country"`
	Currency              string `json:"currency"`
	Sector                string `json:"sector"`
	IndustryGroup         string `json:"industry_group"`
	Industry              string `json:"industry"`
	SubIndustry           string `json:"sub_industry"`
	PrimaryReportCode     string `json:"primary_report_code"`
	PrimaryReportExchange string `json:"primary_report_exchange"`
	Description           string `json:"description"`
	Site                  string `json:"site"`
	DiscLink              string `json:"disc_link"`
	ReportFrequency       string `json:"report_frequency"`
	SectorID              int64  `json:"sector_id"`
	IndustryGroupID       int64  `json:"industry_group_id"`
	IndustryID            int64  `json:"industry_id"`
	SubIndustryID         int64  `json:"sub_industry_id"`
	SPB                   bool   `json:"spb"`
}

// stockDTO — корневой объект ответа эндпоинта по эмитенту. Здесь
// разбирается только блок info — остальные разделы (summary / ratios / ...)
// добавляются по мере появления соответствующих источников.
type stockDTO struct {
	Info infoDTO `json:"info"`
}

// translateStockInfo раскладывает info-блок FinanceMarker в StockInfo.
func translateStockInfo(info *infoDTO) *company.StockInfo {
	return &company.StockInfo{
		IssuerName:            info.Name,
		Sector:                info.Sector,
		IndustryGroup:         info.IndustryGroup,
		Industry:              info.Industry,
		SubIndustry:           info.SubIndustry,
		Country:               info.Country,
		Description:           info.Description,
		Site:                  info.Site,
		DisclosureLink:        info.DiscLink,
		PrimaryReportTicker:   info.PrimaryReportCode,
		SectorID:              info.SectorID,
		IndustryGroupID:       info.IndustryGroupID,
		IndustryID:            info.IndustryID,
		SubIndustryID:         info.SubIndustryID,
		PrimaryReportExchange: exchangeByCode[info.PrimaryReportExchange],
		Exchange:              exchangeByCode[info.Exchange],
		Currency:              currencyByCode[info.Currency],
		ReportFrequency:       reportFrequencyByCode[info.ReportFrequency],
		SPB:                   info.SPB,
	}
}
