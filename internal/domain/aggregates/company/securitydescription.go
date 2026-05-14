package company

import (
	"context"
	"time"
)

// SecurityDescription — описание ценной бумаги: идентификаторы выпуска,
// типизация, торговый листинг и режимы дополнительных сессий. Секция
// отдаётся одним источником целиком.
type SecurityDescription struct {
	// IssueDate — дата начала торгов бумагой. Нулевое значение = поле отсутствует.
	IssueDate time.Time

	// RegistryDate — дата государственной регистрации выпуска. Нулевое значение = отсутствует.
	RegistryDate time.Time

	// Ticker — биржевой код бумаги.
	Ticker string

	// ISIN — международный идентификатор ценной бумаги.
	ISIN string

	// Name — полное наименование выпуска.
	Name string

	// ShortName — краткое наименование выпуска.
	ShortName string

	// IssueName — наименование выпуска.
	IssueName string

	// LatName — латинское наименование выпуска.
	LatName string

	// RegNumber — номер государственной регистрации выпуска.
	RegNumber string

	// SecurityTypeName — человекочитаемое название типа бумаги.
	SecurityTypeName string

	// SecurityGroup — код группы инструмента.
	SecurityGroup string

	// SecurityGroupName — название группы инструмента.
	SecurityGroupName string

	// FaceValue — номинальная стоимость бумаги (как строка, без округления).
	FaceValue string

	// EmitterID — идентификатор эмитента в справочнике источника.
	EmitterID string

	// IssueSize — объём выпуска в штуках.
	IssueSize int64

	// SecurityType — тип бумаги.
	SecurityType SecurityType

	// ListingLevel — котировальный уровень бумаги.
	ListingLevel ListingLevel

	// FaceUnit — валюта номинала бумаги.
	FaceUnit Currency

	// HasProspectus — наличие зарегистрированного проспекта эмиссии.
	HasProspectus bool

	// HasDefault — наличие дефолта по выпуску.
	HasDefault bool

	// HasTechnicalDefault — наличие технического дефолта.
	HasTechnicalDefault bool

	// EmitentMismatchCurrent — эмитент не соответствует требованию на текущий список.
	EmitentMismatchCurrent bool

	// IsQualifiedInvestors — доступ только для квалифицированных инвесторов.
	IsQualifiedInvestors bool

	// MorningSession — допуск к утренней дополнительной сессии.
	MorningSession bool

	// EveningSession — допуск к вечерней дополнительной сессии.
	EveningSession bool

	// WeekendSession — допуск к сессии выходного дня.
	WeekendSession bool
}

// SecurityDescriptionSource — порт источника секции SecurityDescription.
type SecurityDescriptionSource interface {
	// FindByTicker возвращает описание бумаги по тикеру.
	// Возвращает (nil, ErrNotFound), если источник не знает тикер.
	FindByTicker(ctx context.Context, ticker string) (*SecurityDescription, error)
}
