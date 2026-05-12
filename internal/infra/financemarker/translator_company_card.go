package financemarker

import "github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"

// infoDTO — блок `info` ответа /api/fm/v2/stocks/{exchange}:{code}.
type infoDTO struct {
	Code                  string `json:"code"`
	Name                  string `json:"name"`
	Exchange              string `json:"exchange"`
	Country               string `json:"country"`
	Currency              string `json:"currency"`
	Sector                string `json:"sector"`
	Industry              string `json:"industry"`
	IndustryGroup         string `json:"industry_group"`
	PrimaryReportCode     string `json:"primary_report_code"`
	PrimaryReportExchange string `json:"primary_report_exchange"`
	Description           string `json:"description"`
	Site                  string `json:"site"`
	DiscLink              string `json:"disc_link"`
	SectorID              int    `json:"sector_id"`
	IndustryID            int    `json:"industry_id"`
	IndustryGroupID       int    `json:"industry_group_id"`
}

// stockDTO — корневой объект ответа эндпоинта по эмитенту. Здесь
// разбирается только блок info — остальные разделы (summary / ratios / ...)
// добавляются по мере появления соответствующих репозиториев.
type stockDTO struct {
	Info infoDTO `json:"info"`
}

// translateCompanyCard собирает entities.CompanyCard из info-блока FinanceMarker.
func translateCompanyCard(info *infoDTO) entities.CompanyCard {
	return entities.CompanyCard{
		Ticker:                info.Code,
		Exchange:              translateExchange(info.Exchange),
		Name:                  info.Name,
		Sector:                info.Sector,
		Industry:              info.Industry,
		IndustryGroup:         info.IndustryGroup,
		Country:               info.Country,
		Currency:              translateCurrency(info.Currency),
		PrimaryReportTicker:   info.PrimaryReportCode,
		PrimaryReportExchange: translateExchange(info.PrimaryReportExchange),
		Description:           info.Description,
		Site:                  info.Site,
		DiscLink:              info.DiscLink,
		SectorID:              info.SectorID,
		IndustryID:            info.IndustryID,
		IndustryGroupID:       info.IndustryGroupID,
	}
}

// translateExchange переводит строковый код биржи FinanceMarker в domain-enum.
// Незнакомые значения становятся ExchangeUnspecified — набор бирж со временем
// расширяется.
func translateExchange(s string) entities.Exchange {
	switch s {
	case codeExchangeMOEX:
		return entities.ExchangeMOEX
	default:
		return entities.ExchangeUnspecified
	}
}

// translateCurrency переводит ISO 4217-код валюты FinanceMarker в domain-enum.
// Незнакомые значения становятся CurrencyUnspecified.
func translateCurrency(s string) entities.Currency {
	switch s {
	case "RUB":
		return entities.CurrencyRUB
	case "USD":
		return entities.CurrencyUSD
	case "EUR":
		return entities.CurrencyEUR
	default:
		return entities.CurrencyUnspecified
	}
}
