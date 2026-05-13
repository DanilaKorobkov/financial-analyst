package company

import "errors"

// ErrNotFound — компании с таким тикером нет в справочнике.
var ErrNotFound = errors.New("company not found")
