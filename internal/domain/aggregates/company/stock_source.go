package company

import "context"

// StockOptions — набор секций, которые нужно вернуть в Stock. Каждый
// флаг отвечает за свою секцию: запрошены только те, что включены.
type StockOptions struct {
	// WithInfo — запросить карточку эмитента (Info).
	WithInfo bool

	// WithSummary — запросить сводные метрики эмитента (Summary).
	WithSummary bool

	// WithRatios — запросить временной ряд мультипликаторов (Ratios).
	WithRatios bool

	// WithReports — запросить ленту финансовой отчётности (Reports).
	WithReports bool

	// WithDividends — запросить историю и прогноз дивидендов (Dividends).
	WithDividends bool

	// WithIdeas — запросить ленту инвест-идей аналитиков (Ideas).
	WithIdeas bool

	// WithInsiderTransactions — запросить сделки инсайдеров (InsiderTransactions).
	WithInsiderTransactions bool

	// WithOperations — запросить ленту операционных метрик (Operations).
	WithOperations bool

	// WithOwners — запросить структуру акционеров (Owners).
	WithOwners bool

	// WithShares — запросить ряд по количеству выпущенных акций (Shares).
	WithShares bool
}

// Stock — секции карточки эмитента, возвращаемые StockSource в одном
// вызове. Поля, не запрошенные у источника, остаются zero-value
// (массивы — nil, объекты — нулевыми).
type Stock struct {
	// Ratios — временной ряд мультипликаторов по отчётным периодам.
	Ratios []StockRatio

	// Reports — публикации финансовой отчётности.
	Reports []StockReport

	// Dividends — история и прогноз дивидендных выплат.
	Dividends []StockDividend

	// Ideas — инвест-идеи аналитиков по бумаге.
	Ideas []StockIdea

	// InsiderTransactions — сделки инсайдеров в узком формате.
	InsiderTransactions []StockInsiderTransaction

	// Operations — значения операционных метрик по периодам.
	Operations []StockOperation

	// Owners — структура акционеров по датам среза.
	Owners []StockOwner

	// Shares — количество выпущенных акций по датам среза.
	Shares []StockShare

	// Info — карточка эмитента.
	Info StockInfo

	// Summary — сводные метрики эмитента.
	Summary StockSummary
}

// StockSource — порт источника секций карточки эмитента.
//
// Источник делает один запрос за указанные секции и возвращает Stock,
// в котором заполнены только запрошенные секции. Канонизация и порядок
// секций — забота реализации.
type StockSource interface {
	// FindByTicker возвращает запрошенные секции карточки эмитента.
	// Возвращает ErrNotFound, если источник не знает тикер.
	FindByTicker(ctx context.Context, ticker string, opts StockOptions) (Stock, error)
}
