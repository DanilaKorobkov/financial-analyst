// Package financemarker — реализации domain-портов поверх REST API
// FinanceMarker.ru (https://financemarker.ru/api/swagger-ui/).
package financemarker

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// includeInfo — значение query-параметра include, ограничивающее ответ
// блоком info (классификация, описание, ссылки).
const includeInfo = "info"

// codeExchangeMOEX — строковый код Московской биржи в формате FinanceMarker.
const codeExchangeMOEX = "MOEX"

// ConfigCompanyCardRepository — параметры доступа к FinanceMarker для
// репозитория карточек эмитента.
type ConfigCompanyCardRepository struct {
	// BaseURL — корень FinanceMarker REST API без завершающего слэша,
	// например "https://financemarker.ru/api/fm/v2".
	BaseURL string

	// Token — API-токен из профиля FinanceMarker. Передаётся query-параметром
	// "api_token" во всех запросах.
	Token string

	// Timeout — таймаут на один HTTP-запрос.
	Timeout time.Duration
}

// CompanyCardRepository — реализация entities.CompanyCardRepository поверх
// FinanceMarker /api/fm/v2/stocks/{exchange}:{code}.
type CompanyCardRepository struct {
	client *resty.Client
}

// NewCompanyCardRepository собирает репозиторий вокруг своего resty-клиента.
func NewCompanyCardRepository(cfg ConfigCompanyCardRepository) *CompanyCardRepository {
	return &CompanyCardRepository{client: newRestyClient(cfg.BaseURL, cfg.Token, cfg.Timeout)}
}

// newRestyClient собирает resty-клиент для FinanceMarker: фиксирует базовый
// URL, таймаут и автоматически прокидывает api_token во все запросы.
func newRestyClient(baseURL, token string, timeout time.Duration) *resty.Client {
	return resty.New().
		SetBaseURL(baseURL).
		SetTimeout(timeout).
		SetQueryParam("api_token", token)
}

// FindByTicker запрашивает карточку эмитента и переводит её в entities.CompanyCard.
// Биржу принимаем явно: символ для FM собирается как "{exchange}:{ticker}".
// Перевод HTTP-ошибок в domain/internal-ошибки — в mapHTTPError.
func (r *CompanyCardRepository) FindByTicker(
	ctx context.Context,
	exchange entities.Exchange,
	ticker string,
) (entities.CompanyCard, error) {
	symbol, err := buildSymbol(exchange, ticker)
	if err != nil {
		return entities.CompanyCard{}, err
	}

	resp, err := r.client.R().
		SetContext(ctx).
		SetPathParam("symbol", symbol).
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

	return translateCompanyCard(&dto.Info), nil
}

// buildSymbol собирает path-параметр FinanceMarker вида "{exchange}:{ticker}".
// Неподдерживаемая биржа — ошибка: запрос с непустым кодом, который FM не
// признаёт, всё равно вернёт 404 / 400, и лучше отвалиться явно.
func buildSymbol(exchange entities.Exchange, ticker string) (string, error) {
	code, err := exchangeCode(exchange)
	if err != nil {
		return "", err
	}
	return code + ":" + ticker, nil
}

// exchangeCode возвращает строковый код биржи в формате FinanceMarker.
func exchangeCode(exchange entities.Exchange) (string, error) {
	switch exchange {
	case entities.ExchangeMOEX:
		return codeExchangeMOEX, nil
	case entities.ExchangeUnspecified:
		return "", fmt.Errorf("financemarker: exchange is unspecified")
	default:
		return "", fmt.Errorf("financemarker: unsupported exchange %d", exchange)
	}
}
