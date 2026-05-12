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

// includeInfo — значение query-параметра include, ограничивающее ответ
// блоком info (классификация, описание, ссылки).
const includeInfo = "info"

// CompanyCardRepository — реализация entities.CompanyCardRepository поверх
// FinanceMarker /api/fm/v2/stocks/{exchange}:{code}.
type CompanyCardRepository struct {
	client *resty.Client
}

// NewCompanyCardRepository собирает репозиторий вокруг FinanceMarker-клиента.
func NewCompanyCardRepository(client *resty.Client) *CompanyCardRepository {
	return &CompanyCardRepository{client: client}
}

// FindByTicker запрашивает карточку эмитента и маппит её в entities.CompanyCard.
// Перевод HTTP-ошибок в domain-ошибки — в mapHTTPError.
func (r *CompanyCardRepository) FindByTicker(
	ctx context.Context,
	ticker string,
) (entities.CompanyCard, error) {
	resp, err := r.client.R().
		SetContext(ctx).
		SetPathParam("symbol", exchangeMOEX+":"+ticker).
		SetQueryParam("include", includeInfo).
		Get("/stocks/{symbol}")
	if err != nil {
		return entities.CompanyCard{}, fmt.Errorf("financemarker request: %w", err)
	}
	if err := mapHTTPError(resp); err != nil {
		return entities.CompanyCard{}, err
	}

	var dto stockDTO
	if err := jsonParser.Unmarshal(resp.Body(), &dto); err != nil {
		return entities.CompanyCard{}, fmt.Errorf("decode financemarker payload: %w", err)
	}

	return mapCompanyCard(&dto.Info), nil
}
