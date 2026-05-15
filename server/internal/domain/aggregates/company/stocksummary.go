package company

import "time"

// StockSummary — сводные метрики эмитента «одной строкой»: рыночная
// капитализация, EPS, целевые цены моделей Грэма и Линча, дивидендная
// статистика, среднегодовые темпы роста (CAGR) и агрегированные
// консенсусы по идеям и сделкам инсайдеров. Нулевое значение поля
// означает, что метрика не была отдана источником.
//
// Порядок полей подобран под fieldalignment (time.Time с указателем
// — впереди, далее плоские числовые группы, в конце enum-ы).
type StockSummary struct {
	// ChangedAt — момент пересчёта сводных метрик.
	ChangedAt time.Time

	// Capital — рыночная капитализация эмитента в валюте отчётности
	// (обычно RUB), в миллионах.
	Capital float64

	// EPS — прибыль на акцию (TTM) в валюте отчётности.
	EPS float64

	// PEG — отношение P/E к ожидаемому темпу роста прибыли.
	PEG float64

	// PeterLynchTarget — целевая цена по модели Питера Линча (PEG-based).
	PeterLynchTarget float64

	// GrahamTarget — целевая цена по модели Бенджамина Грэма.
	GrahamTarget float64

	// DividendIndex — композитный индекс дивидендной привлекательности
	// (от 0 до ~1).
	DividendIndex float64

	// DividendYield12M — дивидендная доходность за последние 12 месяцев (%).
	DividendYield12M float64

	// DividendYield3Y — средняя дивидендная доходность за 3 года (%).
	DividendYield3Y float64

	// DividendYield5Y — средняя дивидендная доходность за 5 лет (%).
	DividendYield5Y float64

	// GrowthRevenue3Y — CAGR выручки за 3 года (%).
	GrowthRevenue3Y float64

	// GrowthRevenue5Y — CAGR выручки за 5 лет (%).
	GrowthRevenue5Y float64

	// GrowthEarnings3Y — CAGR чистой прибыли за 3 года (%).
	GrowthEarnings3Y float64

	// GrowthEarnings5Y — CAGR чистой прибыли за 5 лет (%).
	GrowthEarnings5Y float64

	// GrowthEBITDA3Y — CAGR EBITDA за 3 года (%).
	GrowthEBITDA3Y float64

	// GrowthEBITDA5Y — CAGR EBITDA за 5 лет (%).
	GrowthEBITDA5Y float64

	// GrowthAssets3Y — CAGR активов за 3 года (%).
	GrowthAssets3Y float64

	// GrowthAssets5Y — CAGR активов за 5 лет (%).
	GrowthAssets5Y float64

	// GrowthEquity3Y — CAGR собственного капитала за 3 года (%).
	GrowthEquity3Y float64

	// GrowthEquity5Y — CAGR собственного капитала за 5 лет (%).
	GrowthEquity5Y float64

	// GrowthFCF3Y — CAGR свободного денежного потока за 3 года (%).
	GrowthFCF3Y float64

	// GrowthFCF5Y — CAGR свободного денежного потока за 5 лет (%).
	GrowthFCF5Y float64

	// GrowthNetDebt3Y — CAGR чистого долга за 3 года (%).
	GrowthNetDebt3Y float64

	// GrowthNetDebt5Y — CAGR чистого долга за 5 лет (%).
	GrowthNetDebt5Y float64

	// GrowthOperation3Y — CAGR операционной прибыли за 3 года (%).
	GrowthOperation3Y float64

	// GrowthOperation5Y — CAGR операционной прибыли за 5 лет (%).
	GrowthOperation5Y float64

	// IdeaTarget — усреднённая целевая цена по активным идеям.
	IdeaTarget float64

	// IdeaPotential — средний потенциал к целевой цене по активным идеям (%).
	IdeaPotential float64

	// DividendFrequency — историческое число дивидендных выплат в год
	// (1 — годовые, 2 — полугодовые, 4 — квартальные).
	DividendFrequency int

	// DividendStrike — число лет подряд, в течение которых эмитент
	// не пропускал дивидендную выплату.
	DividendStrike int

	// DividendGrowth — число лет подряд, в течение которых дивиденд рос.
	DividendGrowth int

	// DividendGapLast — размер последнего дивидендного гэпа в днях
	// до закрытия.
	DividendGapLast int

	// DividendGapAverage — средний размер дивидендного гэпа в днях
	// до закрытия.
	DividendGapAverage int

	// IdeaBuy — число активных идей с рекомендацией «покупать».
	IdeaBuy int

	// IdeaHold — число активных идей с рекомендацией «держать».
	IdeaHold int

	// IdeaSell — число активных идей с рекомендацией «продавать».
	IdeaSell int

	// IdeaConsensus — агрегированный консенсус по активным идеям.
	IdeaConsensus IdeaConsensus

	// InsiderConsensus — агрегированное направление сделок инсайдеров.
	InsiderConsensus InsiderConsensus
}
