package connect

import (
	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
)

// PackFields — внутренний адаптер FieldValues → map proto-полей,
// открытый для тестов соседнего пакета.
func PackFields(
	values data.FieldValues,
	fieldType func(id string) (data.FieldType, bool),
) (map[string]*companyv1.FieldValue, error) {
	return toProtoFields(values, fieldType)
}

// PackFieldValue — внутренний выбор ветки oneof по типу поля,
// открытый для тестов соседнего пакета.
func PackFieldValue(t data.FieldType, raw any) (*companyv1.FieldValue, error) {
	return encodeFieldValue(t, raw)
}
