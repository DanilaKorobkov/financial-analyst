// Package stockinfo — bundle карточки эмитента поверх блока info
// эндпоинта FinanceMarker /api/fm/v2/stocks/{exchange}:{code}. Делает
// один HTTP-вызов, парсит ответ и раскладывает значения по каноничным
// полям.
//
// Поддерживается только MOEX: единственная биржа, по которой проект
// возвращает карточки.
package stockinfo

import (
	"context"
	"errors"
	"fmt"
	"time"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/cache/httpcache"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/client"
)

const (
	// ID — стабильный идентификатор bundle в реестре.
	ID = "stock-info"

	// includeInfo — значение query-параметра include, ограничивающее ответ
	// блоком info (классификация, описание, ссылки).
	includeInfo = "info"

	// codeExchangeMOEX — строковый код Московской биржи в формате FinanceMarker.
	codeExchangeMOEX = "MOEX"

	// cacheTTL — срок жизни HTTP-ответа /stocks в кеше клиента. Блок info
	// (классификация, описание эмитента, ссылки) меняется крайне редко,
	// поэтому раз в месяц достаточно. Bundle декларирует свой TTL у
	// каждого исходящего запроса через ctx — фактическое хранение
	// принадлежит httpcache-уровню клиента (см. httpcache.WithTTL).
	cacheTTL = 30 * 24 * time.Hour
)

// fields — полный набор полей, которые bundle раскладывает в
// FieldValues после HTTP-ответа FinanceMarker. Расширение списка требует
// синхронной правки translateStockInfo.
var fields = []data.FieldDescriptor{
	{ID: domaincompany.FieldIssuerName, Type: data.TypeString, Description: "Название эмитента у справочника."},
	{ID: domaincompany.FieldSector, Type: data.TypeString, Description: "Сектор эмитента (GICS)."},
	{ID: domaincompany.FieldSectorID, Type: data.TypeInt64, Description: "Числовой код сектора GICS."},
	{ID: domaincompany.FieldIndustryGroup, Type: data.TypeString, Description: "Группа отраслей GICS."},
	{ID: domaincompany.FieldIndustryGroupID, Type: data.TypeInt64, Description: "Числовой код группы отраслей GICS."},
	{ID: domaincompany.FieldIndustry, Type: data.TypeString, Description: "Отрасль эмитента (GICS)."},
	{ID: domaincompany.FieldIndustryID, Type: data.TypeInt64, Description: "Числовой код отрасли GICS."},
	{ID: domaincompany.FieldSubIndustry, Type: data.TypeString, Description: "Под-отрасль эмитента (GICS)."},
	{ID: domaincompany.FieldSubIndustryID, Type: data.TypeInt64, Description: "Числовой код под-отрасли GICS."},
	{ID: domaincompany.FieldCountry, Type: data.TypeString, Description: "Страна регистрации эмитента."},
	{ID: domaincompany.FieldDescription, Type: data.TypeString, Description: "Текстовое описание эмитента."},
	{ID: domaincompany.FieldSite, Type: data.TypeString, Description: "Корпоративный сайт эмитента."},
	{ID: domaincompany.FieldDisclosureLink, Type: data.TypeString, Description: "Ссылка на страницу раскрытия эмитента."},
	{ID: domaincompany.FieldPrimaryReportTicker, Type: data.TypeString, Description: "Тикер основной бумаги, к которой привязана отчётность."},
	{ID: domaincompany.FieldPrimaryReportExchange, Type: data.TypeExchange, Description: "Биржа основной отчётной бумаги."},
	{ID: domaincompany.FieldExchange, Type: data.TypeExchange, Description: "Биржа листинга бумаги."},
	{ID: domaincompany.FieldCurrency, Type: data.TypeCurrency, Description: "Валюта торгов бумагой."},
	{ID: domaincompany.FieldReportFrequency, Type: data.TypeReportFrequency, Description: "Частота публикации отчётности эмитентом."},
	{ID: domaincompany.FieldSPB, Type: data.TypeBool, Description: "Дополнительный листинг на СПБ-бирже."},
}

// Bundle — реализация data.Bundle для блока info FinanceMarker.
// Кеширование живёт уровнем ниже, на http-transport клиента; bundle
// только декларирует TTL у своего запроса через ctx (httpcache.WithTTL).
type Bundle struct {
	client *client.Client
}

// New собирает bundle поверх общего FinanceMarker-клиента.
func New(c *client.Client) *Bundle {
	return &Bundle{client: c}
}

// BundleID — реализация data.Bundle.
func (*Bundle) BundleID() string { return ID }

// Fields — реализация data.Bundle.
func (*Bundle) Fields() []data.FieldDescriptor { return fields }

// Fetch запрашивает карточку эмитента, переводит блок info в плоский
// FieldValues. Сетевые и HTTP-ошибки приходят из общего клиента уже
// классифицированными (см. client.New / classifyError), 404 здесь
// переводится в domaincompany.ErrNotFound.
//
// TTL для http-кеша выставляется на ctx через httpcache.WithTTL. Если
// у клиента кеш не подключён (CacheDir пустой) — аннотация остаётся в
// ctx без эффекта, и запрос идёт в сеть как есть.
func (b *Bundle) Fetch(ctx context.Context, ticker string) (data.FieldValues, error) {
	ctx = httpcache.WithTTL(ctx, cacheTTL)
	symbol := codeExchangeMOEX + ":" + ticker

	var dto stockDTO
	resp, err := b.client.R().
		SetContext(ctx).
		SetPathParam("symbol", symbol).
		SetQueryParam("include", includeInfo).
		SetResult(&dto).
		Get("/stocks/{symbol}")
	if err != nil {
		switch {
		case resp == nil || resp.StatusCode() == 0:
			return nil, fmt.Errorf("financemarker request: %w", err)
		case !resp.IsError():
			return nil, fmt.Errorf("decode financemarker payload: %w", err)
		case errors.Is(err, client.ErrNotFound):
			return nil, domaincompany.ErrNotFound
		default:
			return nil, err //nolint:wrapcheck // err уже сформирован classifyError общего клиента
		}
	}

	return translateStockInfo(&dto.Info), nil
}
