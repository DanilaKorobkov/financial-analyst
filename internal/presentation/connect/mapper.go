// Package connect — Connect-handler для CompanyService.
package connect

import (
	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
)

// toProtoCompany переводит company.Company в proto-сообщение.
func toProtoCompany(c *company.Company) *companyv1.Company {
	return &companyv1.Company{
		Ticker:       c.Ticker,
		Isin:         c.ISIN,
		Name:         c.Name,
		SecurityType: toProtoSecurityType(c.SecurityType),
		ListingLevel: toProtoListingLevel(c.ListingLevel),
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
