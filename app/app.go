// Package app — экспортируемая точка сборки приложения financial-analyst.
//
// Единственный пакет, который импортирует domain / infra / presentation
// одновременно. Используется из cmd/server и tests/integration.
package app

import (
	"net/http"
	"path/filepath"

	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	infracompany "github.com/DanilaKorobkov/financial-analyst/internal/infra/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/bundles/stockinfo"
	fmclient "github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/client"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/bundles/securitydescription"
	moexclient "github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/client"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
)

// New собирает все слои приложения и возвращает готовый http.Handler.
//
// Поток сборки:
//  1. Поднимаем HTTP-клиенты внешних источников (каждый со своим
//     таймаутом и при необходимости — со своим кешем).
//  2. Поднимаем источники секций (SecurityDescriptionSource, StockInfoSource).
//  3. Оборачиваем источники в company.Repository — он собирает агрегат
//     параллельно из своих секций.
//  4. Создаём CompanyService поверх репозитория.
//  5. Поднимаем Connect-сервер.
func New(cfg *Config) (http.Handler, error) {
	moexHTTP := moexclient.New(moexclient.Config{
		BaseURL: cfg.Moex.BaseURL,
		Timeout: cfg.Moex.Timeout,
	})

	fmHTTP := fmclient.New(fmclient.Config{
		BaseURL:  cfg.FinanceMarker.BaseURL,
		Token:    cfg.FinanceMarker.Token,
		Timeout:  cfg.FinanceMarker.Timeout,
		CacheDir: filepath.Join(cfg.FinanceMarker.CacheRootDir, "financemarker"),
	})

	companies := infracompany.NewRepository(infracompany.ConfigRepository{
		SecurityDescription: securitydescription.New(moexHTTP),
		StockInfo:           stockinfo.New(fmHTTP),
	})
	companyService := services.NewCompanyService(services.ConfigCompanyService{
		Companies: companies,
	})
	srv := pconnect.NewServer(pconnect.ConfigServer{Companies: companyService})

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	return mux, nil
}
