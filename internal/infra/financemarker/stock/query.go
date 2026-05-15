package stock

import (
	"errors"
	"strings"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

var (
	// errEmptySections — вызов источника без запрошенных секций. Это
	// контрактная ошибка вызывающего (не имеет смысла слать запрос с
	// пустым include — FinanceMarker отдаст все массивы пустыми), наверх
	// едет как непомеченный internal сбой.
	errEmptySections = errors.New("financemarker stock: no sections requested")

	// canonicalIncludeOrder — фиксированный порядок блоков в строке include.
	// Совпадает с порядком полей company.Stock и порядком разделов в
	// справочнике api-financemarker/references/stock.md. Один use case → одна
	// строка include — это требование справочника (Edge cases stock.md) и условие
	// устойчивого ключа серверного и локального HTTP-кеша.
	canonicalIncludeOrder = []company.StockSection{
		company.StockSectionInfo,
		company.StockSectionSummary,
	}

	// includeCodeBySection переводит секцию в код блока FinanceMarker.
	// StockSectionUnspecified присутствует с пустым кодом и отбрасывается
	// в buildInclude (см. lookup на ok && code != "").
	includeCodeBySection = map[company.StockSection]string{
		company.StockSectionUnspecified: "",
		company.StockSectionInfo:        "info",
		company.StockSectionSummary:     "summary",
	}
)

// buildInclude собирает значение query-параметра include в каноническом
// порядке. Дубликаты и StockSectionUnspecified отбрасываются. Пустой
// результат означает, что вызывающий не запросил ни одной секции —
// функция возвращает errEmptySections.
func buildInclude(sections []company.StockSection) (string, error) {
	requested := make(map[company.StockSection]struct{}, len(sections))
	for _, s := range sections {
		code, ok := includeCodeBySection[s]
		if !ok || code == "" {
			continue
		}
		requested[s] = struct{}{}
	}
	if len(requested) == 0 {
		return "", errEmptySections
	}

	codes := make([]string, 0, len(requested))
	for _, s := range canonicalIncludeOrder {
		if _, ok := requested[s]; ok {
			codes = append(codes, includeCodeBySection[s])
		}
	}
	return strings.Join(codes, ","), nil
}
