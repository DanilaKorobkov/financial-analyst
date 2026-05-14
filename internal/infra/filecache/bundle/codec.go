package bundle

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
)

// valueCodec кодирует/декодирует FieldValues с учётом типов из
// FieldDescriptor каждого поля bundle. Без типизации обычный
// `json.Unmarshal` в `map[string]any` превращает int64 в float64,
// enum-ы в float64, time.Time в string — и round-trip ломается.
// Codec для каждого поля знает целевой Go-тип и unmarshal-ит raw-JSON
// именно в него.
type valueCodec struct {
	fields []data.FieldDescriptor
}

// newValueCodec собирает codec для конкретного bundle. Список полей
// замораживается на время жизни codec'а: добавление поля в bundle
// требует пересоздания Proxy.
func newValueCodec(fields []data.FieldDescriptor) *valueCodec {
	return &valueCodec{fields: fields}
}

// Encode маршалит каждое значение FieldValues по полю bundle в отдельный
// json.RawMessage. Это даёт стабильный набор ключей в выходном объекте
// (по списку Fields) и аккуратный точечный разбор на чтении.
func (c *valueCodec) Encode(values data.FieldValues) (map[string]json.RawMessage, error) {
	out := make(map[string]json.RawMessage, len(c.fields))
	for _, fd := range c.fields {
		raw, err := jsonParser.Marshal(values[fd.ID])
		if err != nil {
			return nil, fmt.Errorf("marshal %s: %w", fd.ID, err)
		}
		out[string(fd.ID)] = raw
	}
	return out, nil
}

// Decode восстанавливает FieldValues из map[id]raw, типизируя значение
// каждого поля согласно FieldDescriptor.Type. Поля, отсутствующие в
// raw-карте, в результат не попадают.
func (c *valueCodec) Decode(raw map[string]json.RawMessage) (data.FieldValues, error) {
	out := make(data.FieldValues, len(c.fields))
	for _, fd := range c.fields {
		rawBytes, ok := raw[string(fd.ID)]
		if !ok {
			continue
		}
		v, err := decodeByType(rawBytes, fd.Type)
		if err != nil {
			return nil, fmt.Errorf("decode %s: %w", fd.ID, err)
		}
		out[fd.ID] = v
	}
	return out, nil
}

// decodeByType разбирает raw в конкретный Go-тип, отвечающий
// data.FieldType. Список ветвей синхронизирован с FieldType-набором
// в `internal/domain/data`; добавление нового типа в data требует
// синхронной правки этой функции — расхождение ловится тестом
// TestDecodeUnknownType.
func decodeByType(raw json.RawMessage, t data.FieldType) (any, error) {
	switch t {
	case data.TypeString:
		return decodeInto[string](raw)
	case data.TypeInt64:
		return decodeInto[int64](raw)
	case data.TypeBool:
		return decodeInto[bool](raw)
	case data.TypeDate:
		return decodeInto[time.Time](raw)
	case data.TypeSecurityType:
		return decodeInto[company.SecurityType](raw)
	case data.TypeListingLevel:
		return decodeInto[company.ListingLevel](raw)
	case data.TypeCurrency:
		return decodeInto[company.Currency](raw)
	case data.TypeExchange:
		return decodeInto[company.Exchange](raw)
	case data.TypeReportFrequency:
		return decodeInto[company.ReportFrequency](raw)
	default:
		return nil, fmt.Errorf("unsupported field type: %d", t)
	}
}

// decodeInto — generic-обёртка вокруг json.Unmarshal: даёт каждой ветке
// decodeByType один шаблонный однострочник вместо трёх строк
// объявление/разбор/возврат.
func decodeInto[T any](raw json.RawMessage) (any, error) {
	var v T
	if err := jsonParser.Unmarshal(raw, &v); err != nil {
		return nil, err //nolint:wrapcheck // ошибка обёрнута выше в decode по имени поля
	}
	return v, nil
}
