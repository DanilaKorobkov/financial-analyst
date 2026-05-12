package moex

import (
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

var (
	// jsonParser — JSON-парсер с поведением, идентичным encoding/json:
	// sorted map keys, html-escape, validateRawMessage. См. rules/golang.md.
	jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary

	// errDescriptionBlockMissing — payload не содержит ожидаемого блока description.
	errDescriptionBlockMissing = fmt.Errorf("description block missing in payload")
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

// parseDescription разбирает блок description как map[name]value, без сохранения порядка.
func parseDescription(raw []byte) (map[string]string, error) {
	var payload extendedPayload
	if err := jsonParser.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("decode extended payload: %w", err)
	}

	for _, block := range payload {
		rawDesc, ok := block["description"]
		if !ok {
			continue
		}
		var fields []descriptionField
		if err := jsonParser.Unmarshal(rawDesc, &fields); err != nil {
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
		return entities.Company{}, fmt.Errorf("empty description block: %w", entities.ErrCompanyNotFound)
	}

	listLevel, err := parseListingLevel(fields["LISTLEVEL"])
	if err != nil {
		return entities.Company{}, fmt.Errorf("LISTLEVEL: %w", err)
	}

	return entities.Company{
		Ticker:       fields["SECID"],
		ISIN:         fields["ISIN"],
		Name:         fields["NAME"],
		SecurityType: parseSecurityType(fields["TYPE"]),
		ListingLevel: listLevel,
	}, nil
}

// parseSecurityType переводит код TYPE блока description MOEX в domain-enum.
// Неизвестные значения становятся SecurityTypeUnspecified — TYPE задаётся
// биржей и со временем расширяется (фонды, облигации и т.п.).
func parseSecurityType(s string) entities.SecurityType {
	switch s {
	case "common_share":
		return entities.SecurityTypeCommonShare
	case "preferred_share":
		return entities.SecurityTypePreferredShare
	case "depositary_receipt":
		return entities.SecurityTypeDepositaryReceipt
	default:
		return entities.SecurityTypeUnspecified
	}
}

// parseListingLevel переводит строковое значение LISTLEVEL ("1" / "2" / "3" / "")
// в entities.ListingLevel. Пустое значение — биржа не указала уровень.
// Любое другое значение — ошибка: список уровней зафиксирован MOEX.
func parseListingLevel(s string) (entities.ListingLevel, error) {
	switch s {
	case "":
		return entities.ListingLevelUnspecified, nil
	case "1":
		return entities.ListingLevelFirst, nil
	case "2":
		return entities.ListingLevelSecond, nil
	case "3":
		return entities.ListingLevelThird, nil
	default:
		return entities.ListingLevelUnspecified, fmt.Errorf("unexpected LISTLEVEL value: %q", s)
	}
}
