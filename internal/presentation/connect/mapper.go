// Package connect — Connect-handler для CompanyService.
package connect

import (
	"math"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// toProtoCompany переводит entities.Company в proto-сообщение.
func toProtoCompany(c *entities.Company) *companyv1.Company {
	return &companyv1.Company{
		Ticker:       c.Ticker,
		Isin:         c.ISIN,
		Name:         c.Name,
		SecurityType: toProtoSecurityType(c.SecurityType),
		ListingLevel: toProtoListingLevel(c.ListingLevel),
	}
}

// toProtoSecurityType переводит domain-enum в proto-enum.
func toProtoSecurityType(t entities.SecurityType) companyv1.SecurityType {
	switch t {
	case entities.SecurityTypeCommonShare:
		return companyv1.SecurityType_SECURITY_TYPE_COMMON_SHARE
	case entities.SecurityTypePreferredShare:
		return companyv1.SecurityType_SECURITY_TYPE_PREFERRED_SHARE
	case entities.SecurityTypeDepositaryReceipt:
		return companyv1.SecurityType_SECURITY_TYPE_DEPOSITARY_RECEIPT
	default:
		return companyv1.SecurityType_SECURITY_TYPE_UNSPECIFIED
	}
}

// toProtoListingLevel переводит domain-enum в proto-enum.
func toProtoListingLevel(level entities.ListingLevel) companyv1.ListingLevel {
	switch level {
	case entities.ListingLevelFirst:
		return companyv1.ListingLevel_LISTING_LEVEL_FIRST
	case entities.ListingLevelSecond:
		return companyv1.ListingLevel_LISTING_LEVEL_SECOND
	case entities.ListingLevelThird:
		return companyv1.ListingLevel_LISTING_LEVEL_THIRD
	default:
		return companyv1.ListingLevel_LISTING_LEVEL_UNSPECIFIED
	}
}

// toProtoCompanyCard переводит entities.CompanyCard в proto-сообщение.
func toProtoCompanyCard(c *entities.CompanyCard) *companyv1.CompanyCard {
	return &companyv1.CompanyCard{
		Ticker:                c.Ticker,
		Exchange:              c.Exchange,
		Name:                  c.Name,
		Sector:                c.Sector,
		SectorId:              toInt32(c.SectorID),
		Industry:              c.Industry,
		IndustryId:            toInt32(c.IndustryID),
		IndustryGroup:         c.IndustryGroup,
		IndustryGroupId:       toInt32(c.IndustryGroupID),
		Country:               c.Country,
		Currency:              c.Currency,
		PrimaryReportTicker:   c.PrimaryReportTicker,
		PrimaryReportExchange: c.PrimaryReportExchange,
		Description:           c.Description,
		Site:                  c.Site,
		DiscLink:              c.DiscLink,
	}
}

// toInt32 безопасно сужает int до int32. GICS-коды (сектор / отрасль /
// группа) укладываются в int32 (значения < 10^9); выход за границы для этих
// доменных полей нереалистичен, поэтому возвращаем 0.
func toInt32(v int) int32 {
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0
	}
	return int32(v)
}
