package app

import (
	"net/http"

	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
)

// New собирает все слои приложения и возвращает готовый http.Handler.
//
// Чистая композиция: HTTP-клиенты внешних источников → репозитории →
// domain-сервисы → Connect handler → http.ServeMux.
func New(cfg *Config) http.Handler {
	companies := moex.NewCompanyRepository(cfg.Moex.BaseURL, cfg.Moex.Timeout)
	companyInfo := services.NewCompanyInfo(companies)

	fmClient := financemarker.NewClient(
		cfg.FinanceMarker.BaseURL,
		cfg.FinanceMarker.Token,
		cfg.FinanceMarker.Timeout,
	)
	companyMetricsRepo := financemarker.NewCompanyMetricsRepository(fmClient)
	companyMetrics := services.NewCompanyMetrics(companyMetricsRepo)

	srv := pconnect.NewServer(companyInfo, companyMetrics)

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	return mux
}
