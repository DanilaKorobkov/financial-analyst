// Package services — domain-сервисы, оркеструют агрегаты и порты.
package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
)

// ErrTickerEmpty — клиент передал пустой тикер.
var ErrTickerEmpty = errors.New("ticker is empty")

// CompanyService собирает карточку эмитента по тикеру. Сам сервис
// не знает ни про источники, ни про каталог полей — состав ответа
// определяет ProfileRepository (по тикеру), а откуда брать значения —
// решает реестр bundles.
type CompanyService struct {
	profiles company.ProfileRepository
	registry *data.Registry
}

// ConfigCompanyService — параметры CompanyService.
type ConfigCompanyService struct {
	// Profiles — репозиторий профилей карточки.
	Profiles company.ProfileRepository
	// Registry — реестр зарегистрированных bundles.
	Registry *data.Registry
}

// NewCompanyService собирает сервис вокруг репозитория профилей и реестра.
func NewCompanyService(cfg ConfigCompanyService) *CompanyService {
	return &CompanyService{profiles: cfg.Profiles, registry: cfg.Registry}
}

// GetCompany проверяет непустоту тикера, спрашивает у репозитория
// профилей, какие поля нужны для этого тикера, и просит реестр собрать
// значения. Тикер передаётся как есть, без нормализации.
//
// Возможные ошибки:
//   - ErrTickerEmpty — пустой тикер;
//   - company.ErrProfileNotFound — для тикера не настроен профиль;
//   - ошибки реестра (включая data.ErrFieldNotFound для незнакомого
//     полю профиля поля) и bundles — пробрасываются с пометкой тикера.
func (s *CompanyService) GetCompany(ctx context.Context, ticker string) (data.FieldValues, error) {
	if ticker == "" {
		return nil, ErrTickerEmpty
	}

	profile, err := s.profiles.FindByTicker(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("resolve profile for %q: %w", ticker, err)
	}

	values, err := s.registry.Fetch(ctx, ticker, profile.FieldIDs)
	if err != nil {
		return nil, fmt.Errorf("get company %q: %w", ticker, err)
	}
	return values, nil
}
