package stock

import (
	"errors"
	"time"

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
// ошибка вызывающего: пустой набор секций означает «делать нечего», но
// формально это запрос к источнику ни о чём.
var errEmptyOptions = errors.New("financemarker stock: no sections requested")

// sectionFetch описывает одну запрашиваемую секцию: какой код передать
// в параметре include, какой TTL поставить кешу и как разложить
// результат в company.Stock. Каждый sectionFetch обслуживается ровно
// одним HTTP-запросом, поэтому у каждой секции свой кеш-ключ
// (через сортированный requestKey в httpcache) и свой срок жизни.
type sectionFetch struct {
	apply func(dto *stockDTO) func(*company.Stock)
	code  string
	ttl   time.Duration
}

// enabledSections возвращает список секций, которые нужно запросить под
// данный набор опций. Порядок — канонический (info, summary, ratios,
// reports, dividends, ideas, insiderTransactions, operations, owners,
// shares); совпадает с порядком полей company.Stock и порядком разделов
// справочника api-financemarker/references/stock.md.
func enabledSections(opts company.StockOptions) []sectionFetch {
	out := make([]sectionFetch, 0, 10)
	if opts.WithInfo {
		out = append(out, sectionFetch{code: includeInfo, ttl: ttlInfo, apply: applyInfo})
	}
	if opts.WithSummary {
		out = append(out, sectionFetch{code: includeSummary, ttl: ttlSummary, apply: applySummary})
	}
	if opts.WithRatios {
		out = append(out, sectionFetch{code: includeRatios, ttl: ttlRatios, apply: applyRatios})
	}
	if opts.WithReports {
		out = append(out, sectionFetch{code: includeReports, ttl: ttlReports, apply: applyReports})
	}
	if opts.WithDividends {
		out = append(out, sectionFetch{code: includeDividends, ttl: ttlDividends, apply: applyDividends})
	}
	if opts.WithIdeas {
		out = append(out, sectionFetch{code: includeIdeas, ttl: ttlIdeas, apply: applyIdeas})
	}
	if opts.WithInsiderTransactions {
		out = append(out, sectionFetch{
			code: includeInsiderTransactions, ttl: ttlInsiderTransactions, apply: applyInsiderTransactions,
		})
	}
	if opts.WithOperations {
		out = append(out, sectionFetch{code: includeOperations, ttl: ttlOperations, apply: applyOperations})
	}
	if opts.WithOwners {
		out = append(out, sectionFetch{code: includeOwners, ttl: ttlOwners, apply: applyOwners})
	}
	if opts.WithShares {
		out = append(out, sectionFetch{code: includeShares, ttl: ttlShares, apply: applyShares})
	}
	return out
}

// applyInfo — apply-функция секции info: парсит секцию из DTO и
// возвращает аппликатор, который запишет её в company.Stock.
func applyInfo(d *stockDTO) func(*company.Stock) {
	v := translateStockInfo(&d.Info)
	return func(s *company.Stock) { s.Info = v }
}

// applySummary — apply-функция секции summary.
func applySummary(d *stockDTO) func(*company.Stock) {
	v := translateStockSummary(&d.Summary)
	return func(s *company.Stock) { s.Summary = v }
}

// applyRatios — apply-функция секции ratios.
func applyRatios(d *stockDTO) func(*company.Stock) {
	v := translateStockRatios(d.Ratios)
	return func(s *company.Stock) { s.Ratios = v }
}

// applyReports — apply-функция секции reports.
func applyReports(d *stockDTO) func(*company.Stock) {
	v := translateStockReports(d.Reports)
	return func(s *company.Stock) { s.Reports = v }
}

// applyDividends — apply-функция секции dividends.
func applyDividends(d *stockDTO) func(*company.Stock) {
	v := translateStockDividends(d.Dividends)
	return func(s *company.Stock) { s.Dividends = v }
}

// applyIdeas — apply-функция секции ideas.
func applyIdeas(d *stockDTO) func(*company.Stock) {
	v := translateStockIdeas(d.Ideas)
	return func(s *company.Stock) { s.Ideas = v }
}

// applyInsiderTransactions — apply-функция секции insiderTransactions.
func applyInsiderTransactions(d *stockDTO) func(*company.Stock) {
	v := translateStockInsiderTransactions(d.InsiderTransactions)
	return func(s *company.Stock) { s.InsiderTransactions = v }
}

// applyOperations — apply-функция секции operations.
func applyOperations(d *stockDTO) func(*company.Stock) {
	v := translateStockOperations(d.Operations)
	return func(s *company.Stock) { s.Operations = v }
}

// applyOwners — apply-функция секции owners.
func applyOwners(d *stockDTO) func(*company.Stock) {
	v := translateStockOwners(d.Owners)
	return func(s *company.Stock) { s.Owners = v }
}

// applyShares — apply-функция секции shares.
func applyShares(d *stockDTO) func(*company.Stock) {
	v := translateStockShares(d.Shares)
	return func(s *company.Stock) { s.Shares = v }
}
