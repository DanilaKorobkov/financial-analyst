// Package companyprofile — реализации company.ProfileRepository.
//
// Static — единственный профиль на все тикеры: те же поля карточки
// для каждой бумаги. Это переходный шаг, пока пользователь не сможет
// собирать свой профиль через GUI; тогда здесь же появится реализация,
// читающая профиль из конфигурации/БД.
package companyprofile

import (
	"context"

	"github.com/samber/lo"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
)

// Static — статический ProfileRepository: один и тот же набор полей
// возвращается на любой непустой тикер.
type Static struct {
	profile company.Profile
}

// NewStatic собирает Static из произвольного набора каноничных id полей.
// Список копируется, чтобы внешние изменения не протекли в репозиторий.
func NewStatic(fieldIDs []string) *Static {
	return &Static{profile: company.Profile{FieldIDs: lo.Clone(fieldIDs)}}
}

// FindByTicker возвращает один и тот же профиль независимо от тикера.
// Тикер игнорируется намеренно — это поведение Static-реализации.
func (s *Static) FindByTicker(_ context.Context, _ string) (company.Profile, error) {
	return s.profile, nil
}
