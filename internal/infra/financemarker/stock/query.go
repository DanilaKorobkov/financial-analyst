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

	// includeRatios — код блока ratios.
	includeRatios = "ratios"

	// includeReports — код блока reports.
	includeReports = "reports"

	// includeDividends — код блока dividends.
	includeDividends = "dividends"

	// includeIdeas — код блока ideas.
	includeIdeas = "ideas"

	// includeInsiderTransactions — код блока insiderTransactions (camelCase
	// сохранён как в FM, расхождение со snake_case остальных блоков —
	// требование контракта эндпоинта).
	includeInsiderTransactions = "insiderTransactions"

	// includeOperations — код блока operations.
	includeOperations = "operations"

	// includeOwners — код блока owners.
	includeOwners = "owners"

	// includeShares — код блока shares.
	includeShares = "shares"
)

// errEmptyOptions — вызов источника без запрошенных секций. Контрактная
// ошибка вызывающего: пустой include заставит FinanceMarker вернуть
// заполненный только info-блок и потратить тариф впустую.
var errEmptyOptions = errors.New("financemarker stock: no sections requested")

// buildInclude собирает значение query-параметра include в каноническом
// порядке (info, summary, ratios, reports, dividends, ideas,
// insiderTransactions, operations, owners, shares) — соответствует
// порядку полей company.Stock и порядку разделов справочника
// api-financemarker/references/stock.md. Один use case → одна строка
// include — условие устойчивого ключа кеша.
func buildInclude(opts company.StockOptions) (string, error) {
	codes := make([]string, 0, 10)
	if opts.WithInfo {
		codes = append(codes, includeInfo)
	}
	if opts.WithSummary {
		codes = append(codes, includeSummary)
	}
	if opts.WithRatios {
		codes = append(codes, includeRatios)
	}
	if opts.WithReports {
		codes = append(codes, includeReports)
	}
	if opts.WithDividends {
		codes = append(codes, includeDividends)
	}
	if opts.WithIdeas {
		codes = append(codes, includeIdeas)
	}
	if opts.WithInsiderTransactions {
		codes = append(codes, includeInsiderTransactions)
	}
	if opts.WithOperations {
		codes = append(codes, includeOperations)
	}
	if opts.WithOwners {
		codes = append(codes, includeOwners)
	}
	if opts.WithShares {
		codes = append(codes, includeShares)
	}
	if len(codes) == 0 {
		return "", errEmptyOptions
	}
	return strings.Join(codes, ","), nil
}
