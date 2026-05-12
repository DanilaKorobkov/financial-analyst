package app

import (
	"net/http"

	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
)

// New собирает все слои приложения и возвращает готовый http.Handler.
//
// Чистая композиция: resty client → MoexRepository → CompanyInfo →
// Connect handler → http.ServeMux.
func New(cfg Config) http.Handler {
	companies := moex.NewCompanyRepository(cfg.Moex.BaseURL, cfg.Moex.Timeout)
	companyInfo := services.NewCompanyInfo(companies)
	srv := pconnect.NewServer(companyInfo)

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	return mux
}
