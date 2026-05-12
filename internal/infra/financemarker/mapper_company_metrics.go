package financemarker

import (
	"time"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// summaryDTO — блок `summary` ответа /api/fm/v2/stocks/{exchange}:{code}.
//
// FinanceMarker не присылает многие поля для отдельных типов эмитентов
// (банки, фонды) — отсутствующие float-поля разбираются как 0.
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
	DividendYield12m   float64 `json:"dividend_yield_12m"`
	DividendYield3y    float64 `json:"dividend_yield_3y"`
	DividendYield5y    float64 `json:"dividend_yield_5y"`
	GrowthRevenue3y    float64 `json:"growth_revenue_3y"`
	GrowthRevenue5y    float64 `json:"growth_revenue_5y"`
	GrowthEarnings3y   float64 `json:"growth_earnings_3y"`
	GrowthEarnings5y   float64 `json:"growth_earnings_5y"`
	GrowthEbitda3y     float64 `json:"growth_ebitda_3y"`
	GrowthEbitda5y     float64 `json:"growth_ebitda_5y"`
	GrowthAssets3y     float64 `json:"growth_assets_3y"`
	GrowthAssets5y     float64 `json:"growth_assets_5y"`
	GrowthEquity3y     float64 `json:"growth_equity_3y"`
	GrowthEquity5y     float64 `json:"growth_equity_5y"`
	GrowthFCF3y        float64 `json:"growth_fcf_3y"`
	GrowthFCF5y        float64 `json:"growth_fcf_5y"`
	GrowthNetDebt3y    float64 `json:"growth_netdebt_3y"`
	GrowthNetDebt5y    float64 `json:"growth_netdebt_5y"`
	GrowthOperation3y  float64 `json:"growth_operation_3y"`
	GrowthOperation5y  float64 `json:"growth_operation_5y"`
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

// stockDTO — корневой объект ответа эндпоинта по эмитенту.
type stockDTO struct {
	Info    infoDTO    `json:"info"`
	Summary summaryDTO `json:"summary"`
}

// changedAtLayout — формат даты `changed_at` FinanceMarker (ISO 8601 без таймзоны).
const changedAtLayout = "2006-01-02T15:04:05"

// mapCompanyMetrics собирает entities.CompanyMetrics из stockDTO.
//
// ChangedAt разбирается в МСК-таймзоне (биржа работает по этому времени).
// Некорректное значение оставляет нулевой time.Time.
func mapCompanyMetrics(dto *stockDTO) entities.CompanyMetrics {
	return entities.CompanyMetrics{
		Card:               mapCompanyCard(&dto.Info),
		Description:        dto.Info.Description,
		Site:               dto.Info.Site,
		DiscLink:           dto.Info.DiscLink,
		Capital:            dto.Summary.Capital,
		EPS:                dto.Summary.EPS,
		PEG:                dto.Summary.PEG,
		PeterLynchTarget:   dto.Summary.PeterLynchTarget,
		GrahamTarget:       dto.Summary.GrahamTarget,
		DividendFrequency:  dto.Summary.DividendFrequency,
		DividendStrike:     dto.Summary.DividendStrike,
		DividendGrowth:     dto.Summary.DividendGrowth,
		DividendIndex:      dto.Summary.DividendIndex,
		DividendYield12m:   dto.Summary.DividendYield12m,
		DividendYield3y:    dto.Summary.DividendYield3y,
		DividendYield5y:    dto.Summary.DividendYield5y,
		DividendGapLast:    dto.Summary.DividendGapLast,
		DividendGapAverage: dto.Summary.DividendGapAverage,
		GrowthRevenue3y:    dto.Summary.GrowthRevenue3y,
		GrowthRevenue5y:    dto.Summary.GrowthRevenue5y,
		GrowthEarnings3y:   dto.Summary.GrowthEarnings3y,
		GrowthEarnings5y:   dto.Summary.GrowthEarnings5y,
		GrowthEbitda3y:     dto.Summary.GrowthEbitda3y,
		GrowthEbitda5y:     dto.Summary.GrowthEbitda5y,
		GrowthAssets3y:     dto.Summary.GrowthAssets3y,
		GrowthAssets5y:     dto.Summary.GrowthAssets5y,
		GrowthEquity3y:     dto.Summary.GrowthEquity3y,
		GrowthEquity5y:     dto.Summary.GrowthEquity5y,
		GrowthFCF3y:        dto.Summary.GrowthFCF3y,
		GrowthFCF5y:        dto.Summary.GrowthFCF5y,
		GrowthNetDebt3y:    dto.Summary.GrowthNetDebt3y,
		GrowthNetDebt5y:    dto.Summary.GrowthNetDebt5y,
		GrowthOperation3y:  dto.Summary.GrowthOperation3y,
		GrowthOperation5y:  dto.Summary.GrowthOperation5y,
		IdeaBuy:            dto.Summary.IdeaBuy,
		IdeaHold:           dto.Summary.IdeaHold,
		IdeaSell:           dto.Summary.IdeaSell,
		IdeaConsensus:      parseIdeaConsensus(dto.Summary.IdeaConsensus),
		IdeaTarget:         dto.Summary.IdeaTarget,
		IdeaPotential:      dto.Summary.IdeaPotential,
		InsiderConsensus:   parseInsiderConsensus(dto.Summary.InsiderConsensus),
		ChangedAt:          parseChangedAt(dto.Summary.ChangedAt),
	}
}

// parseIdeaConsensus переводит код consensus FinanceMarker в domain-enum.
// Неизвестные значения — IdeaConsensusUnspecified.
func parseIdeaConsensus(s string) entities.IdeaConsensus {
	switch s {
	case "BUY":
		return entities.IdeaConsensusBuy
	case "HOLD":
		return entities.IdeaConsensusHold
	case "SELL":
		return entities.IdeaConsensusSell
	default:
		return entities.IdeaConsensusUnspecified
	}
}

// parseInsiderConsensus переводит перекос сделок инсайдеров в domain-enum.
// Неизвестные значения — InsiderConsensusUnspecified.
func parseInsiderConsensus(s string) entities.InsiderConsensus {
	switch s {
	case "BUYS":
		return entities.InsiderConsensusBuys
	case "SELLS":
		return entities.InsiderConsensusSells
	case "MIXED":
		return entities.InsiderConsensusMixed
	default:
		return entities.InsiderConsensusUnspecified
	}
}

// parseChangedAt разбирает строку changed_at в time.Time. Пустое или
// некорректное значение — нулевой time.Time.
func parseChangedAt(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(changedAtLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}
