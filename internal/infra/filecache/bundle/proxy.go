// Package bundle — кеширующий Proxy над data.Bundle: на cache hit
// возвращает FieldValues, прочитанные с диска; на cache miss идёт в
// нижележащий bundle и сохраняет результат файлом. Шаблон Proxy:
// внешний контракт совпадает с data.Bundle, decorator вокруг любого
// конкретного bundle.
//
// Используется для bundles с дорогими квотами (FinanceMarker). Свободные
// источники (MOEX ISS) подключаются напрямую, без Proxy.
package bundle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/peterbourgon/diskv/v3"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
)

// cacheMaxSize — лимит in-memory кеша diskv. 0 — отключён: храним всё
// исключительно на диске, чтобы поведение Proxy было предсказуемо
// и одинаково между процессами.
const cacheMaxSize = 0

// jsonParser — drop-in для encoding/json.
var jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary

// ConfigProxy — параметры файлового кеша bundle.
type ConfigProxy struct {
	// Delegate — нижележащий bundle, к которому идёт запрос при cache miss.
	Delegate data.Bundle

	// Dir — каталог хранения файлов кеша. Создаётся diskv при первой записи.
	Dir string

	// TTL — срок жизни записи в кеше. Ноль — без экспирации (запись
	// действительна, пока файл лежит на диске).
	TTL time.Duration
}

// Proxy — кеширующий decorator над data.Bundle. Хранит FieldValues
// в файле под ключом тикера; codec ниже знает каждое поле bundle и
// упаковывает его значение по типу из FieldDescriptor — после
// чтения с диска `int64` остаётся `int64`, а не превращается в `float64`,
// как было бы при обычном `json.Unmarshal` в `map[string]any`.
type Proxy struct {
	delegate data.Bundle
	store    *diskv.Diskv
	codec    *valueCodec
	ttl      time.Duration
}

// envelopeOnDisk — формат записи на диске. ExpiresAt хранится в UTC
// и в RFC3339, чтобы файл оставался самодостаточным и читаемым глазами.
// Нулевое ExpiresAt означает «без экспирации». Values — карта от полного
// id поля к raw-JSON; типизация восстанавливается codec'ом при чтении.
type envelopeOnDisk struct {
	ExpiresAt time.Time                  `json:"expires_at"`
	Values    map[string]json.RawMessage `json:"values"`
}

// NewProxy собирает файловый кеш поверх diskv. Codec строится один раз
// по списку полей bundle и переиспользуется на каждый запрос.
func NewProxy(cfg ConfigProxy) *Proxy {
	store := diskv.New(diskv.Options{
		BasePath:     cfg.Dir,
		CacheSizeMax: cacheMaxSize,
	})
	return &Proxy{
		delegate: cfg.Delegate,
		store:    store,
		codec:    newValueCodec(cfg.Delegate.Fields()),
		ttl:      cfg.TTL,
	}
}

// BundleID — реализация data.Bundle: транзитом из delegate.
func (p *Proxy) BundleID() string { return p.delegate.BundleID() }

// Fields — реализация data.Bundle: транзитом из delegate.
func (p *Proxy) Fields() []data.FieldDescriptor { return p.delegate.Fields() }

// Fetch сначала пытается отдать значения из файла кеша. При промахе или
// истёкшем TTL обращается к delegate и, если запрос успешен, перезаписывает
// файл свежей записью. Ошибки delegate на диск не пишутся.
func (p *Proxy) Fetch(ctx context.Context, ticker string) (data.FieldValues, error) {
	key := cacheKey(ticker)
	if values, hit := p.readCache(key); hit {
		return values, nil
	}

	values, err := p.delegate.Fetch(ctx, ticker)
	if err != nil {
		return nil, err //nolint:wrapcheck // ошибка делегата идёт наверх как есть
	}

	if writeErr := p.writeCache(key, values); writeErr != nil {
		return nil, writeErr
	}
	return values, nil
}

// cacheKey собирает ключ diskv для тикера. Тикер пропускается через
// url.PathEscape: на выходе всегда валидный для diskv ключ без
// разделителей путей и NUL-байтов — сохранить в кеш можно любой тикер.
func cacheKey(ticker string) string {
	return url.PathEscape(ticker) + ".json"
}

// readCache читает конверт из diskv, проверяет срок жизни и декодирует
// значения по каталогу типов. Отсутствие файла, битый JSON, ошибка
// ввода-вывода, истёкший ExpiresAt или ошибка декодирования — всё
// трактуется как cache miss: лучше сходить в delegate, чем отдать
// сломанную или просроченную запись.
func (p *Proxy) readCache(key string) (data.FieldValues, bool) {
	raw, err := p.store.Read(key)
	if err != nil {
		return nil, false
	}
	var envelope envelopeOnDisk
	if unmarshalErr := jsonParser.Unmarshal(raw, &envelope); unmarshalErr != nil {
		return nil, false
	}
	if !envelope.ExpiresAt.IsZero() && !time.Now().UTC().Before(envelope.ExpiresAt) {
		return nil, false
	}
	values, err := p.codec.Decode(envelope.Values)
	if err != nil {
		return nil, false
	}
	return values, true
}

// writeCache упаковывает FieldValues в конверт через codec и кладёт его
// в diskv. Сам diskv пишет через tmp-файл + rename, поэтому параллельные
// писатели не дают читателю «полуписанного» файла. При TTL == 0 ExpiresAt
// остаётся нулевым — такая запись не протухает.
func (p *Proxy) writeCache(key string, values data.FieldValues) error {
	encoded, err := p.codec.Encode(values)
	if err != nil {
		return fmt.Errorf("filecache bundle: encode values: %w", err)
	}
	envelope := envelopeOnDisk{Values: encoded}
	if p.ttl > 0 {
		envelope.ExpiresAt = time.Now().UTC().Add(p.ttl)
	}
	raw, err := jsonParser.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("filecache bundle: marshal envelope: %w", err)
	}
	if err := p.store.Write(key, raw); err != nil {
		return fmt.Errorf("filecache bundle: write: %w", err)
	}
	return nil
}
