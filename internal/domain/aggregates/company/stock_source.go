package company

import "context"

const (
	// StockSectionUnspecified — нулевое значение (не используется как запрос).
	StockSectionUnspecified StockSection = iota

	// StockSectionInfo — карточка эмитента (классификация, описание, ссылки).
	StockSectionInfo

	// StockSectionSummary — сводные метрики эмитента «одной строкой».
	StockSectionSummary
)

// StockSection — блок данных карточки эмитента, который источник может
// заполнить в одном вызове. Перечисление определяет, какие именно секции
// возвращает StockSource.
type StockSection int

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
	FindByTicker(ctx context.Context, ticker string, sections []StockSection) (Stock, error)
}
