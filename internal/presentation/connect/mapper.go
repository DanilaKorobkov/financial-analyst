// Package connect — Connect-handler для CompanyService.
package connect

import (
	"time"

	"github.com/samber/lo"
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
	stockPeriodFrequencyProto = map[company.StockPeriodFrequency]companyv1.StockPeriodFrequency{
		company.StockPeriodFrequencyUnspecified: companyv1.StockPeriodFrequency_STOCK_PERIOD_FREQUENCY_UNSPECIFIED,
		company.StockPeriodFrequencyYearly:      companyv1.StockPeriodFrequency_STOCK_PERIOD_FREQUENCY_YEARLY,
		company.StockPeriodFrequencyHalfYearly:  companyv1.StockPeriodFrequency_STOCK_PERIOD_FREQUENCY_HALF_YEARLY,
		company.StockPeriodFrequencyQuarterly:   companyv1.StockPeriodFrequency_STOCK_PERIOD_FREQUENCY_QUARTERLY,
		company.StockPeriodFrequencyYearToMonth: companyv1.StockPeriodFrequency_STOCK_PERIOD_FREQUENCY_YEAR_TO_MONTH,
	}
	reportStandardProto = map[company.ReportStandard]companyv1.ReportStandard{
		company.ReportStandardUnspecified: companyv1.ReportStandard_REPORT_STANDARD_UNSPECIFIED,
		company.ReportStandardIFRS:        companyv1.ReportStandard_IFRS,
		company.ReportStandardRAS:         companyv1.ReportStandard_RAS,
		company.ReportStandardGAAP:        companyv1.ReportStandard_GAAP,
	}
	dividendTypeProto = map[company.DividendType]companyv1.DividendType{
		company.DividendTypeUnspecified: companyv1.DividendType_DIVIDEND_TYPE_UNSPECIFIED,
		company.DividendTypeYearly:      companyv1.DividendType_DIVIDEND_TYPE_YEARLY,
		company.DividendTypeFirstHalf:   companyv1.DividendType_DIVIDEND_TYPE_FIRST_HALF,
		company.DividendTypeSecondHalf:  companyv1.DividendType_DIVIDEND_TYPE_SECOND_HALF,
		company.DividendTypeQ1:          companyv1.DividendType_DIVIDEND_TYPE_Q1,
		company.DividendTypeQ2:          companyv1.DividendType_DIVIDEND_TYPE_Q2,
		company.DividendTypeQ3:          companyv1.DividendType_DIVIDEND_TYPE_Q3,
		company.DividendTypeQ4:          companyv1.DividendType_DIVIDEND_TYPE_Q4,
		company.DividendTypeSpecial:     companyv1.DividendType_DIVIDEND_TYPE_SPECIAL,
	}
	ideaStatusProto = map[company.IdeaStatus]companyv1.IdeaStatus{
		company.IdeaStatusUnspecified: companyv1.IdeaStatus_IDEA_STATUS_UNSPECIFIED,
		company.IdeaStatusActive:      companyv1.IdeaStatus_ACTIVE,
		company.IdeaStatusClosed:      companyv1.IdeaStatus_CLOSED,
	}
	insiderTransactionTypeProto = map[company.InsiderTransactionType]companyv1.InsiderTransactionType{
		company.InsiderTransactionTypeUnspecified: companyv1.InsiderTransactionType_INSIDER_TRANSACTION_TYPE_UNSPECIFIED,
		company.InsiderTransactionTypePurchase:    companyv1.InsiderTransactionType_PURCHASE,
		company.InsiderTransactionTypeSale:        companyv1.InsiderTransactionType_SALE,
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
		Info:                toProtoStockInfo(&s.Info),
		Summary:             toProtoStockSummary(&s.Summary),
		Ratios:              toProtoStockRatios(s.Ratios),
		Reports:             toProtoStockReports(s.Reports),
		Dividends:           toProtoStockDividends(s.Dividends),
		Ideas:               toProtoStockIdeas(s.Ideas),
		InsiderTransactions: toProtoStockInsiderTransactions(s.InsiderTransactions),
		Operations:          toProtoStockOperations(s.Operations),
		Owners:              toProtoStockOwners(s.Owners),
		Shares:              toProtoStockShares(s.Shares),
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

func toProtoStockPeriod(p company.StockPeriod) *companyv1.StockPeriod {
	return &companyv1.StockPeriod{
		Year:      int64(p.Year),
		Month:     int64(p.Month),
		Frequency: stockPeriodFrequencyProto[p.Frequency],
		Standard:  reportStandardProto[p.Standard],
	}
}

func toProtoStockRatios(rs []company.StockRatio) []*companyv1.StockRatio {
	if len(rs) == 0 {
		return nil
	}
	return lo.Map(rs, func(r company.StockRatio, _ int) *companyv1.StockRatio {
		return &companyv1.StockRatio{
			Period:            toProtoStockPeriod(r.Period),
			ChangedAt:         encodeDate(r.ChangedAt),
			Capital:           r.Capital,
			Pe:                r.PE,
			Pbv:               r.PBV,
			Ps:                r.PS,
			Pcf:               r.PCF,
			Pfcf:              r.PFCF,
			Pffo:              r.PFFO,
			Evs:               r.EVS,
			EvEbitda:          r.EVEBITDA,
			EvEbit:            r.EVEBIT,
			DebtEbitda:        r.DebtEBITDA,
			NetDebtEbitda:     r.NetDebtEBITDA,
			DebtEquity:        r.DebtEquity,
			DebtRatio:         r.DebtRatio,
			CurrentRatio:      r.CurrentRatio,
			InterestCoverage:  r.InterestCoverage,
			GrossMargin:       r.GrossMargin,
			OperationMargin:   r.OperationMargin,
			EbitdaMargin:      r.EBITDAMargin,
			NetMargin:         r.NetMargin,
			Ros:               r.ROS,
			Roe:               r.ROE,
			Roa:               r.ROA,
			Roic:              r.ROIC,
			Roce:              r.ROCE,
			Dpr:               r.DPR,
			CapexRevenue:      r.CapexRevenue,
			NetWorkingCapital: r.NetWorkingCapital,
			Active:            r.Active,
		}
	})
}

// toProtoStockReports переводит срез StockReport в proto-срез.
// Каждая запись отчёта собирается группами строк (заголовок, P&L,
// денежный поток, баланс, per-share) — иначе одна функция-сборщик
// превышает порог funlen на 95 полях отчёта.
func toProtoStockReports(rs []company.StockReport) []*companyv1.StockReport {
	if len(rs) == 0 {
		return nil
	}
	return lo.Map(rs, func(r company.StockReport, _ int) *companyv1.StockReport {
		out := protoReportHeader(&r)
		applyProtoReportProfitAndLoss(out, &r)
		applyProtoReportCashFlow(out, &r)
		applyProtoReportBalance(out, &r)
		applyProtoReportPerShare(out, &r)
		return out
	})
}

// protoReportHeader заполняет в proto-StockReport поля периода,
// валюты, масштаба, ссылок и флага предварительной публикации.
func protoReportHeader(r *company.StockReport) *companyv1.StockReport {
	return &companyv1.StockReport{
		Period:      toProtoStockPeriod(r.Period),
		ChangedAt:   encodeDate(r.ChangedAt),
		Currency:    currencyProto[r.Currency],
		Amount:      r.Amount,
		Link:        r.Link,
		LinkPress:   r.LinkPress,
		LinkUpdate:  r.LinkUpdate,
		Preliminary: r.Preliminary,
	}
}

// applyProtoReportProfitAndLoss заполняет в proto-StockReport строки
// отчёта о прибылях и убытках.
func applyProtoReportProfitAndLoss(out *companyv1.StockReport, r *company.StockReport) {
	out.Revenue = r.Revenue
	out.CostOfSales = r.CostOfSales
	out.GrossProfit = r.GrossProfit
	out.SelGenAdmExpenses = r.SelGenAdmExpenses
	out.OperatingIncome = r.OperatingIncome
	out.Ebit = r.EBIT
	out.Ebitda = r.EBITDA
	out.EbitdaAdjusted = r.EBITDAAdjusted
	out.Earnings = r.Earnings
	out.EarningsWoTax = r.EarningsWoTax
	out.EarningsComprehensive = r.EarningsComprehensive
	out.EarningsComprehensiveStockHolders = r.EarningsComprehensiveStockHolders
	out.EarningsStockHolders = r.EarningsStockHolders
	out.EarningsContinuingOperations = r.EarningsContinuingOperations
	out.TotalExpenses = r.TotalExpenses
	out.OtherOperatingIncome = r.OtherOperatingIncome
	out.DeprDeplAmort = r.DeprDeplAmort
	out.ResearchAndDevelopment = r.ResearchAndDevelopment
	out.Ffo = r.FFO
	out.InterestIncome = r.InterestIncome
	out.InterestExpense = r.InterestExpense
	out.InterestNet = r.InterestNet
	out.CommissionIncome = r.CommissionIncome
	out.CommissionExpense = r.CommissionExpense
	out.CommissionNet = r.CommissionNet
}

// applyProtoReportCashFlow заполняет в proto-StockReport строки отчёта
// о движении денежных средств.
func applyProtoReportCashFlow(out *companyv1.StockReport, r *company.StockReport) {
	out.Cfo = r.CFO
	out.Cfi = r.CFI
	out.Cff = r.CFF
	out.Fcf = r.FCF
	out.FcfAdjusted = r.FCFAdjusted
	out.NetChangeInCash = r.NetChangeInCash
	out.Capex = r.Capex
	out.PpePurchase = r.PPEPurchase
	out.IntangiblesPurchase = r.IntangiblesPurchase
	out.RepurchaseOfStock = r.RepurchaseOfStock
	out.IssuanceOfDebt = r.IssuanceOfDebt
	out.PaymentsOfDebt = r.PaymentsOfDebt
	out.NetIssuanceOfDebt = r.NetIssuanceOfDebt
	out.PaymentsForDividends = r.PaymentsForDividends
	out.CashPaidForInterest = r.CashPaidForInterest
	out.CashPaidForTax = r.CashPaidForTax
}

// applyProtoReportBalance заполняет в proto-StockReport строки баланса.
func applyProtoReportBalance(out *companyv1.StockReport, r *company.StockReport) {
	out.TotalAssets = r.TotalAssets
	out.CurrentAssets = r.CurrentAssets
	out.LongTermAssets = r.LongTermAssets
	out.CashAndEquiv = r.CashAndEquiv
	out.ShortTermInvestments = r.ShortTermInvestments
	out.CashEquivStInvesments = r.CashEquivSTInvestments
	out.AccountsReceivable = r.AccountsReceivable
	out.OtherReceivable = r.OtherReceivable
	out.TotalReceivable = r.TotalReceivable
	out.Inventories = r.Inventories
	out.PropertyPlantEquipment = r.PropertyPlantEquipment
	out.PpeRou = r.PPERoU
	out.RightOfUseAssets = r.RightOfUseAssets
	out.IntangibleAssets = r.IntangibleAssets
	out.Goodwill = r.Goodwill
	out.GoodwillIntangibleAssets = r.GoodwillIntangibleAssets
	out.IntangibleAndTangibleAssets = r.IntangibleAndTangibleAssets
	out.LongTermInvestments = r.LongTermInvestments
	out.TotalLiabilities = r.TotalLiabilities
	out.CurrentLiabilities = r.CurrentLiabilities
	out.LongTermLiabilities = r.LongTermLiabilities
	out.CurrentDebt = r.CurrentDebt
	out.LongTermDebt = r.LongTermDebt
	out.TotalDebt = r.TotalDebt
	out.NetDebt = r.NetDebt
	out.NetDebtAdjusted = r.NetDebtAdjusted
	out.CurLongDebt = r.CurLongDebt
	out.CurLongLease = r.CurLongLease
	out.CurrentLease = r.CurrentLease
	out.LongTermLease = r.LongTermLease
	out.AccountsPayable = r.AccountsPayable
	out.OtherPayable = r.OtherPayable
	out.TotalPayable = r.TotalPayable
	out.Equity = r.Equity
	out.EquityStockHolders = r.EquityStockHolders
	out.RetainedEarnings = r.RetainedEarnings
	out.SharePremium = r.SharePremium
	out.TreasuryStock = r.TreasuryStock
}

// applyProtoReportPerShare заполняет в proto-StockReport per-share метрики.
func applyProtoReportPerShare(out *companyv1.StockReport, r *company.StockReport) {
	out.EarningsPs = r.EarningsPS
	out.EquityPs = r.EquityPS
	out.RevenuePs = r.RevenuePS
	out.EbitdaPs = r.EBITDAPS
	out.EbitdaAdjustedPs = r.EBITDAAdjustedPS
	out.FcfPs = r.FCFPS
	out.FcfAdjustedPs = r.FCFAdjustedPS
}

func toProtoStockDividends(ds []company.StockDividend) []*companyv1.StockDividend {
	if len(ds) == 0 {
		return nil
	}
	return lo.Map(ds, func(d company.StockDividend, _ int) *companyv1.StockDividend {
		return &companyv1.StockDividend{
			LastBuyDate:     encodeDate(d.LastBuyDate),
			ReestrCloseDate: encodeDate(d.ReestrCloseDate),
			ChangedAt:       encodeDate(d.ChangedAt),
			LastBuyPrice:    d.LastBuyPrice,
			DivAmount:       d.DivAmount,
			DivPercent:      d.DivPercent,
			Year:            d.Year,
			Link:            d.Link,
			Currency:        currencyProto[d.Currency],
			Type:            dividendTypeProto[d.Type],
		}
	})
}

func toProtoStockIdeas(is []company.StockIdea) []*companyv1.StockIdea {
	if len(is) == 0 {
		return nil
	}
	return lo.Map(is, func(i company.StockIdea, _ int) *companyv1.StockIdea {
		return &companyv1.StockIdea{
			DateIn:          encodeDate(i.DateIn),
			DateOut:         encodeDate(i.DateOut),
			CloseDate:       encodeDate(i.CloseDate),
			UpdateDate:      encodeDate(i.UpdateDate),
			ChangedAt:       encodeDate(i.ChangedAt),
			Community:       i.Community,
			Idea:            i.Idea,
			CloseComment:    i.CloseComment,
			CloseLink:       i.CloseLink,
			Id:              i.ID,
			CommunityId:     i.CommunityID,
			DurationInMonth: i.DurationInMonth,
			PriceIn:         i.PriceIn,
			PriceOut:        i.PriceOut,
			PriceDay:        i.PriceDay,
			ProfitPotential: i.ProfitPotential,
			ProfitActual:    i.ProfitActual,
			StopLoss:        i.StopLoss,
			ClosePrice:      i.ClosePrice,
			UpdatePrice:     i.UpdatePrice,
			Status:          ideaStatusProto[i.Status],
		}
	})
}

func toProtoStockInsiderTransactions(ts []company.StockInsiderTransaction) []*companyv1.StockInsiderTransaction {
	if len(ts) == 0 {
		return nil
	}
	return lo.Map(ts, func(t company.StockInsiderTransaction, _ int) *companyv1.StockInsiderTransaction {
		return &companyv1.StockInsiderTransaction{
			TransactionDate: encodeDate(t.TransactionDate),
			Insider:         t.Insider,
			InsiderTitle:    t.InsiderTitle,
			Type:            insiderTransactionTypeProto[t.Type],
		}
	})
}

func toProtoStockOperations(os []company.StockOperation) []*companyv1.StockOperation {
	if len(os) == 0 {
		return nil
	}
	return lo.Map(os, func(o company.StockOperation, _ int) *companyv1.StockOperation {
		return &companyv1.StockOperation{
			Period:         toProtoStockPeriod(o.Period),
			MetricId:       o.MetricID,
			Unit:           o.Unit,
			OriginalUnit:   o.OriginalUnit,
			Link:           o.Link,
			LinkUpdate:     o.LinkUpdate,
			Amount:         o.Amount,
			OriginalAmount: o.OriginalAmount,
			Value:          o.Value,
			OriginalValue:  o.OriginalValue,
			Curs:           o.Curs,
		}
	})
}

func toProtoStockOwners(os []company.StockOwner) []*companyv1.StockOwner {
	if len(os) == 0 {
		return nil
	}
	return lo.Map(os, func(o company.StockOwner, _ int) *companyv1.StockOwner {
		return &companyv1.StockOwner{
			ChangedAt: encodeDate(o.ChangedAt),
			Owner:     o.Owner,
			Link:      o.Link,
			Period:    toProtoStockPeriod(o.Period),
			Own:       o.Own,
		}
	})
}

func toProtoStockShares(ss []company.StockShare) []*companyv1.StockShare {
	if len(ss) == 0 {
		return nil
	}
	return lo.Map(ss, func(s company.StockShare, _ int) *companyv1.StockShare {
		return &companyv1.StockShare{
			Ticker: s.Ticker,
			Period: toProtoStockPeriod(s.Period),
			Num:    s.Num,
		}
	})
}
