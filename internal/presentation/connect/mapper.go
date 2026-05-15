// Package connect — Connect-handler для CompanyService.
package connect

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

// Таблицы перевода доменных enum-ов в proto enum-ы. Численные значения
// домена и proto-схемы фиксированы и совпадают, но явные таблицы
// держат маппинг под контролем: добавление нового значения в домен
// либо появляется здесь сознательно, либо упадёт в default-ветке
// proto-enum как Unspecified.
var (
	securityTypeProto = map[company.SecurityType]companyv1.SecurityType{
		company.SecurityTypeUnspecified:       companyv1.SecurityType_SECURITY_TYPE_UNSPECIFIED,
		company.SecurityTypeCommonShare:       companyv1.SecurityType_COMMON_SHARE,
		company.SecurityTypePreferredShare:    companyv1.SecurityType_PREFERRED_SHARE,
		company.SecurityTypeDepositaryReceipt: companyv1.SecurityType_DEPOSITARY_RECEIPT,
	}
	listingLevelProto = map[company.ListingLevel]companyv1.ListingLevel{
		company.ListingLevelUnspecified: companyv1.ListingLevel_LISTING_LEVEL_UNSPECIFIED,
		company.ListingLevelFirst:       companyv1.ListingLevel_FIRST,
		company.ListingLevelSecond:      companyv1.ListingLevel_SECOND,
		company.ListingLevelThird:       companyv1.ListingLevel_THIRD,
	}
	currencyProto = map[company.Currency]companyv1.Currency{
		company.CurrencyUnspecified: companyv1.Currency_CURRENCY_UNSPECIFIED,
		company.CurrencyRUB:         companyv1.Currency_RUB,
		company.CurrencyUSD:         companyv1.Currency_USD,
		company.CurrencyEUR:         companyv1.Currency_EUR,
	}
	exchangeProto = map[company.Exchange]companyv1.Exchange{
		company.ExchangeUnspecified: companyv1.Exchange_EXCHANGE_UNSPECIFIED,
		company.ExchangeMOEX:        companyv1.Exchange_MOEX,
	}
	reportFrequencyProto = map[company.ReportFrequency]companyv1.ReportFrequency{
		company.ReportFrequencyUnspecified: companyv1.ReportFrequency_REPORT_FREQUENCY_UNSPECIFIED,
		company.ReportFrequencyYearly:      companyv1.ReportFrequency_YEARLY,
		company.ReportFrequencyQuarterly:   companyv1.ReportFrequency_QUARTERLY,
	}
	ideaConsensusProto = map[company.IdeaConsensus]companyv1.IdeaConsensus{
		company.IdeaConsensusUnspecified: companyv1.IdeaConsensus_IDEA_CONSENSUS_UNSPECIFIED,
		company.IdeaConsensusBuy:         companyv1.IdeaConsensus_BUY,
		company.IdeaConsensusHold:        companyv1.IdeaConsensus_HOLD,
		company.IdeaConsensusSell:        companyv1.IdeaConsensus_SELL,
	}
	insiderConsensusProto = map[company.InsiderConsensus]companyv1.InsiderConsensus{
		company.InsiderConsensusUnspecified: companyv1.InsiderConsensus_INSIDER_CONSENSUS_UNSPECIFIED,
		company.InsiderConsensusBuys:        companyv1.InsiderConsensus_BUYS,
		company.InsiderConsensusSells:       companyv1.InsiderConsensus_SELLS,
		company.InsiderConsensusMixed:       companyv1.InsiderConsensus_MIXED,
	}
)

// toProtoCompany переводит domain-агрегат в proto-сообщение. Принимает
// указатель: домен оперирует value-типами, но Company содержит две
// тяжёлые секции (общий размер ~512 байт) — копировать на каждый вызов
// маппера незачем.
func toProtoCompany(c *company.Company) *companyv1.Company {
	return &companyv1.Company{
		SecurityDescription: toProtoSecurityDescription(&c.SecurityDescription),
		Stock:               toProtoStock(&c.Stock),
	}
}

func toProtoStock(s *company.Stock) *companyv1.Stock {
	return &companyv1.Stock{
		Info:    toProtoStockInfo(&s.Info),
		Summary: toProtoStockSummary(&s.Summary),
	}
}

func toProtoSecurityDescription(d *company.SecurityDescription) *companyv1.SecurityDescription {
	return &companyv1.SecurityDescription{
		Ticker:                 d.Ticker,
		Isin:                   d.ISIN,
		Name:                   d.Name,
		ShortName:              d.ShortName,
		IssueName:              d.IssueName,
		LatName:                d.LatName,
		RegNumber:              d.RegNumber,
		SecurityTypeName:       d.SecurityTypeName,
		SecurityGroup:          d.SecurityGroup,
		SecurityGroupName:      d.SecurityGroupName,
		SecurityType:           securityTypeProto[d.SecurityType],
		ListingLevel:           listingLevelProto[d.ListingLevel],
		FaceValue:              d.FaceValue,
		FaceUnit:               currencyProto[d.FaceUnit],
		IssueSize:              d.IssueSize,
		IssueDate:              encodeDate(d.IssueDate),
		RegistryDate:           encodeDate(d.RegistryDate),
		EmitterId:              d.EmitterID,
		HasProspectus:          d.HasProspectus,
		HasDefault:             d.HasDefault,
		HasTechnicalDefault:    d.HasTechnicalDefault,
		EmitentMismatchCurrent: d.EmitentMismatchCurrent,
		IsQualifiedInvestors:   d.IsQualifiedInvestors,
		MorningSession:         d.MorningSession,
		EveningSession:         d.EveningSession,
		WeekendSession:         d.WeekendSession,
	}
}

func toProtoStockInfo(i *company.StockInfo) *companyv1.StockInfo {
	return &companyv1.StockInfo{
		IssuerName:            i.IssuerName,
		Sector:                i.Sector,
		IndustryGroup:         i.IndustryGroup,
		Industry:              i.Industry,
		SubIndustry:           i.SubIndustry,
		SectorId:              i.SectorID,
		IndustryGroupId:       i.IndustryGroupID,
		IndustryId:            i.IndustryID,
		SubIndustryId:         i.SubIndustryID,
		Country:               i.Country,
		Description:           i.Description,
		Site:                  i.Site,
		DisclosureLink:        i.DisclosureLink,
		PrimaryReportTicker:   i.PrimaryReportTicker,
		PrimaryReportExchange: exchangeProto[i.PrimaryReportExchange],
		Exchange:              exchangeProto[i.Exchange],
		Currency:              currencyProto[i.Currency],
		ReportFrequency:       reportFrequencyProto[i.ReportFrequency],
		Spb:                   i.SPB,
	}
}

func toProtoStockSummary(s *company.StockSummary) *companyv1.StockSummary {
	return &companyv1.StockSummary{
		Capital:            s.Capital,
		Eps:                s.EPS,
		Peg:                s.PEG,
		PeterLynchTarget:   s.PeterLynchTarget,
		GrahamTarget:       s.GrahamTarget,
		DividendFrequency:  int64(s.DividendFrequency),
		DividendStrike:     int64(s.DividendStrike),
		DividendGrowth:     int64(s.DividendGrowth),
		DividendIndex:      s.DividendIndex,
		DividendYield_12M:  s.DividendYield12M,
		DividendYield_3Y:   s.DividendYield3Y,
		DividendYield_5Y:   s.DividendYield5Y,
		DividendGapLast:    int64(s.DividendGapLast),
		DividendGapAverage: int64(s.DividendGapAverage),
		GrowthRevenue_3Y:   s.GrowthRevenue3Y,
		GrowthRevenue_5Y:   s.GrowthRevenue5Y,
		GrowthEarnings_3Y:  s.GrowthEarnings3Y,
		GrowthEarnings_5Y:  s.GrowthEarnings5Y,
		GrowthEbitda_3Y:    s.GrowthEBITDA3Y,
		GrowthEbitda_5Y:    s.GrowthEBITDA5Y,
		GrowthAssets_3Y:    s.GrowthAssets3Y,
		GrowthAssets_5Y:    s.GrowthAssets5Y,
		GrowthEquity_3Y:    s.GrowthEquity3Y,
		GrowthEquity_5Y:    s.GrowthEquity5Y,
		GrowthFcf_3Y:       s.GrowthFCF3Y,
		GrowthFcf_5Y:       s.GrowthFCF5Y,
		GrowthNetDebt_3Y:   s.GrowthNetDebt3Y,
		GrowthNetDebt_5Y:   s.GrowthNetDebt5Y,
		GrowthOperation_3Y: s.GrowthOperation3Y,
		GrowthOperation_5Y: s.GrowthOperation5Y,
		IdeaBuy:            int64(s.IdeaBuy),
		IdeaHold:           int64(s.IdeaHold),
		IdeaSell:           int64(s.IdeaSell),
		IdeaTarget:         s.IdeaTarget,
		IdeaPotential:      s.IdeaPotential,
		IdeaConsensus:      ideaConsensusProto[s.IdeaConsensus],
		InsiderConsensus:   insiderConsensusProto[s.InsiderConsensus],
		ChangedAt:          encodeDate(s.ChangedAt),
	}
}

// encodeDate переводит time.Time в Timestamp. Нулевой time.Time
// (поле не отдано источником) превращается в nil — поле в proto-сообщении
// будет отсутствовать.
func encodeDate(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
