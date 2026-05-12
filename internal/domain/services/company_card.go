package services

import (
	"context"
	"fmt"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// CompanyCard — сервис поиска карточки эмитента (классификация, описание,
// ссылки).
type CompanyCard struct {
	cards entities.CompanyCardRepository
}

// NewCompanyCard собирает сервис вокруг репозитория карточек эмитентов.
func NewCompanyCard(cards entities.CompanyCardRepository) *CompanyCard {
	return &CompanyCard{cards: cards}
}

// FindByTicker проверяет непустоту тикера и делегирует поиск репозиторию.
// Тикер передаётся как есть, без нормализации.
func (s *CompanyCard) FindByTicker(
	ctx context.Context,
	ticker string,
) (entities.CompanyCard, error) {
	if ticker == "" {
		return entities.CompanyCard{}, ErrTickerEmpty
	}
	card, err := s.cards.FindByTicker(ctx, ticker)
	if err != nil {
		return entities.CompanyCard{}, fmt.Errorf("find card %q: %w", ticker, err)
	}
	return card, nil
}
