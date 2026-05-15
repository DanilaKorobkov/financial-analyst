package company

// StockOperation — значение одной операционной метрики эмитента за один
// отчётный период. Метрики (`car_loans`, `gmv_incl_services`, ...)
// идентифицируются строковым MetricID; расшифровка — отдельная справочная
// лента FM operation_metrics, в этот агрегат не входит.
type StockOperation struct {
	// MetricID — идентификатор операционной метрики.
	MetricID string

	// Unit — единица измерения итогового Value.
	Unit string

	// OriginalUnit — единица, в которой метрика опубликована эмитентом.
	OriginalUnit string

	// Link — ссылка на исходный отчёт эмитента.
	Link string

	// LinkUpdate — ссылка на обновление отчёта.
	LinkUpdate string

	// Period — отчётный период записи.
	Period StockPeriod

	// Amount — множитель итогового Value, применённый при пересчёте.
	Amount int64

	// OriginalAmount — множитель исходного значения у эмитента.
	OriginalAmount int64

	// Value — значение метрики в единицах Unit с учётом Amount.
	Value float64

	// OriginalValue — значение метрики «как опубликовал эмитент».
	OriginalValue float64

	// Curs — курс пересчёта, если метрика в валюте отличной от рубля.
	Curs float64
}
