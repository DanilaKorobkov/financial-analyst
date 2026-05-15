package company

import "errors"

// ErrNotFound — компании (или одной из её секций) по такому тикеру нет.
// Возвращается источниками секций как маркер «нет данных» и репозиторием —
// когда хотя бы одна секция не нашлась. Транспортный слой переводит её
// в 404 / CodeNotFound.
var ErrNotFound = errors.New("company not found")
