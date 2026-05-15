package securitydescription

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
)

// dateLayout — MOEX отдаёт даты блока description в ISO-форме без
// времени. Парсим с фиксированным UTC, чтобы дата не уезжала на день
// в зависимости от часового пояса исполнителя.
const dateLayout = "2006-01-02"

var (
	// jsonParser — JSON-парсер с поведением, идентичным encoding/json:
	// sorted map keys, html-escape, validateRawMessage. См. rules/golang.md.
	jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary

	// errDescriptionBlockMissing — payload не содержит ожидаемого блока description.
	errDescriptionBlockMissing = fmt.Errorf("description block missing in payload")

	// securityTypeByCode переводит код TYPE блока description MOEX в domain-enum.
	// Неизвестные значения возвращают zero (SecurityTypeUnspecified) — TYPE
	// задаётся биржей и со временем расширяется (фонды, облигации и т.п.).
	securityTypeByCode = map[string]company.SecurityType{
		"common_share":       company.SecurityTypeCommonShare,
		"preferred_share":    company.SecurityTypePreferredShare,
		"depositary_receipt": company.SecurityTypeDepositaryReceipt,
	}

	// listingLevelByCode переводит строковое значение LISTLEVEL ("1" / "2" / "3")
	// в domain-enum. Пустое значение — биржа не указала уровень (Unspecified).
	// Любое другое значение в parseListingLevel — ошибка: список уровней
	// зафиксирован MOEX.
	listingLevelByCode = map[string]company.ListingLevel{
		"":  company.ListingLevelUnspecified,
		"1": company.ListingLevelFirst,
		"2": company.ListingLevelSecond,
		"3": company.ListingLevelThird,
	}

	// faceUnitByCode переводит код FACEUNIT MOEX в Currency-enum.
	// MOEX использует исторический "SUR" для рубля; ISO-коды USD/EUR — как есть.
	// Неизвестные значения возвращают zero (CurrencyUnspecified) — список валют
	// со временем расширяется.
	faceUnitByCode = map[string]company.Currency{
		"SUR": company.CurrencyRUB,
		"RUB": company.CurrencyRUB,
		"USD": company.CurrencyUSD,
		"EUR": company.CurrencyEUR,
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

// descriptionParser — накопитель ошибок при разборе типизированных
// полей блока description. После первой ошибки следующие парсы
// возвращают zero и err не перезаписывается. Цель — держать
// mapDescription плоским по cyclomatic-complexity: вместо двенадцати
// `if err != nil` идёт одна проверка p.err в конце.
type descriptionParser struct {
	fields map[string]string
	err    error
}

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

// mapDescription раскладывает map[name]value блока description по полям
// SecurityDescription. Пустой map → company.ErrNotFound (тикер не найден).
func mapDescription(fields map[string]string) (company.SecurityDescription, error) {
	if len(fields) == 0 {
		return company.SecurityDescription{}, fmt.Errorf("empty description block: %w", company.ErrNotFound)
	}

	p := descriptionParser{fields: fields}
	listLevel := p.listingLevel("LISTLEVEL")
	issueSize := p.int64("ISSUESIZE")
	issueDate := p.date("ISSUEDATE")
	registryDate := p.date("REGISTRY_DATE")
	hasProspectus := p.bool("HASPROSPECTUS")
	hasDefault := p.bool("HASDEFAULT")
	hasTechnicalDefault := p.bool("HASTECHNICALDEFAULT")
	emitentMismatch := p.bool("EMITENTMISMATCHCUR")
	isQualified := p.bool("ISQUALIFIEDINVESTORS")
	morningSession := p.bool("MORNINGSESSION")
	eveningSession := p.bool("EVENINGSESSION")
	weekendSession := p.bool("WEEKENDSESSION")
	if p.err != nil {
		return company.SecurityDescription{}, p.err
	}

	return company.SecurityDescription{
		Ticker:                 fields["SECID"],
		ISIN:                   fields["ISIN"],
		Name:                   fields["NAME"],
		ShortName:              fields["SHORTNAME"],
		IssueName:              fields["ISSUENAME"],
		LatName:                fields["LATNAME"],
		RegNumber:              fields["REGNUMBER"],
		SecurityTypeName:       fields["TYPENAME"],
		SecurityGroup:          fields["GROUP"],
		SecurityGroupName:      fields["GROUPNAME"],
		FaceValue:              fields["FACEVALUE"],
		EmitterID:              fields["EMITTER_ID"],
		SecurityType:           securityTypeByCode[fields["TYPE"]],
		ListingLevel:           listLevel,
		FaceUnit:               faceUnitByCode[fields["FACEUNIT"]],
		IssueSize:              issueSize,
		IssueDate:              issueDate,
		RegistryDate:           registryDate,
		HasProspectus:          hasProspectus,
		HasDefault:             hasDefault,
		HasTechnicalDefault:    hasTechnicalDefault,
		EmitentMismatchCurrent: emitentMismatch,
		IsQualifiedInvestors:   isQualified,
		MorningSession:         morningSession,
		EveningSession:         eveningSession,
		WeekendSession:         weekendSession,
	}, nil
}

func (p *descriptionParser) listingLevel(key string) company.ListingLevel {
	if p.err != nil {
		return company.ListingLevelUnspecified
	}
	v, err := parseListingLevel(p.fields[key])
	if err != nil {
		p.err = fmt.Errorf("%s: %w", key, err)
	}
	return v
}

func (p *descriptionParser) int64(key string) int64 {
	if p.err != nil {
		return 0
	}
	v, err := parseOptionalInt64(p.fields[key])
	if err != nil {
		p.err = fmt.Errorf("%s: %w", key, err)
	}
	return v
}

func (p *descriptionParser) date(key string) time.Time {
	if p.err != nil {
		return time.Time{}
	}
	v, err := parseOptionalDate(p.fields[key])
	if err != nil {
		p.err = fmt.Errorf("%s: %w", key, err)
	}
	return v
}

func (p *descriptionParser) bool(key string) bool {
	if p.err != nil {
		return false
	}
	v, err := parseOptionalBool(p.fields[key])
	if err != nil {
		p.err = fmt.Errorf("%s: %w", key, err)
	}
	return v
}

// parseListingLevel — единственная функция, где нужна валидация (неизвестный
// уровень — ошибка, а не Unspecified): список уровней зафиксирован MOEX и
// должен быть пополнен сознательно, если биржа добавит новый.
func parseListingLevel(s string) (company.ListingLevel, error) {
	level, ok := listingLevelByCode[s]
	if !ok {
		return company.ListingLevelUnspecified, fmt.Errorf("unexpected LISTLEVEL value: %q", s)
	}
	return level, nil
}

// parseOptionalInt64 — пустая строка переводится в 0 (поле отсутствует),
// непустая обязана разобраться в int64. Так отличаем «не отдал» от
// «отдал явный 0».
func parseOptionalInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse int64 %q: %w", s, err)
	}
	return v, nil
}

// parseOptionalBool — MOEX отдаёт булевы поля как "0" / "1" / "" (нет).
func parseOptionalBool(s string) (bool, error) {
	switch s {
	case "", "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, fmt.Errorf("unexpected bool value: %q", s)
	}
}

// parseOptionalDate — MOEX отдаёт даты в формате YYYY-MM-DD; пустая
// строка — поле отсутствует, возвращаем zero time.Time.
func parseOptionalDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	t, err := time.ParseInLocation(dateLayout, s, time.UTC)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse date %q: %w", s, err)
	}
	return t, nil
}
