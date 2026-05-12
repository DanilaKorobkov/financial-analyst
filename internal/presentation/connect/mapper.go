// Package connect — Connect-handler для CompanyService.
package connect

import (
	"math"

	"google.golang.org/protobuf/types/known/timestamppb"

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
	}
}

// toInt32 безопасно сужает int до int32. GICS-коды и счётчики дивидендов /
// инвестидей укладываются в int32 (значения < 10^9); за пределами диапазона
// возвращаем 0 — выход за границы для этих доменных полей нереалистичен.
func toInt32(v int) int32 {
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0
	}
	return int32(v)
}

// toProtoCompanyMetrics переводит entities.CompanyMetrics в proto-сообщение.
func toProtoCompanyMetrics(m *entities.CompanyMetrics) *companyv1.CompanyMetrics {
	var changedAt *timestamppb.Timestamp
	if !m.ChangedAt.IsZero() {
		changedAt = timestamppb.New(m.ChangedAt)
	}
	return &companyv1.CompanyMetrics{
		Card:               toProtoCompanyCard(&m.Card),
		Description:        m.Description,
		Site:               m.Site,
		DiscLink:           m.DiscLink,
		Capital:            m.Capital,
		Eps:                m.EPS,
		Peg:                m.PEG,
		PeterLynchTarget:   m.PeterLynchTarget,
		GrahamTarget:       m.GrahamTarget,
		DividendFrequency:  toInt32(m.DividendFrequency),
		DividendStrike:     toInt32(m.DividendStrike),
		DividendGrowth:     toInt32(m.DividendGrowth),
		DividendIndex:      m.DividendIndex,
		DividendYield_12M:  m.DividendYield12m,
		DividendYield_3Y:   m.DividendYield3y,
		DividendYield_5Y:   m.DividendYield5y,
		DividendGapLast:    toInt32(m.DividendGapLast),
		DividendGapAverage: toInt32(m.DividendGapAverage),
		GrowthRevenue_3Y:   m.GrowthRevenue3y,
		GrowthRevenue_5Y:   m.GrowthRevenue5y,
		GrowthEarnings_3Y:  m.GrowthEarnings3y,
		GrowthEarnings_5Y:  m.GrowthEarnings5y,
		GrowthEbitda_3Y:    m.GrowthEbitda3y,
		GrowthEbitda_5Y:    m.GrowthEbitda5y,
		GrowthAssets_3Y:    m.GrowthAssets3y,
		GrowthAssets_5Y:    m.GrowthAssets5y,
		GrowthEquity_3Y:    m.GrowthEquity3y,
		GrowthEquity_5Y:    m.GrowthEquity5y,
		GrowthFcf_3Y:       m.GrowthFCF3y,
		GrowthFcf_5Y:       m.GrowthFCF5y,
		GrowthNetdebt_3Y:   m.GrowthNetDebt3y,
		GrowthNetdebt_5Y:   m.GrowthNetDebt5y,
		GrowthOperation_3Y: m.GrowthOperation3y,
		GrowthOperation_5Y: m.GrowthOperation5y,
		IdeaBuy:            toInt32(m.IdeaBuy),
		IdeaHold:           toInt32(m.IdeaHold),
		IdeaSell:           toInt32(m.IdeaSell),
		IdeaConsensus:      toProtoIdeaConsensus(m.IdeaConsensus),
		IdeaTarget:         m.IdeaTarget,
		IdeaPotential:      m.IdeaPotential,
		InsiderConsensus:   toProtoInsiderConsensus(m.InsiderConsensus),
		ChangedAt:          changedAt,
	}
}

// toProtoIdeaConsensus переводит domain-enum в proto-enum.
func toProtoIdeaConsensus(c entities.IdeaConsensus) companyv1.IdeaConsensus {
	switch c {
	case entities.IdeaConsensusBuy:
		return companyv1.IdeaConsensus_IDEA_CONSENSUS_BUY
	case entities.IdeaConsensusHold:
		return companyv1.IdeaConsensus_IDEA_CONSENSUS_HOLD
	case entities.IdeaConsensusSell:
		return companyv1.IdeaConsensus_IDEA_CONSENSUS_SELL
	default:
		return companyv1.IdeaConsensus_IDEA_CONSENSUS_UNSPECIFIED
	}
}

// toProtoInsiderConsensus переводит domain-enum в proto-enum.
func toProtoInsiderConsensus(c entities.InsiderConsensus) companyv1.InsiderConsensus {
	switch c {
	case entities.InsiderConsensusBuys:
		return companyv1.InsiderConsensus_INSIDER_CONSENSUS_BUYS
	case entities.InsiderConsensusSells:
		return companyv1.InsiderConsensus_INSIDER_CONSENSUS_SELLS
	case entities.InsiderConsensusMixed:
		return companyv1.InsiderConsensus_INSIDER_CONSENSUS_MIXED
	default:
		return companyv1.InsiderConsensus_INSIDER_CONSENSUS_UNSPECIFIED
	}
}
