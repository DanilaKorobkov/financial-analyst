package company

import "errors"

// ErrNotFound — компании по такому тикеру нет.
var ErrNotFound = errors.New("company not found")
