package connect

import (
	"context"
	"errors"

	connectrpc "connectrpc.com/connect"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/companycard"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
)

// Server реализует companyv1connect.CompanyServiceHandler поверх domain-сервиса.
type Server struct {
	cards *services.CompanyCardService
}

// NewServer собирает Connect-сервер вокруг доменного сервиса.
func NewServer(cards *services.CompanyCardService) *Server {
	return &Server{cards: cards}
}

// GetCompany — unary-метод CompanyService.GetCompany.
func (s *Server) GetCompany(
	ctx context.Context,
	req *connectrpc.Request[companyv1.GetCompanyRequest],
) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	card, err := s.cards.GetCard(ctx, req.Msg.GetTicker())
	if err != nil {
		return nil, mapDomainError(err)
	}
	return connectrpc.NewResponse(&companyv1.GetCompanyResponse{
		Company: toProtoCompany(&card),
	}), nil
}

// mapDomainError переводит domain-ошибки в Connect-коды.
func mapDomainError(err error) error {
	switch {
	case errors.Is(err, services.ErrTickerEmpty):
		return connectrpc.NewError(connectrpc.CodeInvalidArgument, err)
	case errors.Is(err, companycard.ErrNotFound):
		return connectrpc.NewError(connectrpc.CodeNotFound, err)
	default:
		return connectrpc.NewError(connectrpc.CodeInternal, err)
	}
}
