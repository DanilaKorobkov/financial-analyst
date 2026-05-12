package entities

import "errors"

// ErrCompanyNotFound — справочник не нашёл компанию по тикеру.
var ErrCompanyNotFound = errors.New("company not found")

// ErrNotFound — внешний источник не нашёл запрошенный объект по идентификатору.
var ErrNotFound = errors.New("not found")
