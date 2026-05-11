package moex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// descriptionField — одна строка блока description в ответе ISS.
//
// Значение ISS всегда отдаёт строкой (поле type указывает исходный тип на
// стороне биржи). См. .claude/skills/api-moex/references/find_security_description.md.
type descriptionField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// extendedPayload — формат iss.json=extended: массив [meta, blocks].
type extendedPayload []map[string]json.RawMessage

// errDescriptionBlockMissing — payload не содержит ожидаемого блока description.
var errDescriptionBlockMissing = fmt.Errorf("description block missing in payload")

// parseDescription разбирает блок description как map[name]value, без сохранения порядка.
func parseDescription(raw []byte) (map[string]string, error) {
	var payload extendedPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("decode extended payload: %w", err)
	}

	for _, block := range payload {
		rawDesc, ok := block["description"]
		if !ok {
			continue
		}
		var fields []descriptionField
		if err := json.Unmarshal(rawDesc, &fields); err != nil {
			return nil, fmt.Errorf("decode description block: %w", err)
		}
		out := make(map[string]string, len(fields))
		for _, f := range fields {
			out[f.Name] = f.Value
		}
		return out, nil
	}

	return nil, errDescriptionBlockMissing
}

// mapDescription собирает entities.Company из map[name]value блока description.
// Пустой map → entities.ErrCompanyNotFound (тикер не найден).
func mapDescription(fields map[string]string) (entities.Company, error) {
	if len(fields) == 0 {
		return entities.Company{}, entities.ErrCompanyNotFound
	}

	issueSize, err := parseInt64(fields["ISSUESIZE"])
	if err != nil {
		return entities.Company{}, fmt.Errorf("ISSUESIZE: %w", err)
	}
	faceValue, err := parseFloat64(fields["FACEVALUE"])
	if err != nil {
		return entities.Company{}, fmt.Errorf("FACEVALUE: %w", err)
	}
	issueDate, err := parseDate(fields["ISSUEDATE"])
	if err != nil {
		return entities.Company{}, fmt.Errorf("ISSUEDATE: %w", err)
	}
	listLevel, err := parseInt(fields["LISTLEVEL"])
	if err != nil {
		return entities.Company{}, fmt.Errorf("LISTLEVEL: %w", err)
	}
	emitterID, err := parseInt64(fields["EMITTER_ID"])
	if err != nil {
		return entities.Company{}, fmt.Errorf("EMITTER_ID: %w", err)
	}

	return entities.Company{
		Ticker:       fields["SECID"],
		ISIN:         fields["ISIN"],
		Name:         fields["NAME"],
		ShortName:    fields["SHORTNAME"],
		RegNumber:    fields["REGNUMBER"],
		SecurityType: fields["TYPE"],
		Group:        fields["GROUP"],
		IssueSize:    issueSize,
		FaceValue:    faceValue,
		FaceUnit:     fields["FACEUNIT"],
		IssueDate:    issueDate,
		ListingLevel: listLevel,
		Sessions: entities.Sessions{
			Morning: parseBool(fields["MORNINGSESSION"]),
			Evening: parseBool(fields["EVENINGSESSION"]),
			Weekend: parseBool(fields["WEEKENDSESSION"]),
		},
		EmitterID: emitterID,
	}, nil
}

func parseInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

func parseInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func parseFloat64(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

func parseBool(s string) bool {
	return s == "1"
}

func parseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.ParseInLocation("2006-01-02", s, time.UTC)
}
