// Package moex — реализация entities.CompanyRepository поверх MOEX ISS REST API.
package moex

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// CompanyRepository ходит в /iss/securities/{TICKER}.json (блок description).
type CompanyRepository struct {
	client *resty.Client
}

// NewCompanyRepository собирает репозиторий с resty-клиентом, привязанным к
// baseURL (корень MOEX ISS без завершающего слэша, например
// https://iss.moex.com/iss). timeout применяется к каждому запросу.
func NewCompanyRepository(baseURL string, timeout time.Duration) *CompanyRepository {
	// Параметры iss.* — это управляющие query-параметры самого MOEX ISS,
	// определены в их публичной справке (https://iss.moex.com/iss/reference/),
	// общие для всех эндпоинтов:
	//   - iss.json=extended  — JSON в виде массива записей вместо
	//     «columns + data»; удобнее разбирать (один объект на запись).
	//   - iss.meta=off       — не присылать блок описания колонок (типы,
	//     длины); экономит трафик и упрощает разбор.
	//   - iss.only=description — из всех блоков ответа (securities, boards,
	//     marketdata и т.п.) вернуть только блок description с реквизитами
	//     эмитента — это всё, что нужно CompanyRepository.
	client := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(timeout).
		SetQueryParams(map[string]string{
			"iss.json": "extended",
			"iss.meta": "off",
			"iss.only": "description",
		})
	return &CompanyRepository{client: client}
}

// FindByTicker запрашивает description у MOEX и маппит в entities.Company.
//
// Возвращает entities.ErrCompanyNotFound, если ISS вернула пустой блок
// description (тикер не существует).
func (r *CompanyRepository) FindByTicker(ctx context.Context, ticker string) (entities.Company, error) {
	resp, err := r.client.R().
		SetContext(ctx).
		SetPathParam("ticker", ticker).
		Get("/securities/{ticker}.json")
	if err != nil {
		return entities.Company{}, fmt.Errorf("moex request: %w", err)
	}
	if resp.IsError() {
		return entities.Company{}, fmt.Errorf("moex http status %d", resp.StatusCode())
	}

	fields, err := parseDescription(resp.Body())
	if err != nil {
		return entities.Company{}, fmt.Errorf("parse description: %w", err)
	}

	return mapDescription(fields)
}
