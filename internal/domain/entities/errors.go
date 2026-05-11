package entities

import "errors"

// ErrCompanyNotFound — справочник не нашёл компанию по тикеру.
var ErrCompanyNotFound = errors.New("company not found")
