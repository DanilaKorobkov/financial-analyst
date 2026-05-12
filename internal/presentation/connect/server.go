package connect

import (
	"context"
	"errors"

	connectrpc "connectrpc.com/connect"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
)

// Server реализует companyv1connect.CompanyServiceHandler поверх domain-сервисов.
type Server struct {
	companyInfo    *services.CompanyInfo
	companyMetrics *services.CompanyMetrics
}

// NewServer собирает Connect-сервер вокруг доменных сервисов.
func NewServer(
	companyInfo *services.CompanyInfo,
	companyMetrics *services.CompanyMetrics,
) *Server {
	return &Server{
		companyInfo:    companyInfo,
		companyMetrics: companyMetrics,
	}
}

// GetCompany — unary-метод CompanyService.GetCompany.
func (s *Server) GetCompany(
	ctx context.Context,
	req *connectrpc.Request[companyv1.GetCompanyRequest],
) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	company, err := s.companyInfo.Lookup(ctx, req.Msg.GetTicker())
	if err != nil {
		return nil, mapDomainError(err)
	}
	return connectrpc.NewResponse(&companyv1.GetCompanyResponse{
		Company: toProtoCompany(&company),
	}), nil
}

// GetCompanyMetrics — unary-метод CompanyService.GetCompanyMetrics.
func (s *Server) GetCompanyMetrics(
	ctx context.Context,
	req *connectrpc.Request[companyv1.GetCompanyMetricsRequest],
) (*connectrpc.Response[companyv1.GetCompanyMetricsResponse], error) {
	metrics, err := s.companyMetrics.FindByTicker(ctx, req.Msg.GetTicker())
	if err != nil {
		return nil, mapDomainError(err)
	}
	return connectrpc.NewResponse(&companyv1.GetCompanyMetricsResponse{
		Metrics: toProtoCompanyMetrics(&metrics),
	}), nil
}

// mapDomainError переводит domain-ошибки в Connect-коды.
func mapDomainError(err error) error {
	switch {
	case errors.Is(err, services.ErrTickerEmpty):
		return connectrpc.NewError(connectrpc.CodeInvalidArgument, err)
	case errors.Is(err, entities.ErrCompanyNotFound),
		errors.Is(err, entities.ErrNotFound):
		return connectrpc.NewError(connectrpc.CodeNotFound, err)
	case errors.Is(err, entities.ErrUnauthorized):
		return connectrpc.NewError(connectrpc.CodeUnauthenticated, err)
	case errors.Is(err, entities.ErrQuotaExceeded):
		return connectrpc.NewError(connectrpc.CodeResourceExhausted, err)
	default:
		return connectrpc.NewError(connectrpc.CodeInternal, err)
	}
}
