package data

import (
	"context"
	"fmt"

	"github.com/sourcegraph/conc/pool"
)

// Registry — реестр зарегистрированных bundles. Заполняется в
// composition root на старте и далее не меняется; чтение по lookup
// потокобезопасно при условии, что регистрация завершилась.
//
// Bundle сам знает только свой BundleID; привязка к провайдеру держится
// здесь, в реестре, и задаётся через RegisterProvider — реестр сам
// итерирует Provider.Bundles() и записывает каждый. Сам Provider и его
// bundles про реестр ничего не знают.
type Registry struct {
	bundles map[bundleKey]Registered
	byField map[Field]Registered
	fields  map[Field]FieldDescriptor
}

type bundleKey struct {
	Provider string
	Bundle   string
}

// Registered — bundle вместе с ID провайдера, под которым он
// зарегистрирован в реестре. Используется для обхода реестра наружу
// (см. *Registry.Bundles): сам Bundle про своего провайдера ничего
// не знает.
type Registered struct {
	// Bundle — зарегистрированная реализация.
	Bundle Bundle

	// ProviderID — ID провайдера, под которым bundle зарегистрирован.
	ProviderID string
}

// Provider — точка сборки адаптера внешнего источника (MOEX, FinanceMarker
// и др.). Знает свой стабильный ID и какие bundles ему принадлежат —
// уже со всеми внутренними обвязками (например, файловым кешем).
// Конструктор конкретного Provider принимает всё, что ему нужно (клиент,
// кеш и т. п.); composition root не знает деталей — только итерирует
// список Provider и отдаёт его реестру.
type Provider interface {
	// ID — стабильный идентификатор провайдера в реестре.
	ID() string

	// Bundles возвращает все bundles, которые провайдер хочет
	// зарегистрировать. Реестр сам пробежит по ним и положит под ID
	// провайдера; Provider и Bundle про реестр ничего не знают.
	Bundles() []Bundle
}

// NewRegistry собирает пустой реестр.
func NewRegistry() *Registry {
	return &Registry{
		bundles: make(map[bundleKey]Registered),
		byField: make(map[Field]Registered),
		fields:  make(map[Field]FieldDescriptor),
	}
}

// RegisterProvider регистрирует все bundles провайдера в реестре под
// ID этого провайдера. Двойная регистрация той же пары (provider, bundle)
// или того же id поля — ошибка: реестр строится один раз и подмена
// реализации в нём — почти всегда баг сборки.
func (r *Registry) RegisterProvider(p Provider) error {
	providerID := p.ID()
	for _, b := range p.Bundles() {
		if err := r.Register(providerID, b); err != nil {
			return fmt.Errorf("provider %s register: %w", providerID, err)
		}
	}
	return nil
}

// Register добавляет bundle в реестр под указанным providerID. Прямой
// путь записи мимо Provider — удобен в тестах реестра, где поднимать
// полноценный Provider избыточно.
func (r *Registry) Register(providerID string, b Bundle) error {
	entry := Registered{ProviderID: providerID, Bundle: b}
	if err := r.checkDuplicates(entry); err != nil {
		return err
	}
	r.add(entry)
	return nil
}

// checkDuplicates проверяет, что bundle и его поля ещё не зарегистрированы.
// Регистрация атомарна: если хотя бы одна проверка падает, реестр
// не изменяется. Поэтому проверка вынесена в отдельный шаг до add.
func (r *Registry) checkDuplicates(e Registered) error {
	bk := bundleKey{Provider: e.ProviderID, Bundle: e.Bundle.BundleID()}
	if _, exists := r.bundles[bk]; exists {
		return fmt.Errorf("%w: %s/%s", ErrBundleAlreadyRegistered, e.ProviderID, e.Bundle.BundleID())
	}
	for _, fd := range e.Bundle.Fields() {
		if _, exists := r.byField[fd.ID]; exists {
			return fmt.Errorf("%w: %s", ErrFieldAlreadyRegistered, fd.ID)
		}
	}
	return nil
}

// add вносит bundle и все его поля в индексы реестра. Вызывается только
// после успешного checkDuplicates, поэтому не возвращает ошибку.
func (r *Registry) add(e Registered) {
	r.bundles[bundleKey{Provider: e.ProviderID, Bundle: e.Bundle.BundleID()}] = e
	for _, fd := range e.Bundle.Fields() {
		r.byField[fd.ID] = e
		r.fields[fd.ID] = fd
	}
}

// FieldByID — O(1) lookup описания поля по идентификатору. Второй
// возвращаемый — false, если такого поля в реестре нет.
func (r *Registry) FieldByID(id Field) (FieldDescriptor, bool) {
	fd, ok := r.fields[id]
	return fd, ok
}

// Bundle ищет bundle по паре (provider, bundle).
// Возвращает ErrBundleNotFound, если такого bundle в реестре нет.
func (r *Registry) Bundle(provider, bundle string) (Bundle, error) {
	e, ok := r.bundles[bundleKey{Provider: provider, Bundle: bundle}]
	if !ok {
		return nil, fmt.Errorf("%w: %s/%s", ErrBundleNotFound, provider, bundle)
	}
	return e.Bundle, nil
}

// BundleFor ищет bundle, который отдаёт указанное поле. Возвращает
// ErrFieldNotFound, если такого поля в реестре нет.
func (r *Registry) BundleFor(fieldID Field) (Bundle, error) {
	e, ok := r.byField[fieldID]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrFieldNotFound, fieldID)
	}
	return e.Bundle, nil
}

// Fetch собирает значения для запрошенного списка полей.
//
// Реестр сам решает, какие bundles нужно вызвать: каждое поле
// разрешается в свой bundle через индекс byField; одинаковые bundles
// схлопываются, чтобы один источник не дёргался дважды на запрос.
// Все нужные bundles вызываются параллельно; первая ошибка отменяет
// контекст остальных и поднимается наверх как есть (sentinel-ошибки
// сохраняются для errors.Is). Итоговая карта содержит только запрошенные
// поля — лишние значения, которые bundle вернул заодно, отбрасываются.
//
// Возможные ошибки:
//   - ErrFieldNotFound, обёрнутый через fmt.Errorf, — реестр не знает
//     хотя бы одно из запрошенных полей; до Fetch ни один bundle
//     не вызывается;
//   - ошибка из Bundle.Fetch — пробрасывается с пометкой источника.
func (r *Registry) Fetch(ctx context.Context, ticker string, fieldIDs []Field) (FieldValues, error) {
	if len(fieldIDs) == 0 {
		return FieldValues{}, nil
	}

	entries, err := r.resolveBundles(fieldIDs)
	if err != nil {
		return nil, err
	}

	parts, err := r.fetchAll(ctx, ticker, entries)
	if err != nil {
		return nil, err
	}

	return filterFieldValues(parts, fieldIDs), nil
}

// resolveBundles переводит список fieldIDs в список уникальных
// зарегистрированных bundles. Каждый bundle попадает в результат
// ровно один раз, даже если из него запрошено несколько полей. Любое
// неизвестное поле — ошибка до запуска горутин: дешевле упасть на
// старте, чем поднимать пул и сетевые вызовы ради провального ответа.
func (r *Registry) resolveBundles(fieldIDs []Field) ([]Registered, error) {
	seen := make(map[bundleKey]bool, len(fieldIDs))
	entries := make([]Registered, 0, len(fieldIDs))
	for _, fieldID := range fieldIDs {
		e, ok := r.byField[fieldID]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrFieldNotFound, fieldID)
		}
		bk := bundleKey{Provider: e.ProviderID, Bundle: e.Bundle.BundleID()}
		if seen[bk] {
			continue
		}
		seen[bk] = true
		entries = append(entries, e)
	}
	return entries, nil
}

// fetchAll параллельно зовёт Fetch у всех bundles. Первая ошибка
// отменяет контекст остальных и возвращается как есть.
func (*Registry) fetchAll(ctx context.Context, ticker string, entries []Registered) ([]FieldValues, error) {
	p := pool.NewWithResults[FieldValues]().
		WithErrors().
		WithFirstError().
		WithContext(ctx).
		WithCancelOnError()
	for _, e := range entries {
		p.Go(func(ctx context.Context) (FieldValues, error) {
			values, err := e.Bundle.Fetch(ctx, ticker)
			if err != nil {
				return nil, fmt.Errorf("%s/%s: %w", e.ProviderID, e.Bundle.BundleID(), err)
			}
			return values, nil
		})
	}
	parts, err := p.Wait()
	if err != nil {
		return nil, fmt.Errorf("fetch bundles: %w", err)
	}
	return parts, nil
}

// filterFieldValues сводит partial-результаты в итоговую карту только
// по запрошенным fieldIDs. Bundles обязаны вернуть полный набор своих
// полей, но в ответ должны попасть только те, что попросил клиент.
func filterFieldValues(parts []FieldValues, fieldIDs []Field) FieldValues {
	merged := make(FieldValues, len(fieldIDs))
	for _, fieldID := range fieldIDs {
		for _, p := range parts {
			if v, ok := p[fieldID]; ok {
				merged[fieldID] = v
				break
			}
		}
	}
	return merged
}

// Bundles — все зарегистрированные bundles вместе с ID их провайдера.
// Порядок не определён. Используется в app/ для startup-сверки реестра
// со спецификацией.
func (r *Registry) Bundles() []Registered {
	out := make([]Registered, 0, len(r.bundles))
	for _, e := range r.bundles {
		out = append(out, e)
	}
	return out
}
