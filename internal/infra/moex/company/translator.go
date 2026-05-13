package company

import (
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
)

var (
	// jsonParser — JSON-парсер с поведением, идентичным encoding/json:
	// sorted map keys, html-escape, validateRawMessage. См. rules/golang.md.
	jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary

	// errDescriptionBlockMissing — payload не содержит ожидаемого блока description.
	errDescriptionBlockMissing = fmt.Errorf("description block missing in payload")

	// securityTypeByCode переводит код TYPE блока description MOEX в domain-enum.
	// Неизвестные значения возвращают zero (SecurityTypeUnspecified) — TYPE
	// задаётся биржей и со временем расширяется (фонды, облигации и т.п.).
	securityTypeByCode = map[string]domaincompany.SecurityType{
		"common_share":       domaincompany.SecurityTypeCommonShare,
		"preferred_share":    domaincompany.SecurityTypePreferredShare,
		"depositary_receipt": domaincompany.SecurityTypeDepositaryReceipt,
	}

	// listingLevelByCode переводит строковое значение LISTLEVEL ("1" / "2" / "3")
	// в domain-enum. Пустое значение — биржа не указала уровень (Unspecified).
	// Любое другое значение в parseListingLevel — ошибка: список уровней
	// зафиксирован MOEX.
	listingLevelByCode = map[string]domaincompany.ListingLevel{
		"":  domaincompany.ListingLevelUnspecified,
		"1": domaincompany.ListingLevelFirst,
		"2": domaincompany.ListingLevelSecond,
		"3": domaincompany.ListingLevelThird,
	}
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

// mapDescription собирает domaincompany.Identity из map[name]value блока description.
// Пустой map → domaincompany.ErrNotFound (тикер не найден).
func mapDescription(fields map[string]string) (domaincompany.Identity, error) {
	if len(fields) == 0 {
		return domaincompany.Identity{}, fmt.Errorf("empty description block: %w", domaincompany.ErrNotFound)
	}

	listLevel, err := parseListingLevel(fields["LISTLEVEL"])
	if err != nil {
		return domaincompany.Identity{}, fmt.Errorf("LISTLEVEL: %w", err)
	}

	return domaincompany.Identity{
		Ticker:       fields["SECID"],
		ISIN:         fields["ISIN"],
		Name:         fields["NAME"],
		SecurityType: securityTypeByCode[fields["TYPE"]],
		ListingLevel: listLevel,
	}, nil
}

// parseListingLevel — единственная функция, где нужна валидация (неизвестный
// уровень — ошибка, а не Unspecified): список уровней зафиксирован MOEX и
// должен быть пополнен сознательно, если биржа добавит новый.
func parseListingLevel(s string) (domaincompany.ListingLevel, error) {
	level, ok := listingLevelByCode[s]
	if !ok {
		return domaincompany.ListingLevelUnspecified, fmt.Errorf("unexpected LISTLEVEL value: %q", s)
	}
	return level, nil
}
