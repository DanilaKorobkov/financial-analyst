package connect

import (
	"context"
	"errors"

	connectrpc "connectrpc.com/connect"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
)

// Server реализует companyv1connect.CompanyServiceHandler поверх domain-сервиса.
type Server struct {
	companies *services.CompanyService
}

// NewServer собирает Connect-сервер вокруг доменного сервиса.
func NewServer(companies *services.CompanyService) *Server {
	return &Server{companies: companies}
}

// GetCompany — unary-метод CompanyService.GetCompany.
func (s *Server) GetCompany(
	ctx context.Context,
	req *connectrpc.Request[companyv1.GetCompanyRequest],
) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	found, err := s.companies.GetCompany(ctx, req.Msg.GetTicker())
	if err != nil {
		return nil, mapDomainError(err)
	}
	return connectrpc.NewResponse(&companyv1.GetCompanyResponse{
		Company: toProtoCompany(&found),
	}), nil
}

// mapDomainError переводит domain-ошибки в Connect-коды.
func mapDomainError(err error) error {
	switch {
	case errors.Is(err, services.ErrTickerEmpty):
		return connectrpc.NewError(connectrpc.CodeInvalidArgument, err)
	case errors.Is(err, company.ErrNotFound):
		return connectrpc.NewError(connectrpc.CodeNotFound, err)
	default:
		return connectrpc.NewError(connectrpc.CodeInternal, err)
	}
}
