// Package company — реализация company.ClassificationGateway по шаблону
// Proxy: кеширует классификационные секции компаний на локальной
// файловой системе (поверх diskv) и делегирует «холодные» запросы
// нижележащему gateway. Кешируется только FinanceMarker: у него квота
// на запросы. Свободный MOEX-источник идёт без кеша.
package company

import (
	"context"
	"fmt"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/peterbourgon/diskv/v3"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
)

// cacheMaxSize — лимит in-memory кеша diskv. 0 — отключён: храним всё
// исключительно на диске, чтобы поведение Proxy было предсказуемо
// и одинаково между процессами.
const cacheMaxSize = 0

// jsonParser — drop-in для encoding/json.
var jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary

// ConfigClassificationProxy — параметры файлового кеша классификационной
// секции карточки.
type ConfigClassificationProxy struct {
	// Delegate — нижележащий gateway, к которому идёт запрос при cache miss.
	Delegate domaincompany.ClassificationGateway

	// Dir — каталог хранения файлов кеша. Создаётся diskv при первой записи.
	Dir string

	// TTL — срок жизни записи в кеше. Ноль — без экспирации (запись
	// действительна, пока файл лежит на диске).
	TTL time.Duration
}

// ClassificationProxy — Proxy над company.ClassificationGateway: на
// cache hit возвращает классификацию, прочитанную с диска; на cache miss
// идёт в Delegate и сохраняет результат файлом.
type ClassificationProxy struct {
	delegate domaincompany.ClassificationGateway
	store    *diskv.Diskv
	ttl      time.Duration
}

// NewClassificationProxy собирает файловый кеш поверх diskv.
func NewClassificationProxy(cfg ConfigClassificationProxy) *ClassificationProxy {
	store := diskv.New(diskv.Options{
		BasePath:     cfg.Dir,
		CacheSizeMax: cacheMaxSize,
	})
	return &ClassificationProxy{
		delegate: cfg.Delegate,
		store:    store,
		ttl:      cfg.TTL,
	}
}

// FindByTicker сначала пытается отдать классификацию из файла кеша. При
// промахе или истёкшем TTL обращается к Delegate и, если запрос
// успешен, перезаписывает файл свежей записью. Ошибки Delegate
// (включая domaincompany.ErrNotFound) на диск не пишутся.
func (p *ClassificationProxy) FindByTicker(
	ctx context.Context,
	ticker string,
) (domaincompany.Classification, error) {
	key := cacheKey(ticker)
	if cls, hit := p.readCache(key); hit {
		return cls, nil
	}

	cls, err := p.delegate.FindByTicker(ctx, ticker)
	if err != nil {
		return domaincompany.Classification{}, err //nolint:wrapcheck // ошибка делегата идёт наверх как есть
	}

	if writeErr := p.writeCache(key, &cls); writeErr != nil {
		return domaincompany.Classification{}, writeErr
	}
	return cls, nil
}

// cacheKey собирает ключ diskv для тикера. Тикер пропускается через
// url.PathEscape: на выходе всегда валидный для diskv ключ без
// разделителей путей и NUL-байтов — сохранить в кеш можно любой тикер.
func cacheKey(ticker string) string {
	return url.PathEscape(ticker) + ".json"
}

// readCache читает конверт из diskv и проверяет срок жизни. Отсутствие
// файла, битый JSON, ошибка ввода-вывода или истёкший ExpiresAt — всё
// трактуется как cache miss: лучше сходить в Delegate, чем отдать
// сломанную или просроченную запись.
func (p *ClassificationProxy) readCache(key string) (domaincompany.Classification, bool) {
	raw, err := p.store.Read(key)
	if err != nil {
		return domaincompany.Classification{}, false
	}
	var envelope classificationEnvelope
	if err := jsonParser.Unmarshal(raw, &envelope); err != nil {
		return domaincompany.Classification{}, false
	}
	if !envelope.ExpiresAt.IsZero() && !time.Now().UTC().Before(envelope.ExpiresAt) {
		return domaincompany.Classification{}, false
	}
	return classificationFromDTO(&envelope.Classification), true
}

// writeCache упаковывает классификацию в конверт и кладёт её в diskv.
// Сам diskv пишет через tmp-файл + rename, поэтому параллельные
// писатели не дают читателю «полуписанного» файла. При TTL == 0
// ExpiresAt остаётся нулевым — такая запись не протухает.
func (p *ClassificationProxy) writeCache(key string, cls *domaincompany.Classification) error {
	envelope := classificationEnvelope{Classification: classificationToDTO(cls)}
	if p.ttl > 0 {
		envelope.ExpiresAt = time.Now().UTC().Add(p.ttl)
	}
	raw, err := jsonParser.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("filecache company: marshal envelope: %w", err)
	}
	if err := p.store.Write(key, raw); err != nil {
		return fmt.Errorf("filecache company: write: %w", err)
	}
	return nil
}
