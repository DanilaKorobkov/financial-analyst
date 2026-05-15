package company

import "time"

// StockIdea — инвест-идея аналитика по бумаге. Поля статуса, цены закрытия
// и обновления могут быть нулевыми — это означает, что событие у идеи ещё
// не наступило (например, ещё нет факта закрытия или обновления).
type StockIdea struct {
	// DateIn — дата публикации идеи (вход).
	DateIn time.Time

	// DateOut — плановая дата выхода (когда автор ожидает достижения цели).
	DateOut time.Time

	// CloseDate — дата фактического закрытия идеи или последнего снимка цены.
	CloseDate time.Time

	// UpdateDate — дата последнего обновления (изменение target/stop).
	UpdateDate time.Time

	// ChangedAt — момент последнего изменения записи в FM.
	ChangedAt time.Time

	// Community — имя автора-аналитика или брокера.
	Community string

	// Idea — краткий заголовок идеи.
	Idea string

	// CloseComment — комментарий автора при закрытии идеи.
	CloseComment string

	// CloseLink — ссылка на пост о закрытии идеи.
	CloseLink string

	// ID — идентификатор идеи в FM.
	ID int64

	// CommunityID — числовой ID автора-аналитика.
	CommunityID int64

	// DurationInMonth — срок идеи в полных месяцах (DateOut − DateIn).
	DurationInMonth int64

	// PriceIn — цена входа, заявленная автором.
	PriceIn float64

	// PriceOut — целевая цена.
	PriceOut float64

	// PriceDay — текущая цена бумаги на момент снимка ленты.
	PriceDay float64

	// ProfitPotential — потенциал к target от PriceIn (%).
	ProfitPotential float64

	// ProfitActual — фактическая доходность от PriceIn к PriceDay (%).
	ProfitActual float64

	// StopLoss — уровень stop-loss, если задан автором.
	StopLoss float64

	// ClosePrice — цена закрытия / последнего снимка идеи.
	ClosePrice float64

	// UpdatePrice — цена бумаги в момент обновления идеи.
	UpdatePrice float64

	// Status — статус идеи (активна, закрыта).
	Status IdeaStatus
}
