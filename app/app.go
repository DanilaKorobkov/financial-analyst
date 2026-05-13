package app

import (
	"net/http"

	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
	moexcompany "github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/company"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
)

// New собирает все слои приложения и возвращает готовый http.Handler.
//
// Чистая композиция: общий MOEX-клиент → companyRepository → CompanyInfo →
// Connect handler → http.ServeMux.
func New(cfg Config) http.Handler {
	moexClient := moex.NewClient(moex.ConfigClient{
		BaseURL: cfg.Moex.BaseURL,
		Timeout: cfg.Moex.Timeout,
	})
	companies := moexcompany.NewRepository(moexClient)
	companyInfo := services.NewCompanyInfo(companies)
	srv := pconnect.NewServer(companyInfo)

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	return mux
}
