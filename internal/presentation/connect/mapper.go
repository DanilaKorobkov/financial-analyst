// Package connect — Connect-handler для CompanyService.
package connect

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
)

// Стабильные строковые коды доменных enum-ов. Менять значения нельзя —
// это часть внешнего контракта (попадает в ответ клиента).
// *Unspecified трактуется как «нет значения» — пустая строка, поле
// выпадает из map в encodeEnum.
var (
	securityTypeCode = map[company.SecurityType]string{
		company.SecurityTypeUnspecified:       "",
		company.SecurityTypeCommonShare:       "common_share",
		company.SecurityTypePreferredShare:    "preferred_share",
		company.SecurityTypeDepositaryReceipt: "depositary_receipt",
	}
	listingLevelCode = map[company.ListingLevel]string{
		company.ListingLevelUnspecified: "",
		company.ListingLevelFirst:       "first",
		company.ListingLevelSecond:      "second",
		company.ListingLevelThird:       "third",
	}
	currencyCode = map[company.Currency]string{
		company.CurrencyUnspecified: "",
		company.CurrencyRUB:         "RUB",
		company.CurrencyUSD:         "USD",
		company.CurrencyEUR:         "EUR",
	}
	exchangeCode = map[company.Exchange]string{
		company.ExchangeUnspecified: "",
		company.ExchangeMOEX:        "moex",
	}
	reportFrequencyCode = map[company.ReportFrequency]string{
		company.ReportFrequencyUnspecified: "",
		company.ReportFrequencyYearly:      "yearly",
		company.ReportFrequencyQuarterly:   "quarterly",
	}
)

// fieldTypeResolver — поиск типа поля по id. Сервер строит его поверх
// data.Registry; mapper остаётся независимым от того, откуда берётся
// метаданные.
type fieldTypeResolver func(id data.Field) (data.FieldType, bool)

// toProtoFields упаковывает значения, собранные сервисом, в карту полей
// proto-ответа. Тип каждой ветки oneof выбирается по типу поля,
// найденному в реестре.
// «Значения нет» (нулевая дата, *Unspecified) — поля в map не будет,
// а не пустой строки/нулевой ветки.
func toProtoFields(values data.FieldValues, fieldType fieldTypeResolver) (map[string]*companyv1.FieldValue, error) {
	out := make(map[string]*companyv1.FieldValue, len(values))
	for fieldID, raw := range values {
		t, ok := fieldType(fieldID)
		if !ok {
			return nil, fmt.Errorf("field %q is not in registry", fieldID)
		}
		fv, err := encodeFieldValue(t, raw)
		if err != nil {
			return nil, fmt.Errorf("encode %s: %w", fieldID, err)
		}
		if fv == nil {
			continue
		}
		out[string(fieldID)] = fv
	}
	return out, nil
}

// encodeFieldValue выбирает ветку oneof по типу поля. Возвращает (nil, nil),
// если значение трактуется как «отсутствует» — для нулевого time.Time
// и для *Unspecified enum-ов.
func encodeFieldValue(t data.FieldType, raw any) (*companyv1.FieldValue, error) {
	switch t {
	case data.TypeString:
		return encodeString(raw)
	case data.TypeInt64:
		return encodeInt64(raw)
	case data.TypeBool:
		return encodeBool(raw)
	case data.TypeDate:
		return encodeDate(raw)
	case data.TypeSecurityType:
		return encodeEnum(raw, securityTypeCode)
	case data.TypeListingLevel:
		return encodeEnum(raw, listingLevelCode)
	case data.TypeCurrency:
		return encodeEnum(raw, currencyCode)
	case data.TypeExchange:
		return encodeEnum(raw, exchangeCode)
	case data.TypeReportFrequency:
		return encodeEnum(raw, reportFrequencyCode)
	default:
		return nil, fmt.Errorf("unsupported field type %d", t)
	}
}

func encodeString(raw any) (*companyv1.FieldValue, error) {
	v, ok := raw.(string)
	if !ok {
		return nil, fmt.Errorf("expected string, got %T", raw)
	}
	return stringValue(v), nil
}

func encodeInt64(raw any) (*companyv1.FieldValue, error) {
	v, ok := raw.(int64)
	if !ok {
		return nil, fmt.Errorf("expected int64, got %T", raw)
	}
	return &companyv1.FieldValue{Value: &companyv1.FieldValue_IntValue{IntValue: v}}, nil
}

func encodeBool(raw any) (*companyv1.FieldValue, error) {
	v, ok := raw.(bool)
	if !ok {
		return nil, fmt.Errorf("expected bool, got %T", raw)
	}
	return &companyv1.FieldValue{Value: &companyv1.FieldValue_BoolValue{BoolValue: v}}, nil
}

// encodeDate переводит time.Time в Timestamp. Нулевой time.Time
// (поле не отдано источником) превращается в (nil, nil) — поле выпадает
// из map, что соответствует семантике proto-3 «значения нет».
//
//nolint:nilnil // (nil, nil) — каноничный «пропустить поле» для динамического контракта.
func encodeDate(raw any) (*companyv1.FieldValue, error) {
	v, ok := raw.(time.Time)
	if !ok {
		return nil, fmt.Errorf("expected time.Time, got %T", raw)
	}
	if v.IsZero() {
		return nil, nil
	}
	return &companyv1.FieldValue{Value: &companyv1.FieldValue_TimestampValue{TimestampValue: timestamppb.New(v)}}, nil
}

// encodeEnum переводит доменный enum в string_value по таблице стабильных
// кодов. Пустой код в таблице (стандартный маркер *Unspecified) или
// отсутствие ключа — это «нет значения», поле выпадает из map.
//
//nolint:nilnil // (nil, nil) — каноничный «пропустить поле» для динамического контракта.
func encodeEnum[T comparable](raw any, table map[T]string) (*companyv1.FieldValue, error) {
	v, ok := raw.(T)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", *new(T), raw)
	}
	code, found := table[v]
	if !found || code == "" {
		return nil, nil
	}
	return stringValue(code), nil
}

func stringValue(v string) *companyv1.FieldValue {
	return &companyv1.FieldValue{Value: &companyv1.FieldValue_StringValue{StringValue: v}}
}
