package connect

import (
	"context"
	"errors"

	connectrpc "connectrpc.com/connect"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
)

// fieldRegistry — узкий контракт реестра, нужный только для поиска типа
// поля при упаковке ответа. Сужение до одного метода вместо *data.Registry
// держит presentation слабо связанным с реестром и упрощает тесты.
type fieldRegistry interface {
	FieldByID(id string) (data.FieldDescriptor, bool)
}

// ConfigServer — параметры Connect-сервера.
type ConfigServer struct {
	// Companies — domain-сервис, отдающий значения полей по тикеру.
	Companies *services.CompanyService

	// Registry — источник типов полей для упаковки proto-ответа.
	Registry fieldRegistry
}

// Server реализует companyv1connect.CompanyServiceHandler поверх domain-сервиса.
type Server struct {
	companies *services.CompanyService
	registry  fieldRegistry
}

// NewServer собирает Connect-сервер.
func NewServer(cfg ConfigServer) *Server {
	return &Server{companies: cfg.Companies, registry: cfg.Registry}
}

// GetCompany — unary-метод CompanyService.GetCompany.
func (s *Server) GetCompany(
	ctx context.Context,
	req *connectrpc.Request[companyv1.GetCompanyRequest],
) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	values, err := s.companies.GetCompany(ctx, req.Msg.GetTicker())
	if err != nil {
		return nil, mapDomainError(err)
	}
	fields, err := toProtoFields(values, func(id string) (data.FieldType, bool) {
		fd, ok := s.registry.FieldByID(id)
		return fd.Type, ok
	})
	if err != nil {
		return nil, connectrpc.NewError(connectrpc.CodeInternal, err)
	}
	return connectrpc.NewResponse(&companyv1.GetCompanyResponse{Fields: fields}), nil
}

// mapDomainError переводит domain-ошибки в Connect-коды.
func mapDomainError(err error) error {
	switch {
	case errors.Is(err, services.ErrTickerEmpty):
		return connectrpc.NewError(connectrpc.CodeInvalidArgument, err)
	case errors.Is(err, company.ErrNotFound),
		errors.Is(err, company.ErrProfileNotFound):
		return connectrpc.NewError(connectrpc.CodeNotFound, err)
	case errors.Is(err, data.ErrFieldNotFound):
		return connectrpc.NewError(connectrpc.CodeFailedPrecondition, err)
	default:
		return connectrpc.NewError(connectrpc.CodeInternal, err)
	}
}
