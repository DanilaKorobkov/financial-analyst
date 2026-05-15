package company

import "time"

// StockReport — публикация финансовой отчётности эмитента за один период
// в одном стандарте (МСФО/РСБУ/GAAP). Содержит строки отчёта о прибылях,
// баланса, движения денежных средств и per-share метрики; нулевое значение
// поля = «не отдано источником» (FM возвращает `null` для строк, не
// применимых к эмитенту — банкам, REIT и т.п.).
//
// Все денежные значения даны в исходных единицах с множителем Amount
// (1 = единицы валюты отчётности, 1_000_000 = млн и т.д.).
type StockReport struct {
	// ChangedAt — момент записи.
	ChangedAt time.Time

	// Link — ссылка на исходный отчёт эмитента.
	Link string

	// LinkPress — ссылка на пресс-релиз.
	LinkPress string

	// LinkUpdate — ссылка на обновление отчёта.
	LinkUpdate string

	// Period — отчётный период публикации.
	Period StockPeriod

	// Currency — валюта отчётности эмитента.
	Currency Currency

	// Amount — множитель денежных значений (1, 1_000_000, ...).
	Amount int64

	// Revenue — выручка.
	Revenue float64

	// CostOfSales — себестоимость.
	CostOfSales float64

	// GrossProfit — валовая прибыль.
	GrossProfit float64

	// SelGenAdmExpenses — коммерческие, общехозяйственные и административные расходы.
	SelGenAdmExpenses float64

	// OperatingIncome — операционная прибыль.
	OperatingIncome float64

	// EBIT — прибыль до уплаты процентов и налога на прибыль.
	EBIT float64

	// EBITDA — прибыль до уплаты процентов, налога на прибыль и амортизации.
	EBITDA float64

	// EBITDAAdjusted — скорректированная EBITDA.
	EBITDAAdjusted float64

	// Earnings — чистая прибыль.
	Earnings float64

	// EarningsWoTax — прибыль до налога на прибыль.
	EarningsWoTax float64

	// EarningsComprehensive — совокупный доход.
	EarningsComprehensive float64

	// EarningsComprehensiveStockHolders — совокупный доход акционеров материнской компании.
	EarningsComprehensiveStockHolders float64

	// EarningsStockHolders — чистая прибыль акционеров материнской компании.
	EarningsStockHolders float64

	// EarningsContinuingOperations — чистая прибыль от продолжающейся деятельности.
	EarningsContinuingOperations float64

	// TotalExpenses — суммарные расходы.
	TotalExpenses float64

	// OtherOperatingIncome — прочие операционные доходы.
	OtherOperatingIncome float64

	// DeprDeplAmort — амортизация и истощение.
	DeprDeplAmort float64

	// ResearchAndDevelopment — расходы на R&D.
	ResearchAndDevelopment float64

	// FFO — Funds From Operations (для REIT).
	FFO float64

	// InterestIncome — процентные доходы (для финансовых организаций).
	InterestIncome float64

	// InterestExpense — процентные расходы.
	InterestExpense float64

	// InterestNet — чистый процентный доход.
	InterestNet float64

	// CommissionIncome — комиссионные доходы.
	CommissionIncome float64

	// CommissionExpense — комиссионные расходы.
	CommissionExpense float64

	// CommissionNet — чистый комиссионный доход.
	CommissionNet float64

	// CFO — денежный поток от операционной деятельности.
	CFO float64

	// CFI — денежный поток от инвестиционной деятельности.
	CFI float64

	// CFF — денежный поток от финансовой деятельности.
	CFF float64

	// FCF — свободный денежный поток.
	FCF float64

	// FCFAdjusted — скорректированный свободный денежный поток.
	FCFAdjusted float64

	// NetChangeInCash — чистое изменение денежных средств за период.
	NetChangeInCash float64

	// Capex — капитальные затраты.
	Capex float64

	// PPEPurchase — приобретение основных средств.
	PPEPurchase float64

	// IntangiblesPurchase — приобретение нематериальных активов.
	IntangiblesPurchase float64

	// RepurchaseOfStock — выкуп собственных акций.
	RepurchaseOfStock float64

	// IssuanceOfDebt — привлечение долга.
	IssuanceOfDebt float64

	// PaymentsOfDebt — погашение долга.
	PaymentsOfDebt float64

	// NetIssuanceOfDebt — чистое привлечение долга.
	NetIssuanceOfDebt float64

	// PaymentsForDividends — выплаченные дивиденды.
	PaymentsForDividends float64

	// CashPaidForInterest — выплаченные проценты.
	CashPaidForInterest float64

	// CashPaidForTax — выплаченный налог на прибыль.
	CashPaidForTax float64

	// TotalAssets — суммарные активы.
	TotalAssets float64

	// CurrentAssets — оборотные активы.
	CurrentAssets float64

	// LongTermAssets — внеоборотные активы.
	LongTermAssets float64

	// CashAndEquiv — денежные средства и эквиваленты.
	CashAndEquiv float64

	// ShortTermInvestments — краткосрочные финансовые вложения.
	ShortTermInvestments float64

	// CashEquivSTInvestments — денежные средства, эквиваленты и краткосрочные вложения.
	CashEquivSTInvestments float64

	// AccountsReceivable — торговая дебиторская задолженность.
	AccountsReceivable float64

	// OtherReceivable — прочая дебиторская задолженность.
	OtherReceivable float64

	// TotalReceivable — суммарная дебиторская задолженность.
	TotalReceivable float64

	// Inventories — запасы.
	Inventories float64

	// PropertyPlantEquipment — основные средства.
	PropertyPlantEquipment float64

	// PPERoU — основные средства, включая активы в форме права пользования.
	PPERoU float64

	// RightOfUseAssets — активы в форме права пользования.
	RightOfUseAssets float64

	// IntangibleAssets — нематериальные активы.
	IntangibleAssets float64

	// Goodwill — гудвилл.
	Goodwill float64

	// GoodwillIntangibleAssets — гудвилл и нематериальные активы.
	GoodwillIntangibleAssets float64

	// IntangibleAndTangibleAssets — нематериальные и материальные активы.
	IntangibleAndTangibleAssets float64

	// LongTermInvestments — долгосрочные финансовые вложения.
	LongTermInvestments float64

	// TotalLiabilities — суммарные обязательства.
	TotalLiabilities float64

	// CurrentLiabilities — краткосрочные обязательства.
	CurrentLiabilities float64

	// LongTermLiabilities — долгосрочные обязательства.
	LongTermLiabilities float64

	// CurrentDebt — краткосрочный долг.
	CurrentDebt float64

	// LongTermDebt — долгосрочный долг.
	LongTermDebt float64

	// TotalDebt — суммарный долг.
	TotalDebt float64

	// NetDebt — чистый долг.
	NetDebt float64

	// NetDebtAdjusted — скорректированный чистый долг.
	NetDebtAdjusted float64

	// CurLongDebt — текущая часть долгосрочного долга.
	CurLongDebt float64

	// CurLongLease — текущая часть долгосрочной аренды.
	CurLongLease float64

	// CurrentLease — краткосрочные обязательства по аренде.
	CurrentLease float64

	// LongTermLease — долгосрочные обязательства по аренде.
	LongTermLease float64

	// AccountsPayable — торговая кредиторская задолженность.
	AccountsPayable float64

	// OtherPayable — прочая кредиторская задолженность.
	OtherPayable float64

	// TotalPayable — суммарная кредиторская задолженность.
	TotalPayable float64

	// Equity — собственный капитал.
	Equity float64

	// EquityStockHolders — собственный капитал акционеров материнской компании.
	EquityStockHolders float64

	// RetainedEarnings — нераспределённая прибыль.
	RetainedEarnings float64

	// SharePremium — эмиссионный доход.
	SharePremium float64

	// TreasuryStock — выкупленные собственные акции.
	TreasuryStock float64

	// EarningsPS — прибыль на акцию.
	EarningsPS float64

	// EquityPS — собственный капитал на акцию.
	EquityPS float64

	// RevenuePS — выручка на акцию.
	RevenuePS float64

	// EBITDAPS — EBITDA на акцию.
	EBITDAPS float64

	// EBITDAAdjustedPS — скорректированная EBITDA на акцию.
	EBITDAAdjustedPS float64

	// FCFPS — свободный денежный поток на акцию.
	FCFPS float64

	// FCFAdjustedPS — скорректированный свободный денежный поток на акцию.
	FCFAdjustedPS float64

	// Preliminary — true для предварительной публикации отчёта.
	Preliminary bool
}
