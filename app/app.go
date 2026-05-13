package app

import (
	"net/http"

	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/filecache/companycard"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker"
	fmcard "github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/companycard"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
	moexcard "github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/companycard"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
)

// New собирает все слои приложения и возвращает готовый http.Handler.
//
// Чистая композиция: общий клиент → gateway → (опциональный кеш) →
// CompanyCardService → Connect handler → http.ServeMux.
//
// Карточка собирается из двух секций параллельно:
//   - идентификация — MOEX (без кеша, свободный источник);
//   - классификация — FinanceMarker (через файловый кеш с TTL, у источника
//     квота на запросы).
func New(cfg *Config) http.Handler {
	moexClient := moex.NewClient(moex.ConfigClient{
		BaseURL: cfg.Moex.BaseURL,
		Timeout: cfg.Moex.Timeout,
	})
	identities := moexcard.NewIdentityGateway(moexClient)

	fmClient := financemarker.NewClient(financemarker.ConfigClient{
		BaseURL: cfg.FinanceMarker.BaseURL,
		Token:   cfg.FinanceMarker.Token,
		Timeout: cfg.FinanceMarker.Timeout,
	})
	classifications := companycard.NewClassificationProxy(companycard.ConfigClassificationProxy{
		Delegate: fmcard.NewClassificationGateway(fmClient),
		Dir:      cfg.ClassCache.Dir,
		TTL:      cfg.ClassCache.TTL,
	})

	cardService := services.NewCompanyCardService(identities, classifications)
	srv := pconnect.NewServer(cardService)

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	return mux
}
