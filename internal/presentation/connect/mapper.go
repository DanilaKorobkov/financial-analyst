// Package connect — Connect-handler для CompanyService.
package connect

import (
	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
)

// toProtoCompany переводит company.Company в proto-сообщение.
func toProtoCompany(c *company.Company) *companyv1.Company {
	return &companyv1.Company{
		Ticker:              c.Ticker,
		Isin:                c.ISIN,
		Name:                c.Name,
		SecurityType:        toProtoSecurityType(c.SecurityType),
		ListingLevel:        toProtoListingLevel(c.ListingLevel),
		Exchange:            toProtoExchange(c.Exchange),
		Currency:            toProtoCurrency(c.Currency),
		Sector:              c.Sector,
		Industry:            c.Industry,
		Country:             c.Country,
		PrimaryReportTicker: c.PrimaryReportTicker,
	}
}

// toProtoSecurityType переводит domain-enum в proto-enum.
func toProtoSecurityType(t company.SecurityType) companyv1.SecurityType {
	switch t {
	case company.SecurityTypeCommonShare:
		return companyv1.SecurityType_SECURITY_TYPE_COMMON_SHARE
	case company.SecurityTypePreferredShare:
		return companyv1.SecurityType_SECURITY_TYPE_PREFERRED_SHARE
	case company.SecurityTypeDepositaryReceipt:
		return companyv1.SecurityType_SECURITY_TYPE_DEPOSITARY_RECEIPT
	default:
		return companyv1.SecurityType_SECURITY_TYPE_UNSPECIFIED
	}
}

// toProtoListingLevel переводит domain-enum в proto-enum.
func toProtoListingLevel(level company.ListingLevel) companyv1.ListingLevel {
	switch level {
	case company.ListingLevelFirst:
		return companyv1.ListingLevel_LISTING_LEVEL_FIRST
	case company.ListingLevelSecond:
		return companyv1.ListingLevel_LISTING_LEVEL_SECOND
	case company.ListingLevelThird:
		return companyv1.ListingLevel_LISTING_LEVEL_THIRD
	default:
		return companyv1.ListingLevel_LISTING_LEVEL_UNSPECIFIED
	}
}

// toProtoExchange переводит domain-enum в proto-enum.
func toProtoExchange(exchange company.Exchange) companyv1.Exchange {
	switch exchange {
	case company.ExchangeMOEX:
		return companyv1.Exchange_EXCHANGE_MOEX
	default:
		return companyv1.Exchange_EXCHANGE_UNSPECIFIED
	}
}

// toProtoCurrency переводит domain-enum в proto-enum.
func toProtoCurrency(currency company.Currency) companyv1.Currency {
	switch currency {
	case company.CurrencyRUB:
		return companyv1.Currency_CURRENCY_RUB
	case company.CurrencyUSD:
		return companyv1.Currency_CURRENCY_USD
	case company.CurrencyEUR:
		return companyv1.Currency_CURRENCY_EUR
	default:
		return companyv1.Currency_CURRENCY_UNSPECIFIED
	}
}
