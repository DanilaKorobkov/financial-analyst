package entities

import "errors"

// ErrCompanyNotFound — справочник не нашёл компанию по тикеру.
var ErrCompanyNotFound = errors.New("company not found")

// ErrNotFound — внешний источник не нашёл запрошенный объект по идентификатору.
var ErrNotFound = errors.New("not found")

// ErrUnauthorized — внешний источник отверг запрос: токен не передан, невалиден
// либо просрочен.
var ErrUnauthorized = errors.New("unauthorized")

// ErrQuotaExceeded — внешний источник отверг запрос: исчерпана квота доступа
// (дневной лимит, подписка).
var ErrQuotaExceeded = errors.New("quota exceeded")
