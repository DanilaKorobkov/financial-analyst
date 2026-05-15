package company

import "context"

// StockOptions — набор секций, которые нужно вернуть в Stock. Каждый
// флаг отвечает за свою секцию: запрошены только те, что включены.
type StockOptions struct {
	// WithInfo — запросить карточку эмитента (Info).
	WithInfo bool

	// WithSummary — запросить сводные метрики эмитента (Summary).
	WithSummary bool
}

// Stock — секции карточки эмитента, возвращаемые StockSource в одном
// вызове. Поля, не запрошенные у источника, остаются zero-value.
type Stock struct {
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
