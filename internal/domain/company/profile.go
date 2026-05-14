package company

import (
	"context"
	"errors"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
)

// ErrProfileNotFound — для тикера не настроен профиль карточки.
var ErrProfileNotFound = errors.New("company profile not found")

// Profile — описание карточки эмитента: какие каноничные поля нужны
// в ответе для конкретного тикера. Состав может различаться по тикеру —
// например, у будущих архетипов «банк», «нефтянка», «телеком» свои
// наборы метрик. Сейчас порт реализован статически, в перспективе —
// чтение из конфигурации пользователя.
type Profile struct {
	// FieldIDs — каноничные идентификаторы полей.
	// Порядок не значим: реестр сам решает, какие bundles вызвать,
	// и параллелит запросы.
	FieldIDs []data.Field
}

// ProfileRepository — порт получения профиля карточки по тикеру.
type ProfileRepository interface {
	// FindByTicker возвращает профиль для запрошенного тикера.
	// Тикер передаётся как есть, без нормализации.
	//
	// Возвращает ErrProfileNotFound, если для тикера нет настроенного
	// профиля.
	FindByTicker(ctx context.Context, ticker string) (Profile, error)
}
