package entities

import "errors"

// ErrMissingCompany — справочник не нашёл компанию по тикеру.
var ErrMissingCompany = errors.New("missing company")
