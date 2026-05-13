// Package financemarker — общий resty-клиент FinanceMarker REST API и
// провайдер-специфичные детали (классификация HTTP-ошибок, error body).
// Подпакеты реализуют конкретные domain-порты поверх этого клиента.
package financemarker

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

// Infra-уровневые ошибки FinanceMarker. Подпакеты при необходимости
// переводят их в доменные sentinel-ы своего агрегата; всё остальное
// едет наверх как непомеченный «внутренний сбой» (presentation → CodeInternal).
var (
	// ErrNotFound — FinanceMarker вернул HTTP 404.
	ErrNotFound = errors.New("financemarker: not found")

	// ErrUnauthorized — токен не принят (401 или 400 token_not_found).
	ErrUnauthorized = errors.New("financemarker: unauthorized")

	// ErrQuotaExceeded — превышен лимит запросов (403).
	ErrQuotaExceeded = errors.New("financemarker: quota exceeded")
)

// ConfigClient — параметры доступа к FinanceMarker REST API.
type ConfigClient struct {
	// BaseURL — корень FinanceMarker REST API без завершающего слэша,
	// например "https://financemarker.ru/api/fm/v2".
	BaseURL string

	// Token — API-токен из профиля FinanceMarker. Передаётся query-параметром
	// "api_token" во всех запросах.
	Token string

	// Timeout — таймаут на один HTTP-запрос.
	Timeout time.Duration
}

// Client — тонкая обёртка над resty.Client, преднастроенная под FinanceMarker:
// jsoniter-парсер, api_token, тип errorBody и middleware, превращающий
// ошибочные HTTP-ответы в типизированные ошибки.
//
// Встраивается в конкретные gateway подпакетов (company, ...).
type Client struct {
	*resty.Client
}

// errorBody — JSON-обёртка ошибочного ответа FinanceMarker. Регистрируется
// клиенту через SetError — resty сам разбирает её на ответах с не-2xx
// статусом и кладёт в resp.Error().
type errorBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewClient собирает Client под FinanceMarker.
func NewClient(cfg ConfigClient) *Client {
	jsonParser := jsoniter.ConfigCompatibleWithStandardLibrary

	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetTimeout(cfg.Timeout).
		SetQueryParam("api_token", cfg.Token).
		SetJSONUnmarshaler(jsonParser.Unmarshal).
		SetJSONMarshaler(jsonParser.Marshal).
		SetError(&errorBody{}).
		OnAfterResponse(func(_ *resty.Client, resp *resty.Response) error {
			return classifyError(resp)
		})
	return &Client{Client: client}
}

// classifyError переводит ошибочный HTTP-ответ FinanceMarker в ошибку слоя
// infra. Вызывается middleware-ом OnAfterResponse, поэтому возвращаемая
// ошибка приходит наверх как err из R().Get(...).
//
// Возвращает infra-sentinel-ы (ErrNotFound / ErrUnauthorized /
// ErrQuotaExceeded) с обёрнутой исходной причиной — подпакет уже переводит
// их в доменные ошибки своего агрегата.
func classifyError(resp *resty.Response) error {
	if !resp.IsError() {
		return nil
	}

	status := resp.StatusCode()

	var message string
	if body, ok := resp.Error().(*errorBody); ok && body != nil {
		message = body.Message
	}

	switch {
	case status == http.StatusNotFound:
		return ErrNotFound
	case status == http.StatusBadRequest && message == "token_not_found":
		return fmt.Errorf("%w: token_not_found", ErrUnauthorized)
	case status == http.StatusUnauthorized:
		return fmt.Errorf("%w: http status %d", ErrUnauthorized, status)
	case status == http.StatusForbidden:
		return fmt.Errorf("%w: http status %d", ErrQuotaExceeded, status)
	default:
		return fmt.Errorf("financemarker http status %d", status)
	}
}
