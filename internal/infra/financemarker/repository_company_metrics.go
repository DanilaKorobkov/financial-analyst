package financemarker

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// exchangeMOEX — единственная биржа, поддерживаемая текущей подпиской
// FinanceMarker. Часть path-параметра запроса: /stocks/{exchange}:{ticker}.
const exchangeMOEX = "MOEX"

// includeInfoSummary — значение query-параметра include, ограничивающее
// ответ блоками info и summary (расширенная карточка эмитента).
const includeInfoSummary = "info,summary"

// CompanyMetricsRepository — реализация entities.CompanyMetricsRepository
// поверх FinanceMarker /api/fm/v2/stocks/{exchange}:{code}.
type CompanyMetricsRepository struct {
	client *resty.Client
}

// NewCompanyMetricsRepository собирает репозиторий вокруг FinanceMarker-клиента.
func NewCompanyMetricsRepository(client *resty.Client) *CompanyMetricsRepository {
	return &CompanyMetricsRepository{client: client}
}

// FindByTicker запрашивает расширенную карточку эмитента и маппит её в
// entities.CompanyMetrics. Перевод HTTP-ошибок в domain-ошибки — в mapHTTPError.
func (r *CompanyMetricsRepository) FindByTicker(
	ctx context.Context,
	ticker string,
) (entities.CompanyMetrics, error) {
	resp, err := r.client.R().
		SetContext(ctx).
		SetPathParam("symbol", exchangeMOEX+":"+ticker).
		SetQueryParam("include", includeInfoSummary).
		Get("/stocks/{symbol}")
	if err != nil {
		return entities.CompanyMetrics{}, fmt.Errorf("financemarker request: %w", err)
	}
	if err := mapHTTPError(resp); err != nil {
		return entities.CompanyMetrics{}, err
	}

	var dto stockDTO
	if err := jsonParser.Unmarshal(resp.Body(), &dto); err != nil {
		return entities.CompanyMetrics{}, fmt.Errorf("decode financemarker payload: %w", err)
	}

	return mapCompanyMetrics(&dto), nil
}
