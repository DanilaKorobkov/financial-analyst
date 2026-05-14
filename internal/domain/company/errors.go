// Package company — domain-аггрегат «компания» (эмитент бумаги):
// идентификаторы каноничных полей, доменные enum-ы (тип бумаги,
// котировальный уровень, валюта, биржа, частота отчётности) и
// доменные ошибки. Сам ответ собирается как набор каноничных полей
// (data.FieldValues) в реестре поверх bundles; типизированный
// агрегат в этом пакете не живёт — это намеренный выбор.
package company

import "errors"

var (
	// ErrNotFound — компании по такому тикеру нет.
	ErrNotFound = errors.New("company not found")

	// ErrProfileNotFound — для тикера не настроен профиль карточки.
	ErrProfileNotFound = errors.New("company profile not found")
)
