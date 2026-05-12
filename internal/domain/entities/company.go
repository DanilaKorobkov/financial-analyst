// Package entities — плоские domain-сущности (без подпапок на агрегат).
package entities

import "time"

// Company — справочная карточка эмитента MOEX. Поля из блока description
// эндпоинта /iss/securities/{TICKER}.json. Примеры значений — по SBER.
type Company struct {
	// IssueDate — дата начала торгов на MOEX. Пример: 2007-07-20.
	IssueDate time.Time

	// Ticker — биржевой код бумаги (SECID). Пример: "SBER".
	Ticker string

	// ISIN — международный идентификатор. Пример: "RU0009029540".
	ISIN string

	// Name — полное название бумаги. Пример: "Сбербанк России ПАО ао".
	Name string

	// ShortName — краткое название. Пример: "Сбербанк".
	ShortName string

	// RegNumber — номер государственной регистрации выпуска. Пример: "10301481B".
	RegNumber string

	// SecurityType — код типа бумаги (поле TYPE блока description).
	// Пример: "common_share" / "preferred_share" / "depositary_receipt".
	SecurityType string

	// Group — код группы бумаг MOEX. Пример: "stock_shares".
	Group string

	// FaceUnit — валюта номинала. Пример: "SUR" (₽), "USD".
	FaceUnit string

	// IssueSize — объём выпуска в штуках. Пример: 21586948000.
	IssueSize int64

	// EmitterID — внутренний ID эмитента в MOEX. Пример: 484.
	EmitterID int64

	// FaceValue — номинальная стоимость в FaceUnit. Пример: 3.0.
	FaceValue float64

	// ListingLevel — котировальный уровень MOEX (1, 2 или 3).
	// 0 = биржа не указала (в proto-ответе — отсутствие optional-поля).
	// Пример: 1.
	ListingLevel int

	// Sessions — допуски бумаги к дополнительным торговым сессиям MOEX.
	// Для SBER: morning=true, evening=true, weekend=true.
	Sessions Sessions
}

// Sessions — допуски бумаги к дополнительным торговым сессиям MOEX.
type Sessions struct {
	// Morning — допуск к утренней сессии. Пример (SBER): true.
	Morning bool

	// Evening — допуск к вечерней сессии. Пример (SBER): true.
	Evening bool

	// Weekend — допуск к сессии выходного дня. Пример (SBER): true.
	Weekend bool
}
