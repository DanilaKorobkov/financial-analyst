package stock

import (
	"strconv"
	"time"

	"github.com/samber/lo"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

const (
	// changedAtLayout — формат datetime, в котором FinanceMarker отдаёт
	// поле changed_at (ISO 8601 без таймзоны, секунды).
	changedAtLayout = "2006-01-02T15:04:05"

	// dateLayout — формат date FM (поля last_buy_date, reestr_close_date,
	// date_in, date_out, close_date, update_date, transaction_date).
	dateLayout = "2006-01-02"
)

var (
	// exchangeByCode переводит строковый код биржи FinanceMarker в domain-enum.
	exchangeByCode = map[string]company.Exchange{
		codeExchangeMOEX: company.ExchangeMOEX,
	}

	// currencyByCode переводит ISO 4217-код валюты FinanceMarker в domain-enum.
	currencyByCode = map[string]company.Currency{
		"RUB": company.CurrencyRUB,
		"USD": company.CurrencyUSD,
		"EUR": company.CurrencyEUR,
	}

	// reportFrequencyByCode переводит код report_frequency FinanceMarker
	// ("Y" / "Q") в domain-enum.
	reportFrequencyByCode = map[string]company.ReportFrequency{
		"Y": company.ReportFrequencyYearly,
		"Q": company.ReportFrequencyQuarterly,
	}

	// ideaConsensusByCode переводит строковый код консенсуса инвест-идей
	// FinanceMarker в domain-enum.
	ideaConsensusByCode = map[string]company.IdeaConsensus{
		"BUY":  company.IdeaConsensusBuy,
		"HOLD": company.IdeaConsensusHold,
		"SELL": company.IdeaConsensusSell,
	}

	// insiderConsensusByCode переводит строковый код консенсуса сделок
	// инсайдеров FinanceMarker в domain-enum.
	insiderConsensusByCode = map[string]company.InsiderConsensus{
		"BUYS":  company.InsiderConsensusBuys,
		"SELLS": company.InsiderConsensusSells,
		"MIXED": company.InsiderConsensusMixed,
	}

	// stockPeriodFrequencyByCode переводит код period записи (Y/Q/H/YTM)
	// в domain-enum гранулярности периода.
	stockPeriodFrequencyByCode = map[string]company.StockPeriodFrequency{
		"Y":   company.StockPeriodFrequencyYearly,
		"Q":   company.StockPeriodFrequencyQuarterly,
		"H":   company.StockPeriodFrequencyHalfYearly,
		"YTM": company.StockPeriodFrequencyYearToMonth,
	}

	// reportStandardByCode переводит код стандарта отчётности FinanceMarker
	// ("МСФО" / "РСБУ" / "GAAP") в domain-enum.
	reportStandardByCode = map[string]company.ReportStandard{
		"МСФО": company.ReportStandardIFRS,
		"РСБУ": company.ReportStandardRAS,
		"GAAP": company.ReportStandardGAAP,
	}

	// dividendTypeByCode переводит код type дивидендной выплаты
	// (Y/S1/S2/Q1..Q4/O) в domain-enum.
	dividendTypeByCode = map[string]company.DividendType{
		"Y":  company.DividendTypeYearly,
		"S1": company.DividendTypeFirstHalf,
		"S2": company.DividendTypeSecondHalf,
		"Q1": company.DividendTypeQ1,
		"Q2": company.DividendTypeQ2,
		"Q3": company.DividendTypeQ3,
		"Q4": company.DividendTypeQ4,
		"O":  company.DividendTypeSpecial,
	}

	// ideaStatusByCode переводит system_status инвест-идеи в domain-enum.
	ideaStatusByCode = map[string]company.IdeaStatus{
		"ACTIVE": company.IdeaStatusActive,
		"CLOSED": company.IdeaStatusClosed,
	}

	// insiderTransactionTypeByCode переводит transaction_type (P/S)
	// в domain-enum направления сделки инсайдера.
	insiderTransactionTypeByCode = map[string]company.InsiderTransactionType{
		"P": company.InsiderTransactionTypePurchase,
		"S": company.InsiderTransactionTypeSale,
	}
)

// stockDTO — корневой объект ответа /api/fm/v2/stocks/{exchange}:{code}.
// Каждое поле соответствует блоку, который может быть запрошен через
// include. Не запрошенные блоки приходят пустыми объектами/массивами и
// раскладываются в zero-value секций.
type stockDTO struct {
	Ratios              []ratioDTO              `json:"ratios"`
	Reports             []reportDTO             `json:"reports"`
	Dividends           []dividendDTO           `json:"dividends"`
	Ideas               []ideaDTO               `json:"ideas"`
	InsiderTransactions []insiderTransactionDTO `json:"insiderTransactions"`
	Operations          []operationDTO          `json:"operations"`
	Owners              []ownerDTO              `json:"owners"`
	Shares              []shareDTO              `json:"shares"`
	Info                infoDTO                 `json:"info"`
	Summary             summaryDTO              `json:"summary"`
}

// infoDTO — блок `info` ответа эндпоинта.
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

// summaryDTO — блок `summary` ответа эндпоинта.
type summaryDTO struct {
	IdeaConsensus      string  `json:"idea_consensus"`
	InsiderConsensus   string  `json:"insider_consensus"`
	ChangedAt          string  `json:"changed_at"`
	Capital            float64 `json:"capital"`
	EPS                float64 `json:"eps"`
	PEG                float64 `json:"peg"`
	PeterLynchTarget   float64 `json:"peter_lynch_target"`
	GrahamTarget       float64 `json:"graham_target"`
	DividendIndex      float64 `json:"dividend_index"`
	DividendYield12M   float64 `json:"dividend_yield_12m"`
	DividendYield3Y    float64 `json:"dividend_yield_3y"`
	DividendYield5Y    float64 `json:"dividend_yield_5y"`
	GrowthRevenue3Y    float64 `json:"growth_revenue_3y"`
	GrowthRevenue5Y    float64 `json:"growth_revenue_5y"`
	GrowthEarnings3Y   float64 `json:"growth_earnings_3y"`
	GrowthEarnings5Y   float64 `json:"growth_earnings_5y"`
	GrowthEBITDA3Y     float64 `json:"growth_ebitda_3y"`
	GrowthEBITDA5Y     float64 `json:"growth_ebitda_5y"`
	GrowthAssets3Y     float64 `json:"growth_assets_3y"`
	GrowthAssets5Y     float64 `json:"growth_assets_5y"`
	GrowthEquity3Y     float64 `json:"growth_equity_3y"`
	GrowthEquity5Y     float64 `json:"growth_equity_5y"`
	GrowthFCF3Y        float64 `json:"growth_fcf_3y"`
	GrowthFCF5Y        float64 `json:"growth_fcf_5y"`
	GrowthNetDebt3Y    float64 `json:"growth_netdebt_3y"`
	GrowthNetDebt5Y    float64 `json:"growth_netdebt_5y"`
	GrowthOperation3Y  float64 `json:"growth_operation_3y"`
	GrowthOperation5Y  float64 `json:"growth_operation_5y"`
	IdeaTarget         float64 `json:"idea_target"`
	IdeaPotential      float64 `json:"idea_potential"`
	DividendFrequency  int     `json:"dividend_frequency"`
	DividendStrike     int     `json:"dividend_strike"`
	DividendGrowth     int     `json:"dividend_growth"`
	DividendGapLast    int     `json:"dividend_gap_last"`
	DividendGapAverage int     `json:"dividend_gap_average"`
	IdeaBuy            int     `json:"idea_buy"`
	IdeaHold           int     `json:"idea_hold"`
	IdeaSell           int     `json:"idea_sell"`
}

// ratioDTO — элемент блока `ratios` ответа эндпоинта.
type ratioDTO struct {
	ChangedAt         string  `json:"changed_at"`
	Period            string  `json:"period"`
	Type              string  `json:"type"`
	Year              int     `json:"year"`
	Month             int     `json:"month"`
	Capital           float64 `json:"capital"`
	PE                float64 `json:"pe"`
	PBV               float64 `json:"pbv"`
	PS                float64 `json:"ps"`
	PCF               float64 `json:"pcf"`
	PFCF              float64 `json:"pfcf"`
	PFFO              float64 `json:"pffo"`
	EVS               float64 `json:"evs"`
	EVEBITDA          float64 `json:"evebitda"`
	EVEBIT            float64 `json:"ev_ebit"`
	DebtEBITDA        float64 `json:"debtebitda"`
	NetDebtEBITDA     float64 `json:"netdebt_ebitda"`
	DebtEquity        float64 `json:"debt_equity"`
	DebtRatio         float64 `json:"debt_ratio"`
	CurrentRatio      float64 `json:"current_ratio"`
	InterestCoverage  float64 `json:"interest_coverage"`
	GrossMargin       float64 `json:"gross_margin"`
	OperationMargin   float64 `json:"operation_margin"`
	EBITDAMargin      float64 `json:"ebitda_margin"`
	NetMargin         float64 `json:"net_margin"`
	ROS               float64 `json:"ros"`
	ROE               float64 `json:"roe"`
	ROA               float64 `json:"roa"`
	ROIC              float64 `json:"roic"`
	ROCE              float64 `json:"roce"`
	DPR               float64 `json:"dpr"`
	CapexRevenue      float64 `json:"capex_revenue"`
	NetWorkingCapital float64 `json:"net_working_capital"`
	Active            bool    `json:"active"`
}

// reportDTO — элемент блока `reports` ответа эндпоинта.
type reportDTO struct {
	ChangedAt                         string  `json:"changed_at"`
	Period                            string  `json:"period"`
	Type                              string  `json:"type"`
	Curr                              string  `json:"curr"`
	Link                              string  `json:"link"`
	LinkPress                         string  `json:"link_press"`
	LinkUpdate                        string  `json:"link_update"`
	Year                              int     `json:"year"`
	Month                             int     `json:"month"`
	Amount                            int64   `json:"amount"`
	Revenue                           float64 `json:"revenue"`
	CostOfSales                       float64 `json:"cost_of_sales"`
	GrossProfit                       float64 `json:"gross_profit"`
	SelGenAdmExpenses                 float64 `json:"sel_gen_adm_expenses"`
	OperatingIncome                   float64 `json:"operating_income"`
	EBIT                              float64 `json:"ebit"`
	EBITDA                            float64 `json:"ebitda"`
	EBITDAAdjusted                    float64 `json:"ebitda_adjusted"`
	Earnings                          float64 `json:"earnings"`
	EarningsWoTax                     float64 `json:"earnings_wo_tax"`
	EarningsComprehensive             float64 `json:"earnings_comprehensive"`
	EarningsComprehensiveStockHolders float64 `json:"earnings_comprehensive_stock_holders"`
	EarningsStockHolders              float64 `json:"earnings_stock_holders"`
	EarningsContinuingOperations      float64 `json:"earnings_continuing_operations"`
	TotalExpenses                     float64 `json:"total_expenses"`
	OtherOperatingIncome              float64 `json:"other_operating_income"`
	DeprDeplAmort                     float64 `json:"depr_depl_amort"`
	ResearchAndDevelopment            float64 `json:"research_and_development"`
	FFO                               float64 `json:"ffo"`
	InterestIncome                    float64 `json:"interest_income"`
	InterestExpense                   float64 `json:"interest_expense"`
	InterestNet                       float64 `json:"interest_net"`
	CommissionIncome                  float64 `json:"commission_income"`
	CommissionExpense                 float64 `json:"commission_expense"`
	CommissionNet                     float64 `json:"commission_net"`
	CFO                               float64 `json:"cfo"`
	CFI                               float64 `json:"cfi"`
	CFF                               float64 `json:"cff"`
	FCF                               float64 `json:"fcf"`
	FCFAdjusted                       float64 `json:"fcf_adjusted"`
	NetChangeInCash                   float64 `json:"net_change_in_cash"`
	Capex                             float64 `json:"capex"`
	PPEPurchase                       float64 `json:"ppe_purchase"`
	IntangiblesPurchase               float64 `json:"intangibles_purchase"`
	RepurchaseOfStock                 float64 `json:"repurchase_of_stock"`
	IssuanceOfDebt                    float64 `json:"issuance_of_debt"`
	PaymentsOfDebt                    float64 `json:"payments_of_debt"`
	NetIssuanceOfDebt                 float64 `json:"net_issuance_of_debt"`
	PaymentsForDividends              float64 `json:"payments_for_dividends"`
	CashPaidForInterest               float64 `json:"cash_paid_for_interest"`
	CashPaidForTax                    float64 `json:"cash_paid_for_tax"`
	TotalAssets                       float64 `json:"total_assets"`
	CurrentAssets                     float64 `json:"current_assets"`
	LongTermAssets                    float64 `json:"long_term_assets"`
	CashAndEquiv                      float64 `json:"cash_and_equiv"`
	ShortTermInvestments              float64 `json:"short_term_investments"`
	CashEquivSTInvestments            float64 `json:"cash_equiv_st_invesments"`
	AccountsReceivable                float64 `json:"accounts_receivable"`
	OtherReceivable                   float64 `json:"other_receivable"`
	TotalReceivable                   float64 `json:"total_receivable"`
	Inventories                       float64 `json:"inventories"`
	PropertyPlantEquipment            float64 `json:"property_plant_equipment"`
	PPERoU                            float64 `json:"ppe_rou"`
	RightOfUseAssets                  float64 `json:"right_of_use_assets"`
	IntangibleAssets                  float64 `json:"intangible_assets"`
	Goodwill                          float64 `json:"goodwill"`
	GoodwillIntangibleAssets          float64 `json:"goodwill_intangible_assets"`
	IntangibleAndTangibleAssets       float64 `json:"intangible_and_tangible_assets"`
	LongTermInvestments               float64 `json:"long_term_investments"`
	TotalLiabilities                  float64 `json:"total_liabilities"`
	CurrentLiabilities                float64 `json:"current_liabilities"`
	LongTermLiabilities               float64 `json:"long_term_liabilities"`
	CurrentDebt                       float64 `json:"current_debt"`
	LongTermDebt                      float64 `json:"long_term_debt"`
	TotalDebt                         float64 `json:"total_debt"`
	NetDebt                           float64 `json:"net_debt"`
	NetDebtAdjusted                   float64 `json:"net_debt_adjusted"`
	CurLongDebt                       float64 `json:"cur_long_debt"`
	CurLongLease                      float64 `json:"cur_long_lease"`
	CurrentLease                      float64 `json:"current_lease"`
	LongTermLease                     float64 `json:"long_term_lease"`
	AccountsPayable                   float64 `json:"accounts_payable"`
	OtherPayable                      float64 `json:"other_payable"`
	TotalPayable                      float64 `json:"total_payable"`
	Equity                            float64 `json:"equity"`
	EquityStockHolders                float64 `json:"equity_stock_holders"`
	RetainedEarnings                  float64 `json:"retained_earnings"`
	SharePremium                      float64 `json:"share_premium"`
	TreasuryStock                     float64 `json:"treasury_stock"`
	EarningsPS                        float64 `json:"earnings_ps"`
	EquityPS                          float64 `json:"equity_ps"`
	RevenuePS                         float64 `json:"revenue_ps"`
	EBITDAPS                          float64 `json:"ebitda_ps"`
	EBITDAAdjustedPS                  float64 `json:"ebitda_adjusted_ps"`
	FCFPS                             float64 `json:"fcf_ps"`
	FCFAdjustedPS                     float64 `json:"fcf_adjusted_ps"`
	Preliminary                       bool    `json:"preliminary"`
}

// dividendDTO — элемент блока `dividends` ответа эндпоинта.
type dividendDTO struct {
	LastBuyDate     string  `json:"last_buy_date"`
	ReestrCloseDate string  `json:"reestr_close_date"`
	ChangedAt       string  `json:"changed_at"`
	DivCurr         string  `json:"div_curr"`
	Type            string  `json:"type"`
	Link            string  `json:"link"`
	LastBuyPrice    float64 `json:"last_buy_price"`
	DivAmount       float64 `json:"div_amount"`
	DivPercent      float64 `json:"div_percent"`
	Year            int64   `json:"year"`
}

// ideaDTO — элемент блока `ideas` ответа эндпоинта.
type ideaDTO struct {
	DateIn          string  `json:"date_in"`
	DateOut         string  `json:"date_out"`
	CloseDate       string  `json:"close_date"`
	UpdateDate      string  `json:"update_date"`
	ChangedAt       string  `json:"changed_at"`
	Community       string  `json:"community"`
	Idea            string  `json:"idea"`
	CloseComment    string  `json:"close_comment"`
	CloseLink       string  `json:"close_link"`
	SystemStatus    string  `json:"system_status"`
	ID              int64   `json:"id"`
	CommunityID     int64   `json:"community_id"`
	DurationInMonth int64   `json:"duration_in_month"`
	PriceIn         float64 `json:"price_in"`
	PriceOut        float64 `json:"price_out"`
	PriceDay        float64 `json:"price_day"`
	ProfitPotential float64 `json:"profit_potential"`
	ProfitActual    float64 `json:"profit_actual"`
	StopLoss        float64 `json:"stop_loss"`
	ClosePrice      float64 `json:"close_price"`
	UpdatePrice     float64 `json:"update_price"`
}

// insiderTransactionDTO — элемент блока `insiderTransactions` ответа
// эндпоинта (узкий формат).
type insiderTransactionDTO struct {
	TransactionDate string `json:"transaction_date"`
	Insider         string `json:"insider"`
	InsiderTitle    string `json:"insider_title"`
	TransactionType string `json:"transaction_type"`
}

// operationDTO — элемент блока `operations` ответа эндпоинта.
type operationDTO struct {
	Period            string  `json:"period"`
	Type              string  `json:"type"`
	OperationMetricID string  `json:"operation_metric_id"`
	Unit              string  `json:"unit"`
	OriginalUnit      string  `json:"original_unit"`
	Link              string  `json:"link"`
	LinkUpdate        string  `json:"link_update"`
	Year              int     `json:"year"`
	Month             int     `json:"month"`
	Amount            int64   `json:"amount"`
	OriginalAmount    int64   `json:"original_amount"`
	Value             float64 `json:"value"`
	OriginalValue     float64 `json:"original_value"`
	Curs              float64 `json:"curs"`
}

// ownerDTO — элемент блока `owners` ответа эндпоинта.
type ownerDTO struct {
	ChangedAt string `json:"changed_at"`
	Period    string `json:"period"`
	Type      string `json:"type"`
	Owner     string `json:"owner"`
	// Own приходит строкой-decimal — парсим в float64 на translator-слое.
	Own   string `json:"own"`
	Link  string `json:"link"`
	Year  int    `json:"year"`
	Month int    `json:"month"`
}

// shareDTO — элемент блока `shares` ответа эндпоинта.
type shareDTO struct {
	Code   string `json:"code"`
	Period string `json:"period"`
	Type   string `json:"type"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
	Num    int64  `json:"num"`
}

// translateStockInfo раскладывает info-блок в company.StockInfo.
func translateStockInfo(info *infoDTO) company.StockInfo {
	return company.StockInfo{
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

// translateStockSummary раскладывает summary-блок в company.StockSummary.
func translateStockSummary(s *summaryDTO) company.StockSummary {
	return company.StockSummary{
		Capital:            s.Capital,
		EPS:                s.EPS,
		PEG:                s.PEG,
		PeterLynchTarget:   s.PeterLynchTarget,
		GrahamTarget:       s.GrahamTarget,
		DividendFrequency:  s.DividendFrequency,
		DividendStrike:     s.DividendStrike,
		DividendGrowth:     s.DividendGrowth,
		DividendIndex:      s.DividendIndex,
		DividendYield12M:   s.DividendYield12M,
		DividendYield3Y:    s.DividendYield3Y,
		DividendYield5Y:    s.DividendYield5Y,
		DividendGapLast:    s.DividendGapLast,
		DividendGapAverage: s.DividendGapAverage,
		GrowthRevenue3Y:    s.GrowthRevenue3Y,
		GrowthRevenue5Y:    s.GrowthRevenue5Y,
		GrowthEarnings3Y:   s.GrowthEarnings3Y,
		GrowthEarnings5Y:   s.GrowthEarnings5Y,
		GrowthEBITDA3Y:     s.GrowthEBITDA3Y,
		GrowthEBITDA5Y:     s.GrowthEBITDA5Y,
		GrowthAssets3Y:     s.GrowthAssets3Y,
		GrowthAssets5Y:     s.GrowthAssets5Y,
		GrowthEquity3Y:     s.GrowthEquity3Y,
		GrowthEquity5Y:     s.GrowthEquity5Y,
		GrowthFCF3Y:        s.GrowthFCF3Y,
		GrowthFCF5Y:        s.GrowthFCF5Y,
		GrowthNetDebt3Y:    s.GrowthNetDebt3Y,
		GrowthNetDebt5Y:    s.GrowthNetDebt5Y,
		GrowthOperation3Y:  s.GrowthOperation3Y,
		GrowthOperation5Y:  s.GrowthOperation5Y,
		IdeaBuy:            s.IdeaBuy,
		IdeaHold:           s.IdeaHold,
		IdeaSell:           s.IdeaSell,
		IdeaTarget:         s.IdeaTarget,
		IdeaPotential:      s.IdeaPotential,
		IdeaConsensus:      ideaConsensusByCode[s.IdeaConsensus],
		InsiderConsensus:   insiderConsensusByCode[s.InsiderConsensus],
		ChangedAt:          parseDateTime(s.ChangedAt),
	}
}

// translateStockRatios раскладывает срез ratios-блока в срез StockRatio.
func translateStockRatios(rs []ratioDTO) []company.StockRatio {
	if len(rs) == 0 {
		return nil
	}
	return lo.Map(rs, func(r ratioDTO, _ int) company.StockRatio {
		return company.StockRatio{
			ChangedAt:         parseDateTime(r.ChangedAt),
			Period:            buildPeriod(r.Year, r.Month, r.Period, r.Type),
			Capital:           r.Capital,
			PE:                r.PE,
			PBV:               r.PBV,
			PS:                r.PS,
			PCF:               r.PCF,
			PFCF:              r.PFCF,
			PFFO:              r.PFFO,
			EVS:               r.EVS,
			EVEBITDA:          r.EVEBITDA,
			EVEBIT:            r.EVEBIT,
			DebtEBITDA:        r.DebtEBITDA,
			NetDebtEBITDA:     r.NetDebtEBITDA,
			DebtEquity:        r.DebtEquity,
			DebtRatio:         r.DebtRatio,
			CurrentRatio:      r.CurrentRatio,
			InterestCoverage:  r.InterestCoverage,
			GrossMargin:       r.GrossMargin,
			OperationMargin:   r.OperationMargin,
			EBITDAMargin:      r.EBITDAMargin,
			NetMargin:         r.NetMargin,
			ROS:               r.ROS,
			ROE:               r.ROE,
			ROA:               r.ROA,
			ROIC:              r.ROIC,
			ROCE:              r.ROCE,
			DPR:               r.DPR,
			CapexRevenue:      r.CapexRevenue,
			NetWorkingCapital: r.NetWorkingCapital,
			Active:            r.Active,
		}
	})
}

// translateStockReports раскладывает срез reports-блока в срез StockReport.
// Каждая запись отчёта собирается группами строк (заголовок, P&L,
// денежный поток, баланс, per-share) — иначе одна функция-сборщик
// превышает порог funlen на 95 полях отчёта.
func translateStockReports(rs []reportDTO) []company.StockReport {
	if len(rs) == 0 {
		return nil
	}
	return lo.Map(rs, func(r reportDTO, _ int) company.StockReport {
		out := translateReportHeader(&r)
		applyReportProfitAndLoss(&out, &r)
		applyReportCashFlow(&out, &r)
		applyReportBalance(&out, &r)
		applyReportPerShare(&out, &r)
		return out
	})
}

// translateReportHeader заполняет в StockReport поля периода, валюты,
// масштаба, ссылок и флага предварительной публикации.
func translateReportHeader(r *reportDTO) company.StockReport {
	return company.StockReport{
		ChangedAt:   parseDateTime(r.ChangedAt),
		Link:        r.Link,
		LinkPress:   r.LinkPress,
		LinkUpdate:  r.LinkUpdate,
		Period:      buildPeriod(r.Year, r.Month, r.Period, r.Type),
		Currency:    currencyByCode[r.Curr],
		Amount:      r.Amount,
		Preliminary: r.Preliminary,
	}
}

// applyReportProfitAndLoss заполняет в StockReport строки отчёта о
// прибылях и убытках (включая чистый процентный и комиссионный доход
// для финансовых организаций).
func applyReportProfitAndLoss(out *company.StockReport, r *reportDTO) {
	out.Revenue = r.Revenue
	out.CostOfSales = r.CostOfSales
	out.GrossProfit = r.GrossProfit
	out.SelGenAdmExpenses = r.SelGenAdmExpenses
	out.OperatingIncome = r.OperatingIncome
	out.EBIT = r.EBIT
	out.EBITDA = r.EBITDA
	out.EBITDAAdjusted = r.EBITDAAdjusted
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
	out.FFO = r.FFO
	out.InterestIncome = r.InterestIncome
	out.InterestExpense = r.InterestExpense
	out.InterestNet = r.InterestNet
	out.CommissionIncome = r.CommissionIncome
	out.CommissionExpense = r.CommissionExpense
	out.CommissionNet = r.CommissionNet
}

// applyReportCashFlow заполняет в StockReport строки отчёта о движении
// денежных средств.
func applyReportCashFlow(out *company.StockReport, r *reportDTO) {
	out.CFO = r.CFO
	out.CFI = r.CFI
	out.CFF = r.CFF
	out.FCF = r.FCF
	out.FCFAdjusted = r.FCFAdjusted
	out.NetChangeInCash = r.NetChangeInCash
	out.Capex = r.Capex
	out.PPEPurchase = r.PPEPurchase
	out.IntangiblesPurchase = r.IntangiblesPurchase
	out.RepurchaseOfStock = r.RepurchaseOfStock
	out.IssuanceOfDebt = r.IssuanceOfDebt
	out.PaymentsOfDebt = r.PaymentsOfDebt
	out.NetIssuanceOfDebt = r.NetIssuanceOfDebt
	out.PaymentsForDividends = r.PaymentsForDividends
	out.CashPaidForInterest = r.CashPaidForInterest
	out.CashPaidForTax = r.CashPaidForTax
}

// applyReportBalance заполняет в StockReport строки баланса (активы,
// обязательства, капитал).
func applyReportBalance(out *company.StockReport, r *reportDTO) {
	out.TotalAssets = r.TotalAssets
	out.CurrentAssets = r.CurrentAssets
	out.LongTermAssets = r.LongTermAssets
	out.CashAndEquiv = r.CashAndEquiv
	out.ShortTermInvestments = r.ShortTermInvestments
	out.CashEquivSTInvestments = r.CashEquivSTInvestments
	out.AccountsReceivable = r.AccountsReceivable
	out.OtherReceivable = r.OtherReceivable
	out.TotalReceivable = r.TotalReceivable
	out.Inventories = r.Inventories
	out.PropertyPlantEquipment = r.PropertyPlantEquipment
	out.PPERoU = r.PPERoU
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

// applyReportPerShare заполняет в StockReport per-share метрики.
func applyReportPerShare(out *company.StockReport, r *reportDTO) {
	out.EarningsPS = r.EarningsPS
	out.EquityPS = r.EquityPS
	out.RevenuePS = r.RevenuePS
	out.EBITDAPS = r.EBITDAPS
	out.EBITDAAdjustedPS = r.EBITDAAdjustedPS
	out.FCFPS = r.FCFPS
	out.FCFAdjustedPS = r.FCFAdjustedPS
}

// translateStockDividends раскладывает срез dividends-блока в срез StockDividend.
func translateStockDividends(ds []dividendDTO) []company.StockDividend {
	if len(ds) == 0 {
		return nil
	}
	return lo.Map(ds, func(d dividendDTO, _ int) company.StockDividend {
		return company.StockDividend{
			LastBuyDate:     parseDate(d.LastBuyDate),
			ReestrCloseDate: parseDate(d.ReestrCloseDate),
			ChangedAt:       parseDateTime(d.ChangedAt),
			LastBuyPrice:    d.LastBuyPrice,
			DivAmount:       d.DivAmount,
			DivPercent:      d.DivPercent,
			Year:            d.Year,
			Link:            d.Link,
			Currency:        currencyByCode[d.DivCurr],
			Type:            dividendTypeByCode[d.Type],
		}
	})
}

// translateStockIdeas раскладывает срез ideas-блока в срез StockIdea.
func translateStockIdeas(is []ideaDTO) []company.StockIdea {
	if len(is) == 0 {
		return nil
	}
	return lo.Map(is, func(i ideaDTO, _ int) company.StockIdea {
		return company.StockIdea{
			DateIn:          parseDate(i.DateIn),
			DateOut:         parseDate(i.DateOut),
			CloseDate:       parseDate(i.CloseDate),
			UpdateDate:      parseDate(i.UpdateDate),
			ChangedAt:       parseDateTime(i.ChangedAt),
			Community:       i.Community,
			Idea:            i.Idea,
			CloseComment:    i.CloseComment,
			CloseLink:       i.CloseLink,
			ID:              i.ID,
			CommunityID:     i.CommunityID,
			DurationInMonth: i.DurationInMonth,
			PriceIn:         i.PriceIn,
			PriceOut:        i.PriceOut,
			PriceDay:        i.PriceDay,
			ProfitPotential: i.ProfitPotential,
			ProfitActual:    i.ProfitActual,
			StopLoss:        i.StopLoss,
			ClosePrice:      i.ClosePrice,
			UpdatePrice:     i.UpdatePrice,
			Status:          ideaStatusByCode[i.SystemStatus],
		}
	})
}

// translateStockInsiderTransactions раскладывает срез insiderTransactions-блока
// в срез StockInsiderTransaction.
func translateStockInsiderTransactions(ts []insiderTransactionDTO) []company.StockInsiderTransaction {
	if len(ts) == 0 {
		return nil
	}
	return lo.Map(ts, func(t insiderTransactionDTO, _ int) company.StockInsiderTransaction {
		return company.StockInsiderTransaction{
			TransactionDate: parseDate(t.TransactionDate),
			Insider:         t.Insider,
			InsiderTitle:    t.InsiderTitle,
			Type:            insiderTransactionTypeByCode[t.TransactionType],
		}
	})
}

// translateStockOperations раскладывает срез operations-блока в срез StockOperation.
func translateStockOperations(os []operationDTO) []company.StockOperation {
	if len(os) == 0 {
		return nil
	}
	return lo.Map(os, func(o operationDTO, _ int) company.StockOperation {
		return company.StockOperation{
			MetricID:       o.OperationMetricID,
			Unit:           o.Unit,
			OriginalUnit:   o.OriginalUnit,
			Link:           o.Link,
			LinkUpdate:     o.LinkUpdate,
			Period:         buildPeriod(o.Year, o.Month, o.Period, o.Type),
			Amount:         o.Amount,
			OriginalAmount: o.OriginalAmount,
			Value:          o.Value,
			OriginalValue:  o.OriginalValue,
			Curs:           o.Curs,
		}
	})
}

// translateStockOwners раскладывает срез owners-блока в срез StockOwner.
func translateStockOwners(os []ownerDTO) []company.StockOwner {
	if len(os) == 0 {
		return nil
	}
	return lo.Map(os, func(o ownerDTO, _ int) company.StockOwner {
		return company.StockOwner{
			ChangedAt: parseDateTime(o.ChangedAt),
			Owner:     o.Owner,
			Link:      o.Link,
			Period:    buildPeriod(o.Year, o.Month, o.Period, o.Type),
			Own:       parseDecimal(o.Own),
		}
	})
}

// translateStockShares раскладывает срез shares-блока в срез StockShare.
func translateStockShares(ss []shareDTO) []company.StockShare {
	if len(ss) == 0 {
		return nil
	}
	return lo.Map(ss, func(s shareDTO, _ int) company.StockShare {
		return company.StockShare{
			Ticker: s.Code,
			Period: buildPeriod(s.Year, s.Month, s.Period, s.Type),
			Num:    s.Num,
		}
	})
}

// buildPeriod собирает value-объект StockPeriod из плоских полей записи FM.
func buildPeriod(year, month int, period, standard string) company.StockPeriod {
	return company.StockPeriod{
		Year:      year,
		Month:     month,
		Frequency: stockPeriodFrequencyByCode[period],
		Standard:  reportStandardByCode[standard],
	}
}

// parseDateTime разбирает FinanceMarker datetime без таймзоны. Пустая
// или невалидная строка превращается в zero time.Time.
func parseDateTime(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	t, err := time.Parse(changedAtLayout, raw)
	if err != nil {
		return time.Time{}
	}
	return t
}

// parseDate разбирает FinanceMarker date (YYYY-MM-DD). Пустая или
// невалидная строка превращается в zero time.Time.
func parseDate(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	t, err := time.Parse(dateLayout, raw)
	if err != nil {
		return time.Time{}
	}
	return t
}

// parseDecimal разбирает строку-decimal (FinanceMarker отдаёт own долей
// акционера именно так). Пустая или невалидная строка превращается в 0.
func parseDecimal(raw string) float64 {
	if raw == "" {
		return 0
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return v
}
