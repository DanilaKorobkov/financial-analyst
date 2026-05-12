// Package financemarker — реализации domain-портов поверх REST API
// FinanceMarker.ru (https://financemarker.ru/api/swagger-ui/).
package financemarker

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// NewClient собирает resty-клиент для FinanceMarker.
//
// baseURL — корень API без завершающего слэша (например,
// "https://financemarker.ru/api/fm/v2"). token прокидывается query-параметром
// "api_token" во все запросы. timeout применяется к каждому HTTP-вызову.
func NewClient(baseURL, token string, timeout time.Duration) *resty.Client {
	return resty.New().
		SetBaseURL(baseURL).
		SetTimeout(timeout).
		SetQueryParam("api_token", token)
}
