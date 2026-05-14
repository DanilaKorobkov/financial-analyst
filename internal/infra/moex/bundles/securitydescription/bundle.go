// Package securitydescription — bundle справочной карточки бумаги
// поверх блока description эндпоинта MOEX ISS /iss/securities/{TICKER}.json.
// Делает один HTTP-вызов, парсит ответ и раскладывает значения по
// каноничным полям.
package securitydescription

import (
	"context"
	"fmt"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/client"
)

// ID — стабильный идентификатор bundle в реестре.
const ID = "security-description"

// fields — полный набор полей, которые bundle раскладывает в
// FieldValues после HTTP-ответа MOEX. Расширение списка требует
// синхронной правки parseDescription.
var fields = []data.FieldDescriptor{
	{ID: domaincompany.FieldTicker, Type: data.TypeString, Description: "Биржевой код бумаги."},
	{ID: domaincompany.FieldISIN, Type: data.TypeString, Description: "Международный идентификатор ценной бумаги."},
	{ID: domaincompany.FieldName, Type: data.TypeString, Description: "Полное наименование выпуска."},
	{ID: domaincompany.FieldShortName, Type: data.TypeString, Description: "Краткое наименование выпуска."},
	{ID: domaincompany.FieldIssueName, Type: data.TypeString, Description: "Наименование выпуска."},
	{ID: domaincompany.FieldLatName, Type: data.TypeString, Description: "Латинское наименование выпуска."},
	{ID: domaincompany.FieldRegNumber, Type: data.TypeString, Description: "Номер государственной регистрации выпуска."},
	{ID: domaincompany.FieldSecurityTypeName, Type: data.TypeString, Description: "Человекочитаемое название типа бумаги."},
	{ID: domaincompany.FieldSecurityGroup, Type: data.TypeString, Description: "Код группы инструмента."},
	{ID: domaincompany.FieldSecurityGroupName, Type: data.TypeString, Description: "Название группы инструмента."},
	{ID: domaincompany.FieldSecurityType, Type: data.TypeSecurityType, Description: "Тип бумаги."},
	{ID: domaincompany.FieldListingLevel, Type: data.TypeListingLevel, Description: "Котировальный уровень бумаги."},
	{ID: domaincompany.FieldFaceValue, Type: data.TypeString, Description: "Номинальная стоимость бумаги."},
	{ID: domaincompany.FieldFaceUnit, Type: data.TypeCurrency, Description: "Валюта номинала бумаги."},
	{ID: domaincompany.FieldIssueSize, Type: data.TypeInt64, Description: "Объём выпуска в штуках."},
	{ID: domaincompany.FieldIssueDate, Type: data.TypeDate, Description: "Дата начала торгов бумагой."},
	{ID: domaincompany.FieldRegistryDate, Type: data.TypeDate, Description: "Дата государственной регистрации выпуска."},
	{ID: domaincompany.FieldEmitterID, Type: data.TypeString, Description: "Идентификатор эмитента в справочнике источника."},
	{ID: domaincompany.FieldHasProspectus, Type: data.TypeBool, Description: "Наличие зарегистрированного проспекта."},
	{ID: domaincompany.FieldHasDefault, Type: data.TypeBool, Description: "Был ли дефолт по выпуску."},
	{ID: domaincompany.FieldHasTechnicalDefault, Type: data.TypeBool, Description: "Был ли технический дефолт по выпуску."},
	{ID: domaincompany.FieldEmitentMismatchCurrent, Type: data.TypeBool, Description: "Эмитент не соответствует требованию на текущий список."},
	{ID: domaincompany.FieldIsQualifiedInvestors, Type: data.TypeBool, Description: "Доступ только для квалифицированных инвесторов."},
	{ID: domaincompany.FieldMorningSession, Type: data.TypeBool, Description: "Допуск к утренней доп. сессии."},
	{ID: domaincompany.FieldEveningSession, Type: data.TypeBool, Description: "Допуск к вечерней доп. сессии."},
	{ID: domaincompany.FieldWeekendSession, Type: data.TypeBool, Description: "Допуск к доп. сессии выходного дня."},
}

// Bundle — реализация data.Bundle для блока description MOEX ISS.
type Bundle struct {
	client *client.Client
}

// New собирает bundle поверх общего MOEX-клиента.
func New(c *client.Client) *Bundle {
	return &Bundle{client: c}
}

// BundleID — реализация data.Bundle.
func (*Bundle) BundleID() string { return ID }

// Fields — реализация data.Bundle.
func (*Bundle) Fields() []data.FieldDescriptor { return fields }

// Fetch запрашивает description у MOEX и возвращает плоский FieldValues.
// Возвращает domaincompany.ErrNotFound, если ISS вернула пустой блок
// description (тикер не существует).
func (b *Bundle) Fetch(ctx context.Context, ticker string) (data.FieldValues, error) {
	resp, err := b.client.R().
		SetContext(ctx).
		SetPathParam("ticker", ticker).
		Get("/securities/{ticker}.json")
	if err != nil {
		if resp == nil || resp.StatusCode() == 0 {
			return nil, fmt.Errorf("moex request: %w", err)
		}
		return nil, err //nolint:wrapcheck // err уже сформирован OnAfterResponse в client.New ("moex http status N")
	}

	parsed, err := parseDescription(resp.Body())
	if err != nil {
		return nil, fmt.Errorf("parse description: %w", err)
	}

	return mapDescription(parsed)
}
