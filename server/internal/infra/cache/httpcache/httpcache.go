// Package httpcache — кеширующий http.RoundTripper, который решает,
// кешировать ли запрос, по TTL из context.Context самого запроса.
//
// TTL декларируется на стороне вызывающего (bundle) через WithTTL и
// едет вместе с запросом по контексту. Если TTL не задан или равен
// нулю, transport прозрачно отдаёт запрос наверх в base RoundTripper —
// никакого кеширования. Если задан положительный TTL, кеш проверяет
// файл по ключу запроса; на промахе идёт в сеть и записывает 2xx-ответ
// с этим TTL.
//
// Слой ничего не знает про конкретного провайдера и про domain — это
// general-purpose http-кеш, который умеет ровно одно: «по TTL из ctx
// материализовать ответ из файла или из сети».
package httpcache

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/DanilaKorobkov/financial-analyst/internal/infra/cache/filecache"
)

// ttlKey — тип-ключ для context.Value. Уникальный пустой struct — это
// канонический способ избежать коллизии ключей с чужими пакетами.
type ttlKey struct{}

// Entry — то, что лежит в файле кеша. Хранится Header и Body, чтобы
// синтетический http.Response после cache hit вёл себя как обычный
// ответ для resty (Content-Type, Content-Length и т. п.). Экспортирован,
// чтобы вызывающий мог собрать filecache.Store[Entry] и отдать его
// в NewTransport.
type Entry struct {
	// Header — заголовки исходного ответа, без модификаций.
	Header http.Header `json:"header"`

	// Body — полное тело исходного ответа.
	Body []byte `json:"body"`

	// StatusCode — HTTP-код исходного ответа (только 2xx попадают в кеш).
	StatusCode int `json:"status_code"`
}

// Transport — кеширующий http.RoundTripper. Заворачивает любой base
// transport и материализует cache hit без обращения в base.
type Transport struct {
	base  http.RoundTripper
	store *filecache.Store[Entry]
}

// WithTTL аннотирует ctx сроком жизни кеш-записи для исходящего HTTP-запроса.
// Bundle вызывает её прямо перед resty-запросом, который должен кешироваться.
// Положительный TTL — кешировать на этот срок; ноль или отрицательный —
// эквивалентно «не аннотировать»: transport прозрачно идёт в сеть.
func WithTTL(ctx context.Context, ttl time.Duration) context.Context {
	if ttl <= 0 {
		return ctx
	}
	return context.WithValue(ctx, ttlKey{}, ttl)
}

// ttlFromContext возвращает TTL, выставленный WithTTL, и признак «было задано».
func ttlFromContext(ctx context.Context) (time.Duration, bool) {
	ttl, ok := ctx.Value(ttlKey{}).(time.Duration)
	return ttl, ok && ttl > 0
}

// NewTransport собирает Transport поверх готового store. Base — обычно
// http.DefaultTransport либо предварительно настроенный пользовательский
// transport.
func NewTransport(base http.RoundTripper, store *filecache.Store[Entry]) *Transport {
	return &Transport{base: base, store: store}
}

// RoundTrip — реализация http.RoundTripper. Логика:
//  1. Нет TTL в ctx — прозрачно в base.
//  2. Есть запись в кеше — возвращаем синтетический ответ без сети.
//  3. Иначе — base.RoundTrip, на 2xx пишем в кеш с TTL из ctx.
//
// Не-2xx и сетевые ошибки в кеш не попадают.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ttl, ok := ttlFromContext(req.Context())
	if !ok {
		return t.base.RoundTrip(req) //nolint:wrapcheck // транзит ошибки base transport
	}

	key := requestKey(req)
	if e, hit := t.store.Get(key); hit {
		return buildResponse(req, e), nil
	}

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return resp, err //nolint:wrapcheck // транзит ошибки base transport
	}

	body, readErr := io.ReadAll(resp.Body)
	if closeErr := resp.Body.Close(); closeErr != nil && readErr == nil {
		readErr = closeErr
	}
	if readErr != nil {
		return nil, fmt.Errorf("httpcache: read response body: %w", readErr)
	}
	resp.Body = io.NopCloser(bytes.NewReader(body))

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		if err := t.store.Put(key, Entry{
			StatusCode: resp.StatusCode,
			Header:     resp.Header.Clone(),
			Body:       body,
		}, ttl); err != nil {
			return nil, fmt.Errorf("httpcache: write entry: %w", err)
		}
	}

	return resp, nil
}

// requestKey собирает детерминированный ключ запроса. Метод + полный
// путь + отсортированные query-параметры за вычетом auth/корреляционных
// (см. isIgnoredQueryParam). Host в ключ не входит: тот же source за
// тем же путём с другим хостом (sandbox vs prod) — это разные конфигурации
// клиента, не разные запросы.
func requestKey(req *http.Request) string {
	query := req.URL.Query()
	keys := make([]string, 0, len(query))
	for k := range query {
		if isIgnoredQueryParam(k) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var qs strings.Builder
	for i, k := range keys {
		if i > 0 {
			qs.WriteByte('&')
		}
		values := query[k]
		sort.Strings(values)
		for j, v := range values {
			if j > 0 {
				qs.WriteByte('&')
			}
			qs.WriteString(url.QueryEscape(k))
			qs.WriteByte('=')
			qs.WriteString(url.QueryEscape(v))
		}
	}

	if qs.Len() == 0 {
		return req.Method + " " + req.URL.Path
	}
	return req.Method + " " + req.URL.Path + "?" + qs.String()
}

// isIgnoredQueryParam решает, надо ли исключить query-параметр из ключа
// кеша. Сюда попадают auth/корреляционные значения, которые не влияют
// на сам payload ответа: их вариативность не должна порождать «новые»
// записи в кеше.
func isIgnoredQueryParam(name string) bool {
	return name == "api_token"
}

// buildResponse материализует http.Response из cache entry. Все поля,
// которые resty/http.Client потенциально читает (Status, StatusCode,
// Header, Body, ContentLength, Request, Proto), заполнены, чтобы
// синтетический ответ был неотличим от настоящего.
func buildResponse(req *http.Request, e Entry) *http.Response {
	body := io.NopCloser(bytes.NewReader(e.Body))
	return &http.Response{
		Status:        fmt.Sprintf("%d %s", e.StatusCode, http.StatusText(e.StatusCode)),
		StatusCode:    e.StatusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        e.Header.Clone(),
		Body:          body,
		ContentLength: int64(len(e.Body)),
		Request:       req,
	}
}
