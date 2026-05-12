package entities

import "time"

// CompanyMetrics — расширенная карточка эмитента: лёгкая шапка CompanyCard
// плюс описание, сводные метрики (EPS, PEG, target-цены моделей Грэма и
// Линча), статистика дивидендов, темпы роста, консенсус инвестидей и
// инсайдерских сделок.
//
// Поля float64 — «как пришло», без округления. Отсутствующие у источника
// значения остаются нулями: domain не различает «ноль» и «нет данных»
// для скалярных метрик. Признаки «нет данных» по enum-полям выражены
// явным значением *Unspecified.
type CompanyMetrics struct {
	// ChangedAt — время последнего пересчёта сводных метрик на стороне источника.
	ChangedAt time.Time

	// Description — описание эмитента в свободной форме.
	Description string

	// Site — корпоративный сайт эмитента.
	Site string

	// DiscLink — ссылка на страницу раскрытия информации (investor relations).
	DiscLink string

	// Card — лёгкая карточка эмитента (идентификация и классификация).
	Card CompanyCard

	// Capital — рыночная капитализация, в валюте Card.Currency, в миллионах.
	Capital float64

	// EPS — прибыль на акцию (TTM), в валюте отчётности.
	EPS float64

	// PEG — P/E, делённое на ожидаемый темп роста прибыли.
	PEG float64

	// PeterLynchTarget — целевая цена по модели Питера Линча (PEG-based).
	PeterLynchTarget float64

	// GrahamTarget — целевая цена по модели Бенджамина Грэма.
	GrahamTarget float64

	// DividendIndex — композитный индекс дивидендной привлекательности (0..~1).
	DividendIndex float64

	// DividendYield12m — дивдоходность за последние 12 месяцев (%).
	DividendYield12m float64

	// DividendYield3y — средняя дивдоходность за 3 года (%).
	DividendYield3y float64

	// DividendYield5y — средняя дивдоходность за 5 лет (%).
	DividendYield5y float64

	// GrowthRevenue3y — CAGR выручки за 3 года (%).
	GrowthRevenue3y float64

	// GrowthRevenue5y — CAGR выручки за 5 лет (%).
	GrowthRevenue5y float64

	// GrowthEarnings3y — CAGR чистой прибыли за 3 года (%).
	GrowthEarnings3y float64

	// GrowthEarnings5y — CAGR чистой прибыли за 5 лет (%).
	GrowthEarnings5y float64

	// GrowthEbitda3y — CAGR EBITDA за 3 года (%). У финкомпаний обычно 0.
	GrowthEbitda3y float64

	// GrowthEbitda5y — CAGR EBITDA за 5 лет (%). У финкомпаний обычно 0.
	GrowthEbitda5y float64

	// GrowthAssets3y — CAGR активов за 3 года (%).
	GrowthAssets3y float64

	// GrowthAssets5y — CAGR активов за 5 лет (%).
	GrowthAssets5y float64

	// GrowthEquity3y — CAGR собственного капитала за 3 года (%).
	GrowthEquity3y float64

	// GrowthEquity5y — CAGR собственного капитала за 5 лет (%).
	GrowthEquity5y float64

	// GrowthFCF3y — CAGR свободного денежного потока за 3 года (%).
	GrowthFCF3y float64

	// GrowthFCF5y — CAGR свободного денежного потока за 5 лет (%).
	GrowthFCF5y float64

	// GrowthNetDebt3y — CAGR чистого долга за 3 года (%). У финкомпаний обычно 0.
	GrowthNetDebt3y float64

	// GrowthNetDebt5y — CAGR чистого долга за 5 лет (%). У финкомпаний обычно 0.
	GrowthNetDebt5y float64

	// GrowthOperation3y — CAGR операционной прибыли за 3 года (%).
	GrowthOperation3y float64

	// GrowthOperation5y — CAGR операционной прибыли за 5 лет (%).
	GrowthOperation5y float64

	// IdeaTarget — усреднённая целевая цена по активным инвестидеям.
	IdeaTarget float64

	// IdeaPotential — средний потенциал к target по активным инвестидеям (%).
	IdeaPotential float64

	// DividendFrequency — историческая частота выплаты дивидендов в году
	// (1 — годовые, 2 — полугодовые, 4 — квартальные).
	DividendFrequency int

	// DividendStrike — сколько лет подряд эмитент не пропускал дивиденд.
	DividendStrike int

	// DividendGrowth — сколько лет подряд дивиденд рос без снижений.
	DividendGrowth int

	// DividendGapLast — длительность последнего дивидендного gap в днях.
	DividendGapLast int

	// DividendGapAverage — средняя длительность дивидендного gap в днях.
	DividendGapAverage int

	// IdeaBuy — число активных инвестидей с рекомендацией BUY.
	IdeaBuy int

	// IdeaHold — число активных инвестидей с рекомендацией HOLD.
	IdeaHold int

	// IdeaSell — число активных инвестидей с рекомендацией SELL.
	IdeaSell int

	// IdeaConsensus — консенсус активных инвестидей.
	IdeaConsensus IdeaConsensus

	// InsiderConsensus — перекос сделок инсайдеров.
	InsiderConsensus InsiderConsensus
}

// IdeaConsensus — агрегированный консенсус активных инвестидей по бумаге.
type IdeaConsensus int

const (
	// IdeaConsensusUnspecified — консенсус не определён или неизвестен.
	IdeaConsensusUnspecified IdeaConsensus = iota
	// IdeaConsensusBuy — преобладают рекомендации BUY.
	IdeaConsensusBuy
	// IdeaConsensusHold — преобладают рекомендации HOLD.
	IdeaConsensusHold
	// IdeaConsensusSell — преобладают рекомендации SELL.
	IdeaConsensusSell
)

// InsiderConsensus — агрегированный перекос сделок инсайдеров по бумаге.
type InsiderConsensus int

const (
	// InsiderConsensusUnspecified — перекос не определён или неизвестен.
	InsiderConsensusUnspecified InsiderConsensus = iota
	// InsiderConsensusBuys — преобладают покупки инсайдеров.
	InsiderConsensusBuys
	// InsiderConsensusSells — преобладают продажи инсайдеров.
	InsiderConsensusSells
	// InsiderConsensusMixed — покупки и продажи примерно равны.
	InsiderConsensusMixed
)
