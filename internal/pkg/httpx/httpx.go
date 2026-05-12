// Package httpx — общие фабрики HTTP-клиентов проекта. Здесь живут
// дефолты, которые одинаковы для всех адаптеров внешних REST-источников:
// JSON-парсер на jsoniter и middleware-обвязка для классификации ошибок.
//
// Провайдер берёт httpx.New(...) как болванку и далее настраивает своё —
// query-параметры, тип ошибочного payload-а, конкретный классификатор.
package httpx

import (
	"time"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

// Config — параметры базового resty-клиента.
type Config struct {
	// BaseURL — корень REST API без завершающего слэша.
	BaseURL string

	// Timeout — таймаут на один HTTP-запрос.
	Timeout time.Duration
}

// New собирает resty-клиент с дефолтами проекта: jsoniter в качестве
// JSON-парсера, заданные BaseURL и Timeout. Дальнейшую провайдер-специфику
// (query-параметры, SetError, OnError) навешивает вызывающая сторона.
func New(cfg Config) *resty.Client {
	jsonParser := jsoniter.ConfigCompatibleWithStandardLibrary

	return resty.New().
		SetBaseURL(cfg.BaseURL).
		SetTimeout(cfg.Timeout).
		SetJSONUnmarshaler(jsonParser.Unmarshal).
		SetJSONMarshaler(jsonParser.Marshal)
}

// OnError вешает middleware, превращающее ошибочные HTTP-ответы в Go-ошибки
// через переданный пользовательский классификатор. Возвращает тот же клиент,
// чтобы было удобно вызывать цепочкой после httpx.New(...).
//
// Классификатор должен вернуть nil, если ответ не ошибочный (resty всё равно
// зовёт OnAfterResponse и на 2xx).
func OnError(c *resty.Client, classify func(*resty.Response) error) *resty.Client {
	return c.OnAfterResponse(func(_ *resty.Client, resp *resty.Response) error {
		return classify(resp)
	})
}
