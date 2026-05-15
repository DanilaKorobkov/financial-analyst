// Package filecache — generic файловый кеш с per-call TTL поверх diskv.
//
// Кеш работает над T, поддерживающим JSON-раунд-трип через jsoniter:
// значение имеет конкретный статический тип, и round-trip через JSON
// ничего не теряет. Слой ничего не знает про domain-типы и FieldDescriptor.
//
// API из двух уровней. Get/Put — низкоуровневые примитивы для случаев,
// где нужна условная запись (например, http-кеш записывает только 2xx).
// LoadOrFetch — удобная обёртка над типичным шаблоном «прочитать, на
// промахе сходить к источнику и записать».
package filecache

import (
	"context"
	"fmt"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/peterbourgon/diskv/v3"
)

// cacheMaxSize — лимит in-memory кеша diskv. 0 — отключён: храним всё
// исключительно на диске, чтобы поведение кеша было предсказуемо
// и одинаково между процессами.
const cacheMaxSize = 0

// jsonParser — drop-in для encoding/json.
var jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary

// Store — типизированный файловый кеш над значениями типа T. T обязан
// корректно проходить jsoniter Marshal/Unmarshal. TTL задаётся per-call:
// один и тот же Store обслуживает записи с разным сроком жизни.
type Store[T any] struct {
	store *diskv.Diskv
}

// Config — параметры конструктора Store.
type Config struct {
	// Dir — каталог хранения файлов кеша. Создаётся diskv при первой записи.
	Dir string
}

// envelope — формат записи на диске. ExpiresAt в UTC RFC3339, чтобы файл
// был самодостаточным и читался глазами. Нулевое ExpiresAt — без
// экспирации. Value — само значение T, упакованное jsoniter'ом.
type envelope[T any] struct {
	ExpiresAt time.Time `json:"expires_at"`
	Value     T         `json:"value"`
}

// New собирает Store поверх diskv. Каталог создаётся лениво при первой
// записи; никаких побочных эффектов на старте.
func New[T any](cfg Config) *Store[T] {
	return &Store[T]{
		store: diskv.New(diskv.Options{
			BasePath:     cfg.Dir,
			CacheSizeMax: cacheMaxSize,
		}),
	}
}

// Get возвращает значение и true, если файл существует, свежий и
// корректно декодируется. В остальных случаях — zero/false без ошибки:
// для вызывающего это cache miss. Любая ошибка ввода-вывода, битый JSON
// или истёкший срок жизни трактуются одинаково — как промах.
func (s *Store[T]) Get(key string) (T, bool) {
	var zero T
	raw, err := s.store.Read(diskKey(key))
	if err != nil {
		return zero, false
	}
	var env envelope[T]
	if err := jsonParser.Unmarshal(raw, &env); err != nil {
		return zero, false
	}
	if !env.ExpiresAt.IsZero() && !time.Now().UTC().Before(env.ExpiresAt) {
		return zero, false
	}
	return env.Value, true
}

// Put упаковывает значение в envelope и кладёт в diskv. TTL == 0 —
// запись без экспирации (живёт пока лежит файл). Сам diskv пишет через
// tmp-файл + rename, поэтому параллельные писатели не дают читателю
// «полуписанного» файла.
func (s *Store[T]) Put(key string, v T, ttl time.Duration) error {
	env := envelope[T]{Value: v}
	if ttl > 0 {
		env.ExpiresAt = time.Now().UTC().Add(ttl)
	}
	raw, err := jsonParser.Marshal(env)
	if err != nil {
		return fmt.Errorf("filecache: marshal envelope: %w", err)
	}
	if err := s.store.Write(diskKey(key), raw); err != nil {
		return fmt.Errorf("filecache: write: %w", err)
	}
	return nil
}

// LoadOrFetch — удобная обёртка над Get/Put для типичного шаблона:
// читаем, на промахе зовём fetch, на успехе записываем. TTL передаётся
// в Put. Ошибки fetch на диск не пишутся и идут наверх как есть.
func (s *Store[T]) LoadOrFetch(ctx context.Context, key string, ttl time.Duration, fetch func(ctx context.Context) (T, error)) (T, error) {
	if v, ok := s.Get(key); ok {
		return v, nil
	}

	v, err := fetch(ctx)
	if err != nil {
		var zero T
		return zero, err
	}

	if writeErr := s.Put(key, v, ttl); writeErr != nil {
		return v, writeErr
	}
	return v, nil
}

// diskKey приводит произвольный ключ к безопасной для diskv форме:
// PathEscape убирает разделители путей и NUL-байты, расширение .json
// делает файл понятным глазами.
func diskKey(key string) string {
	return url.PathEscape(key) + ".json"
}
