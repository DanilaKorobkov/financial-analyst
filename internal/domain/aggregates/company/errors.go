package company

import "errors"

var (
	// ErrNotFound — источник секции не нашёл бумагу по тикеру. Возвращается
	// реализациями <Section>Source как маркер «нет данных», репозиторий
	// различает её, чтобы решить, считать ли всю компанию ненайденной.
	ErrNotFound = errors.New("section not found")

	// ErrCompanyNotFound — репозиторий не собрал ни одной секции: все
	// источники ответили ErrNotFound. Транспортный слой превращает её
	// в 404 / CodeNotFound.
	ErrCompanyNotFound = errors.New("company not found")
)
