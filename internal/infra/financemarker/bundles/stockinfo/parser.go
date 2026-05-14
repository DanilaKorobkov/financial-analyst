package stockinfo

import (
	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
)

var (
	// exchangeByCode переводит строковый код биржи FinanceMarker в domain-enum.
	// Незнакомые значения возвращают zero (ExchangeUnspecified) — набор бирж
	// со временем расширяется.
	exchangeByCode = map[string]domaincompany.Exchange{
		codeExchangeMOEX: domaincompany.ExchangeMOEX,
	}

	// currencyByCode переводит ISO 4217-код валюты FinanceMarker в domain-enum.
	// Незнакомые значения возвращают zero (CurrencyUnspecified).
	currencyByCode = map[string]domaincompany.Currency{
		"RUB": domaincompany.CurrencyRUB,
		"USD": domaincompany.CurrencyUSD,
		"EUR": domaincompany.CurrencyEUR,
	}

	// reportFrequencyByCode переводит код report_frequency FinanceMarker
	// ("Y" / "Q") в domain-enum. Незнакомые значения возвращают zero
	// (ReportFrequencyUnspecified).
	reportFrequencyByCode = map[string]domaincompany.ReportFrequency{
		"Y": domaincompany.ReportFrequencyYearly,
		"Q": domaincompany.ReportFrequencyQuarterly,
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
// добавляются по мере появления соответствующих bundles.
type stockDTO struct {
	Info infoDTO `json:"info"`
}

// translateStockInfo раскладывает info-блок FinanceMarker по каноничным
// полям FieldValues.
func translateStockInfo(info *infoDTO) data.FieldValues {
	return data.FieldValues{
		domaincompany.FieldIssuerName:            info.Name,
		domaincompany.FieldCountry:               info.Country,
		domaincompany.FieldSector:                info.Sector,
		domaincompany.FieldIndustryGroup:         info.IndustryGroup,
		domaincompany.FieldIndustry:              info.Industry,
		domaincompany.FieldSubIndustry:           info.SubIndustry,
		domaincompany.FieldDescription:           info.Description,
		domaincompany.FieldSite:                  info.Site,
		domaincompany.FieldDisclosureLink:        info.DiscLink,
		domaincompany.FieldPrimaryReportTicker:   info.PrimaryReportCode,
		domaincompany.FieldSectorID:              info.SectorID,
		domaincompany.FieldIndustryGroupID:       info.IndustryGroupID,
		domaincompany.FieldIndustryID:            info.IndustryID,
		domaincompany.FieldSubIndustryID:         info.SubIndustryID,
		domaincompany.FieldExchange:              exchangeByCode[info.Exchange],
		domaincompany.FieldPrimaryReportExchange: exchangeByCode[info.PrimaryReportExchange],
		domaincompany.FieldCurrency:              currencyByCode[info.Currency],
		domaincompany.FieldReportFrequency:       reportFrequencyByCode[info.ReportFrequency],
		domaincompany.FieldSPB:                   info.SPB,
	}
}
