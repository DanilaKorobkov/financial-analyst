package connect_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
)

type mapperSuite struct {
	suite.Suite
}

func TestMapperSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(mapperSuite))
}

// resolverByMap собирает резолвер типов поверх явной таблицы, чтобы
// тесты mapper'а не зависели от настоящего реестра.
func resolverByMap(table map[string]data.FieldType) func(id string) (data.FieldType, bool) {
	return func(id string) (data.FieldType, bool) {
		t, ok := table[id]
		return t, ok
	}
}

// TestToProtoFieldsUnknownField: если в FieldValues случайно затесался
// fieldID, которого нет в реестре, mapper падает с ошибкой —
// нормальный путь это исключает, но защита от регрессии нужна.
func (s *mapperSuite) TestToProtoFieldsUnknownField() {
	_, err := pconnect.PackFields(
		data.FieldValues{"unknown::field": "value"},
		resolverByMap(nil),
	)
	s.Require().Error(err)
	s.ErrorContains(err, "unknown::field")
}

// TestToProtoFieldsTypeMismatchWraps: ошибка из encode* оборачивается
// в сообщение с именем поля — оператору сразу видно, какое поле сломалось.
func (s *mapperSuite) TestToProtoFieldsTypeMismatchWraps() {
	_, err := pconnect.PackFields(
		data.FieldValues{company.FieldTicker: 42}, // ожидается string
		resolverByMap(map[string]data.FieldType{company.FieldTicker: data.TypeString}),
	)
	s.Require().Error(err)
	s.ErrorContains(err, company.FieldTicker)
}

// TestEncodeFieldValueMismatchErrors прогоняет каждую ветку oneof
// с заведомо неверным типом значения. Каноничный путь покрывается
// сценариями server_test.go; здесь — только ошибки, чтобы держать
// per-file пороги покрытия.
func (s *mapperSuite) TestEncodeFieldValueMismatchErrors() {
	cases := []struct {
		raw  any
		name string
		t    data.FieldType
	}{
		{raw: 1, name: "string", t: data.TypeString},
		{raw: "not-int", name: "int64", t: data.TypeInt64},
		{raw: "not-bool", name: "bool", t: data.TypeBool},
		{raw: "not-time", name: "date", t: data.TypeDate},
		{raw: "not-enum", name: "security_type", t: data.TypeSecurityType},
		{raw: "not-enum", name: "listing_level", t: data.TypeListingLevel},
		{raw: "not-enum", name: "currency", t: data.TypeCurrency},
		{raw: "not-enum", name: "exchange", t: data.TypeExchange},
		{raw: "not-enum", name: "report_frequency", t: data.TypeReportFrequency},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			_, err := pconnect.PackFieldValue(c.t, c.raw)
			s.Require().Error(err)
		})
	}
}

// TestEncodeFieldValueUnknownType — несуществующая ветка enum FieldType
// (например, добавленная в каталог, но забытая в mapper) приводит к ошибке.
func (s *mapperSuite) TestEncodeFieldValueUnknownType() {
	_, err := pconnect.PackFieldValue(data.FieldType(999), "value")
	s.Require().Error(err)
}

// TestEncodeFieldValueDateZeroSkipped и Unspecified покрывают «пропустить
// поле» (nil, nil): значение нулевое, ошибки нет.
func (s *mapperSuite) TestEncodeFieldValueDateZeroSkipped() {
	got, err := pconnect.PackFieldValue(data.TypeDate, time.Time{})
	s.Require().NoError(err)
	s.Nil(got)
}

func (s *mapperSuite) TestEncodeFieldValueEnumUnspecifiedSkipped() {
	got, err := pconnect.PackFieldValue(data.TypeCurrency, company.CurrencyUnspecified)
	s.Require().NoError(err)
	s.Nil(got)
}
