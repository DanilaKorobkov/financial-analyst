// Package app — экспортируемая точка сборки приложения financial-analyst.
//
// Единственный пакет, который импортирует domain / infra / presentation
// одновременно. Используется из cmd/server и tests/integration.
package app

import (
	"fmt"
	"net/http"

	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/companyprofile"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
)

// New собирает все слои приложения и возвращает готовый http.Handler.
//
// Поток сборки:
//  1. Поднимаем Provider-ы внешних источников (каждый со своим клиентом
//     и при необходимости — со своим кешем).
//  2. Собираем реестр bundles, давая каждому Provider зарегистрировать
//     свои bundles.
//  3. Поднимаем репозиторий профилей карточки эмитента (сейчас статический).
//  4. Создаём CompanyService поверх репозитория профилей и реестра.
func New(cfg *Config) (http.Handler, error) {
	providers := []data.Provider{
		moex.NewProvider(moex.ConfigProvider{
			BaseURL: cfg.Moex.BaseURL,
			Timeout: cfg.Moex.Timeout,
		}),
		financemarker.NewProvider(financemarker.ConfigProvider{
			BaseURL:   cfg.FinanceMarker.BaseURL,
			Token:     cfg.FinanceMarker.Token,
			Timeout:   cfg.FinanceMarker.Timeout,
			CacheRoot: cfg.FinanceMarker.CacheRootDir,
		}),
	}

	registry, err := buildCompanyRegistry(providers)
	if err != nil {
		return nil, fmt.Errorf("build company registry: %w", err)
	}

	profiles := companyprofile.NewDefaultStatic()
	companyService := services.NewCompanyService(services.ConfigCompanyService{
		Profiles: profiles,
		Fetcher:  registry,
	})
	srv := pconnect.NewServer(pconnect.ConfigServer{
		Companies: companyService,
		Registry:  registry,
	})

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	return mux, nil
}
