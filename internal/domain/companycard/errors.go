package companycard

import "errors"

// ErrNotFound — карточки эмитента по такой паре биржа/тикер нет.
var ErrNotFound = errors.New("company card not found")
