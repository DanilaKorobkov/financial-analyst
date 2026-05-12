// Package financemarker — реализации domain-портов поверх REST API
// FinanceMarker.ru (https://financemarker.ru/api/swagger-ui/).
package financemarker

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

const (
	// includeInfo — значение query-параметра include, ограничивающее ответ
	// блоком info (классификация, описание, ссылки).
	includeInfo = "info"

	// codeExchangeMOEX — строковый код Московской биржи в формате FinanceMarker.
	codeExchangeMOEX = "MOEX"
)

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

// newRestyClient собирает resty-клиент для FinanceMarker. Декларативная
// конфигурация: базовый URL, таймаут, api_token, jsoniter-парсер, тип
// ошибочного body (errorBody) и middleware, который превращает ошибочные
// HTTP-ответы в типизированные ошибки. После этого репозиторий пишет только
// запрос и тип ответа — статус и payload разбираются прозрачно.
func newRestyClient(baseURL, token string, timeout time.Duration) *resty.Client {
	jsonParser := jsoniter.ConfigCompatibleWithStandardLibrary

	return resty.New().
		SetBaseURL(baseURL).
		SetTimeout(timeout).
		SetQueryParam("api_token", token).
		SetJSONUnmarshaler(jsonParser.Unmarshal).
		SetJSONMarshaler(jsonParser.Marshal).
		SetError(&errorBody{}).
		OnAfterResponse(func(_ *resty.Client, resp *resty.Response) error {
			return classifyError(resp)
		})
}

// FindByTicker запрашивает карточку эмитента и переводит её в entities.CompanyCard.
// Биржу принимаем явно: символ для FM собирается как "{exchange}:{ticker}".
// Сетевые и HTTP-ошибки приходят из resty-клиента уже классифицированными
// (см. newRestyClient / classifyError), сюда — только транспорт.
func (r *CompanyCardRepository) FindByTicker(
	ctx context.Context,
	exchange entities.Exchange,
	ticker string,
) (entities.CompanyCard, error) {
	symbol, err := buildSymbol(exchange, ticker)
	if err != nil {
		return entities.CompanyCard{}, err
	}

	var dto stockDTO
	resp, err := r.client.R().
		SetContext(ctx).
		SetPathParam("symbol", symbol).
		SetQueryParam("include", includeInfo).
		SetResult(&dto).
		Get("/stocks/{symbol}")
	if err != nil {
		switch {
		case resp == nil || resp.StatusCode() == 0:
			return entities.CompanyCard{}, fmt.Errorf("financemarker request: %w", err)
		case !resp.IsError():
			return entities.CompanyCard{}, fmt.Errorf("decode financemarker payload: %w", err)
		default:
			return entities.CompanyCard{}, err //nolint:wrapcheck // err уже сформирован нашим OnAfterResponse (classifyError)
		}
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
