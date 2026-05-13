package app

import (
	"net/http"

	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	fccompany "github.com/DanilaKorobkov/financial-analyst/internal/infra/filecache/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker"
	fmcompany "github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
	moexcompany "github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/company"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
)

// New собирает все слои приложения и возвращает готовый http.Handler.
//
// Чистая композиция: общий клиент → gateway → (опциональный кеш) →
// CompanyService → Connect handler → http.ServeMux.
//
// Компания собирается из двух секций параллельно:
//   - идентификация — MOEX (без кеша, свободный источник);
//   - классификация — FinanceMarker (через файловый кеш с TTL, у источника
//     квота на запросы).
func New(cfg *Config) http.Handler {
	moexClient := moex.NewClient(moex.ConfigClient{
		BaseURL: cfg.Moex.BaseURL,
		Timeout: cfg.Moex.Timeout,
	})
	identities := moexcompany.NewIdentityGateway(moexClient)

	fmClient := financemarker.NewClient(financemarker.ConfigClient{
		BaseURL: cfg.FinanceMarker.BaseURL,
		Token:   cfg.FinanceMarker.Token,
		Timeout: cfg.FinanceMarker.Timeout,
	})
	classifications := fccompany.NewClassificationProxy(fccompany.ConfigClassificationProxy{
		Delegate: fmcompany.NewClassificationGateway(fmClient),
		Dir:      cfg.ClassCache.Dir,
		TTL:      cfg.ClassCache.TTL,
	})

	companyService := services.NewCompanyService(identities, classifications)
	srv := pconnect.NewServer(companyService)

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	return mux
}
