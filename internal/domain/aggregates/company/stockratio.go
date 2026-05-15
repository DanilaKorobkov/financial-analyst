package company

import "time"

// StockRatio — мультипликаторы и финансовые коэффициенты эмитента за один
// отчётный период. Нулевое значение поля означает, что метрика не была
// отдана источником (FM отдаёт `null` для мультипликаторов, не применимых
// к эмитенту — например, gross_margin/EBITDA-margin у банка).
type StockRatio struct {
	// ChangedAt — момент пересчёта строки.
	ChangedAt time.Time

	// Period — отчётный период записи.
	Period StockPeriod

	// Capital — рыночная капитализация на конец периода (в млн валюты отчётности).
	Capital float64

	// PE — Price / Earnings (P/E, TTM).
	PE float64

	// PBV — Price / Book Value (P/B).
	PBV float64

	// PS — Price / Sales (P/S).
	PS float64

	// PCF — Price / Operating Cash Flow (P/CF).
	PCF float64

	// PFCF — Price / Free Cash Flow (P/FCF).
	PFCF float64

	// PFFO — Price / FFO (для REIT/недвижимости).
	PFFO float64

	// EVS — EV / Sales.
	EVS float64

	// EVEBITDA — EV / EBITDA.
	EVEBITDA float64

	// EVEBIT — EV / EBIT.
	EVEBIT float64

	// DebtEBITDA — Debt / EBITDA.
	DebtEBITDA float64

	// NetDebtEBITDA — NetDebt / EBITDA.
	NetDebtEBITDA float64

	// DebtEquity — Debt / Equity (D/E).
	DebtEquity float64

	// DebtRatio — Total Debt / Total Assets (доля долга в активах, 0..1).
	DebtRatio float64

	// CurrentRatio — текущая ликвидность (current assets / current liabilities).
	CurrentRatio float64

	// InterestCoverage — покрытие процентов (EBIT / Interest Expense).
	InterestCoverage float64

	// GrossMargin — валовая маржа (%).
	GrossMargin float64

	// OperationMargin — операционная маржа (%).
	OperationMargin float64

	// EBITDAMargin — EBITDA-маржа (%).
	EBITDAMargin float64

	// NetMargin — чистая маржа (%).
	NetMargin float64

	// ROS — Return on Sales (%).
	ROS float64

	// ROE — Return on Equity (%).
	ROE float64

	// ROA — Return on Assets (%).
	ROA float64

	// ROIC — Return on Invested Capital (%).
	ROIC float64

	// ROCE — Return on Capital Employed (%).
	ROCE float64

	// DPR — Dividend Payout Ratio (Dividends / Earnings).
	DPR float64

	// CapexRevenue — CAPEX / Revenue (%).
	CapexRevenue float64

	// NetWorkingCapital — чистый оборотный капитал (в валюте отчётности).
	NetWorkingCapital float64

	// Active — true для последней (актуальной) строки ряда.
	Active bool
}
