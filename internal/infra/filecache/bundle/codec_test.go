package bundle_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	fcbundle "github.com/DanilaKorobkov/financial-analyst/internal/infra/filecache/bundle"
)

type codecSuite struct {
	suite.Suite
}

func TestCodecSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(codecSuite))
}

// TestEncodeDecodeAllFieldTypes гоняет round-trip по всем девяти
// FieldType-веткам сразу. После encode→decode значения и их Go-типы
// должны совпасть с исходными — иначе теряется типизация (как было бы
// при обычном json.Unmarshal в map[string]any).
func (s *codecSuite) TestEncodeDecodeAllFieldTypes() {
	fields := []data.FieldDescriptor{
		{ID: company.FieldTicker, Type: data.TypeString},
		{ID: company.FieldIssueSize, Type: data.TypeInt64},
		{ID: company.FieldHasProspectus, Type: data.TypeBool},
		{ID: company.FieldIssueDate, Type: data.TypeDate},
		{ID: company.FieldSecurityType, Type: data.TypeSecurityType},
		{ID: company.FieldListingLevel, Type: data.TypeListingLevel},
		{ID: company.FieldFaceUnit, Type: data.TypeCurrency},
		{ID: company.FieldExchange, Type: data.TypeExchange},
		{ID: company.FieldReportFrequency, Type: data.TypeReportFrequency},
	}
	codec := fcbundle.NewValueCodec(fields)

	values := data.FieldValues{
		company.FieldTicker:          "SBER",
		company.FieldIssueSize:       int64(21586948000),
		company.FieldHasProspectus:   true,
		company.FieldIssueDate:       time.Date(2007, 7, 20, 0, 0, 0, 0, time.UTC),
		company.FieldSecurityType:    company.SecurityTypeCommonShare,
		company.FieldListingLevel:    company.ListingLevelFirst,
		company.FieldFaceUnit:        company.CurrencyRUB,
		company.FieldExchange:        company.ExchangeMOEX,
		company.FieldReportFrequency: company.ReportFrequencyQuarterly,
	}

	encoded, err := codec.Encode(values)
	s.Require().NoError(err)

	decoded, err := codec.Decode(encoded)
	s.Require().NoError(err)
	s.Equal(values, decoded)
}

func (s *codecSuite) TestDecodeBadJSONForField() {
	fields := []data.FieldDescriptor{{ID: company.FieldIssueSize, Type: data.TypeInt64}}
	codec := fcbundle.NewValueCodec(fields)

	// число ожидалось, пришла строка — Unmarshal вернёт ошибку,
	// которую codec обернёт именем поля.
	raw := map[string]json.RawMessage{company.FieldIssueSize: json.RawMessage(`"not a number"`)}
	_, err := codec.Decode(raw)

	s.Require().Error(err)
	s.ErrorContains(err, company.FieldIssueSize)
}

// TestDecodeSkipsMissingFields: если значение поля отсутствует в
// raw-карте (например, кеш частичный), codec не падает — просто
// не кладёт ключ в результат.
func (s *codecSuite) TestDecodeSkipsMissingFields() {
	fields := []data.FieldDescriptor{
		{ID: company.FieldTicker, Type: data.TypeString},
		{ID: company.FieldIssueSize, Type: data.TypeInt64},
	}
	codec := fcbundle.NewValueCodec(fields)

	raw := map[string]json.RawMessage{company.FieldTicker: json.RawMessage(`"SBER"`)}
	out, err := codec.Decode(raw)

	s.Require().NoError(err)
	s.Equal(data.FieldValues{company.FieldTicker: "SBER"}, out)
}

// TestDecodeUnsupportedFieldType: если в FieldDescriptor проставлен
// несуществующий FieldType (например, добавили константу в data,
// но забыли ветку в decodeByType) — codec возвращает ошибку.
func (s *codecSuite) TestDecodeUnsupportedFieldType() {
	fields := []data.FieldDescriptor{{ID: company.FieldTicker, Type: data.FieldType(9999)}}
	codec := fcbundle.NewValueCodec(fields)

	raw := map[string]json.RawMessage{company.FieldTicker: json.RawMessage(`"x"`)}
	_, err := codec.Decode(raw)

	s.Require().Error(err)
	s.ErrorContains(err, "unsupported field type")
}
