package stock

import (
	"errors"
	"strings"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

const (
	// includeInfo — код блока info в query-параметре include.
	includeInfo = "info"

	// includeSummary — код блока summary.
	includeSummary = "summary"
)

// errEmptyOptions — вызов источника без запрошенных секций. Контрактная
// ошибка вызывающего: пустой include заставит FinanceMarker вернуть
// заполненный только info-блок и потратить тариф впустую.
var errEmptyOptions = errors.New("financemarker stock: no sections requested")

// buildInclude собирает значение query-параметра include в каноническом
// порядке (info, summary) — соответствует порядку полей company.Stock и
// порядку разделов справочника api-financemarker/references/stock.md.
// Один use case → одна строка include — условие устойчивого ключа кеша.
func buildInclude(opts company.StockOptions) (string, error) {
	codes := make([]string, 0, 2)
	if opts.WithInfo {
		codes = append(codes, includeInfo)
	}
	if opts.WithSummary {
		codes = append(codes, includeSummary)
	}
	if len(codes) == 0 {
		return "", errEmptyOptions
	}
	return strings.Join(codes, ","), nil
}
