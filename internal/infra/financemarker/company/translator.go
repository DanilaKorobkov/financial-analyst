package company

import (
	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
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
)

// infoDTO — блок `info` ответа /api/fm/v2/stocks/{exchange}:{code}.
// Сюда вынесены только те поля, которые забирает gateway; остальные
// поля блока (description, site, disc_link, sector_id и т.п.) пока не
// нужны.
type infoDTO struct {
	Exchange          string `json:"exchange"`
	Country           string `json:"country"`
	Currency          string `json:"currency"`
	Sector            string `json:"sector"`
	Industry          string `json:"industry"`
	PrimaryReportCode string `json:"primary_report_code"`
}

// stockDTO — корневой объект ответа эндпоинта по эмитенту. Здесь
// разбирается только блок info — остальные разделы (summary / ratios / ...)
// добавляются по мере появления соответствующих gateway.
type stockDTO struct {
	Info infoDTO `json:"info"`
}

// translateClassification собирает domaincompany.Classification из info-блока FinanceMarker.
func translateClassification(info *infoDTO) domaincompany.Classification {
	return domaincompany.Classification{
		Exchange:            exchangeByCode[info.Exchange],
		Currency:            currencyByCode[info.Currency],
		Sector:              info.Sector,
		Industry:            info.Industry,
		Country:             info.Country,
		PrimaryReportTicker: info.PrimaryReportCode,
	}
}
