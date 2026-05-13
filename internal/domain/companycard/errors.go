package companycard

import "errors"

// ErrNotFound — карточки эмитента по такому тикеру нет.
var ErrNotFound = errors.New("company card not found")
