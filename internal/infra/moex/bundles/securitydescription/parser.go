package securitydescription

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
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

	// faceUnitByCode переводит код FACEUNIT MOEX в Currency-enum.
	// MOEX использует исторический "SUR" для рубля; ISO-коды USD/EUR — как есть.
	// Неизвестные значения возвращают zero (CurrencyUnspecified) — список валют
	// со временем расширяется.
	faceUnitByCode = map[string]domaincompany.Currency{
		"SUR": domaincompany.CurrencyRUB,
		"RUB": domaincompany.CurrencyRUB,
		"USD": domaincompany.CurrencyUSD,
		"EUR": domaincompany.CurrencyEUR,
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

// mapDescription раскладывает map[name]value блока description по
// каноничным полям FieldValues. Пустой map → domaincompany.ErrNotFound
// (тикер не найден).
func mapDescription(fields map[string]string) (data.FieldValues, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty description block: %w", domaincompany.ErrNotFound)
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
		return nil, p.err
	}

	return data.FieldValues{
		domaincompany.FieldTicker:                 fields["SECID"],
		domaincompany.FieldISIN:                   fields["ISIN"],
		domaincompany.FieldName:                   fields["NAME"],
		domaincompany.FieldShortName:              fields["SHORTNAME"],
		domaincompany.FieldIssueName:              fields["ISSUENAME"],
		domaincompany.FieldLatName:                fields["LATNAME"],
		domaincompany.FieldRegNumber:              fields["REGNUMBER"],
		domaincompany.FieldSecurityTypeName:       fields["TYPENAME"],
		domaincompany.FieldSecurityGroup:          fields["GROUP"],
		domaincompany.FieldSecurityGroupName:      fields["GROUPNAME"],
		domaincompany.FieldFaceValue:              fields["FACEVALUE"],
		domaincompany.FieldEmitterID:              fields["EMITTER_ID"],
		domaincompany.FieldIssueDate:              issueDate,
		domaincompany.FieldRegistryDate:           registryDate,
		domaincompany.FieldIssueSize:              issueSize,
		domaincompany.FieldSecurityType:           securityTypeByCode[fields["TYPE"]],
		domaincompany.FieldListingLevel:           listLevel,
		domaincompany.FieldFaceUnit:               faceUnitByCode[fields["FACEUNIT"]],
		domaincompany.FieldHasProspectus:          hasProspectus,
		domaincompany.FieldHasDefault:             hasDefault,
		domaincompany.FieldHasTechnicalDefault:    hasTechnicalDefault,
		domaincompany.FieldEmitentMismatchCurrent: emitentMismatch,
		domaincompany.FieldIsQualifiedInvestors:   isQualified,
		domaincompany.FieldMorningSession:         morningSession,
		domaincompany.FieldEveningSession:         eveningSession,
		domaincompany.FieldWeekendSession:         weekendSession,
	}, nil
}

func (p *descriptionParser) listingLevel(key string) domaincompany.ListingLevel {
	if p.err != nil {
		return domaincompany.ListingLevelUnspecified
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
func parseListingLevel(s string) (domaincompany.ListingLevel, error) {
	level, ok := listingLevelByCode[s]
	if !ok {
		return domaincompany.ListingLevelUnspecified, fmt.Errorf("unexpected LISTLEVEL value: %q", s)
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
