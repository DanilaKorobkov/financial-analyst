package stocksummary

import (
	"time"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

// changedAtLayout — формат datetime, в котором FinanceMarker отдаёт
// поле changed_at у блока summary (ISO 8601 без таймзоны, секунды).
const changedAtLayout = "2006-01-02T15:04:05"

var (
	// ideaConsensusByCode переводит строковый код консенсуса инвест-идей
	// FinanceMarker в domain-enum. Незнакомые значения возвращают zero
	// (IdeaConsensusUnspecified).
	ideaConsensusByCode = map[string]company.IdeaConsensus{
		"BUY":  company.IdeaConsensusBuy,
		"HOLD": company.IdeaConsensusHold,
		"SELL": company.IdeaConsensusSell,
	}

	// insiderConsensusByCode переводит строковый код консенсуса сделок
	// инсайдеров FinanceMarker в domain-enum. Незнакомые значения возвращают
	// zero (InsiderConsensusUnspecified).
	insiderConsensusByCode = map[string]company.InsiderConsensus{
		"BUYS":  company.InsiderConsensusBuys,
		"SELLS": company.InsiderConsensusSells,
		"MIXED": company.InsiderConsensusMixed,
	}
)

// summaryDTO — блок `summary` ответа /api/fm/v2/stocks/{exchange}:{code}.
// Поля, не отданные источником, остаются zero-value у DTO и далее
// переезжают как zero-value в company.StockSummary.
//
// Порядок полей — под fieldalignment: string-поля (с указателем) идут
// первыми, далее плоские числовые группы.
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

// stockDTO — корневой объект ответа эндпоинта по эмитенту. Здесь
// разбирается только блок summary — остальные разделы (info / ratios / ...)
// принадлежат соседним источникам.
type stockDTO struct {
	Summary summaryDTO `json:"summary"`
}

// translateStockSummary раскладывает summary-блок FinanceMarker
// в StockSummary. Невалидный или пустой changed_at превращается
// в zero time.Time.
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
		ChangedAt:          parseChangedAt(s.ChangedAt),
	}
}

// parseChangedAt разбирает FinanceMarker datetime без таймзоны. Пустая
// или невалидная строка превращается в zero time.Time.
func parseChangedAt(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	t, err := time.Parse(changedAtLayout, raw)
	if err != nil {
		return time.Time{}
	}
	return t
}
