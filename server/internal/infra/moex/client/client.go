// Package client — общий resty-клиент MOEX ISS REST API и
// провайдер-специфичная классификация HTTP-ошибок.
package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

// Config — параметры доступа к MOEX ISS REST API.
type Config struct {
	// BaseURL — корень MOEX ISS без завершающего слэша,
	// например "https://iss.moex.com/iss".
	BaseURL string

	// Timeout — таймаут на один HTTP-запрос.
	Timeout time.Duration
}

// Client — тонкая обёртка над resty.Client, преднастроенная под MOEX ISS:
// jsoniter-парсер, query-параметры iss.* и middleware, переводящий не-2xx
// статусы в типизированную ошибку.
type Client struct {
	*resty.Client
}

// New собирает Client под MOEX ISS.
//
// Параметры iss.* — это управляющие query-параметры самого MOEX ISS,
// определены в их публичной справке (https://iss.moex.com/iss/reference/),
// общие для всех эндпоинтов:
//   - iss.json=extended    — JSON в виде массива записей вместо
//     «columns + data»; удобнее разбирать (один объект на запись).
//   - iss.meta=off         — не присылать блок описания колонок (типы,
//     длины); экономит трафик и упрощает разбор.
//   - iss.only=description — из всех блоков ответа (securities, boards,
//     marketdata и т.п.) вернуть только блок description с реквизитами
//     эмитента — это всё, что нужно справочным репозиториям.
//
// Если bundle-у потребуется другой набор iss.only — он переопределит
// query-параметр на уровне своего R().SetQueryParam(...).
func New(cfg Config) *Client {
	jsonParser := jsoniter.ConfigCompatibleWithStandardLibrary

	c := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetTimeout(cfg.Timeout).
		SetQueryParams(map[string]string{
			"iss.json": "extended",
			"iss.meta": "off",
			"iss.only": "description",
		}).
		SetJSONUnmarshaler(jsonParser.Unmarshal).
		SetJSONMarshaler(jsonParser.Marshal).
		OnAfterResponse(func(_ *resty.Client, resp *resty.Response) error {
			if resp.IsError() {
				return fmt.Errorf("moex http status %d", resp.StatusCode())
			}
			return nil
		})
	return &Client{Client: c}
}
