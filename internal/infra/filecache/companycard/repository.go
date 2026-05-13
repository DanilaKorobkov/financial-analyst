// Package companycard — реализация companycard.Repository по шаблону
// Proxy: кеширует карточки эмитентов на локальной файловой системе
// (поверх diskv) и делегирует «холодные» запросы нижележащему репозиторию.
package companycard

import (
	"context"
	"fmt"
	"net/url"

	jsoniter "github.com/json-iterator/go"
	"github.com/peterbourgon/diskv/v3"

	domaincard "github.com/DanilaKorobkov/financial-analyst/internal/domain/companycard"
)

// cacheMaxSize — лимит in-memory кеша diskv. 0 — отключён: храним всё
// исключительно на диске, чтобы поведение Proxy было предсказуемо
// и одинаково между процессами.
const cacheMaxSize = 0

// jsonParser — drop-in для encoding/json (см. internal/infra/financemarker).
var jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary

// ConfigRepository — параметры файлового кеша карточек.
type ConfigRepository struct {
	// Delegate — нижележащий репозиторий, к которому идёт запрос при cache miss.
	Delegate domaincard.Repository

	// Dir — каталог хранения файлов кеша. Создаётся diskv при первой записи.
	Dir string
}

// Repository — Proxy над domaincard.Repository: на cache hit возвращает
// карточку, прочитанную с диска; на cache miss идёт в Delegate и
// сохраняет результат файлом.
type Repository struct {
	delegate domaincard.Repository
	store    *diskv.Diskv
}

// NewRepository собирает файловый кеш-репозиторий поверх diskv.
func NewRepository(cfg ConfigRepository) *Repository {
	store := diskv.New(diskv.Options{
		BasePath:     cfg.Dir,
		CacheSizeMax: cacheMaxSize,
	})
	return &Repository{
		delegate: cfg.Delegate,
		store:    store,
	}
}

// FindByTicker сначала пытается отдать карточку из файла кеша. При
// промахе обращается к Delegate и, если запрос успешен, сохраняет
// карточку на диск. Ошибки Delegate (включая domaincard.ErrNotFound) на
// диск не пишутся.
func (r *Repository) FindByTicker(
	ctx context.Context,
	exchange domaincard.Exchange,
	ticker string,
) (domaincard.Card, error) {
	key := cacheKey(exchange, ticker)
	if card, hit := r.readCache(key); hit {
		return card, nil
	}

	card, err := r.delegate.FindByTicker(ctx, exchange, ticker)
	if err != nil {
		return domaincard.Card{}, err //nolint:wrapcheck // ошибка делегата идёт наверх как есть
	}

	if writeErr := r.writeCache(key, &card); writeErr != nil {
		return domaincard.Card{}, writeErr
	}
	return card, nil
}

// cacheKey собирает ключ diskv для пары (exchange, ticker). Биржа
// кодируется числовым значением enum-а, чтобы переименование констант
// не разменивало уже сохранённые файлы. Тикер пропускается через
// url.PathEscape: на выходе всегда валидный для diskv ключ без
// разделителей путей и NUL-байтов — сохранить в кеш можно любой тикер.
func cacheKey(exchange domaincard.Exchange, ticker string) string {
	return fmt.Sprintf("%d_%s.json", int(exchange), url.PathEscape(ticker))
}

// readCache читает карточку из diskv. Отсутствие файла, битый JSON или
// ошибка ввода-вывода — всё трактуется как cache miss: лучше сходить
// в Delegate, чем отдать сломанную запись.
func (r *Repository) readCache(key string) (domaincard.Card, bool) {
	raw, err := r.store.Read(key)
	if err != nil {
		return domaincard.Card{}, false
	}
	var card domaincard.Card
	if err := jsonParser.Unmarshal(raw, &card); err != nil {
		return domaincard.Card{}, false
	}
	return card, true
}

// writeCache кодирует карточку в JSON и кладёт её в diskv. Сам diskv
// пишет через tmp-файл + rename, поэтому параллельные писатели не дают
// читателю «полуписанного» файла.
func (r *Repository) writeCache(key string, card *domaincard.Card) error {
	raw, err := jsonParser.Marshal(card)
	if err != nil {
		return fmt.Errorf("filecache companycard: marshal card: %w", err)
	}
	if err := r.store.Write(key, raw); err != nil {
		return fmt.Errorf("filecache companycard: write: %w", err)
	}
	return nil
}
