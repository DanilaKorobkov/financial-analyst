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
	companyInfo *services.CompanyInfo
	companyCard *services.CompanyCard
}

// NewServer собирает Connect-сервер вокруг доменных сервисов.
func NewServer(
	companyInfo *services.CompanyInfo,
	companyCard *services.CompanyCard,
) *Server {
	return &Server{
		companyInfo: companyInfo,
		companyCard: companyCard,
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

// GetCompanyCard — unary-метод CompanyService.GetCompanyCard.
func (s *Server) GetCompanyCard(
	ctx context.Context,
	req *connectrpc.Request[companyv1.GetCompanyCardRequest],
) (*connectrpc.Response[companyv1.GetCompanyCardResponse], error) {
	card, err := s.companyCard.FindByTicker(ctx, req.Msg.GetTicker())
	if err != nil {
		return nil, mapDomainError(err)
	}
	return connectrpc.NewResponse(&companyv1.GetCompanyCardResponse{
		Card: toProtoCompanyCard(&card),
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
