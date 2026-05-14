package connect_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	connectrpc "connectrpc.com/connect"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/services"
	pconnect "github.com/DanilaKorobkov/financial-analyst/internal/presentation/connect"
	company_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/company"
	data_mock "github.com/DanilaKorobkov/financial-analyst/mocks/internal_/domain/data"
)

type serverSuite struct {
	suite.Suite

	profiles            *company_mock.ProfileRepository
	securityDescription *data_mock.Bundle
	stockInfo           *data_mock.Bundle
	server              *httptest.Server
	client              companyv1connect.CompanyServiceClient
	fieldIDs            []data.Field
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(serverSuite))
}

func (s *serverSuite) SetupTest() {
	s.profiles = company_mock.NewProfileRepository(s.T())
	s.securityDescription = data_mock.NewBundle(s.T())
	s.stockInfo = data_mock.NewBundle(s.T())
	s.securityDescription.EXPECT().BundleID().Return("security-description").Maybe()
	s.securityDescription.EXPECT().
		Fields().
		Return([]data.FieldDescriptor{
			{ID: company.FieldTicker, Type: data.TypeString},
			{ID: company.FieldISIN, Type: data.TypeString},
			{ID: company.FieldName, Type: data.TypeString},
			{ID: company.FieldSecurityType, Type: data.TypeSecurityType},
			{ID: company.FieldListingLevel, Type: data.TypeListingLevel},
			{ID: company.FieldIssueDate, Type: data.TypeDate},
			{ID: company.FieldRegistryDate, Type: data.TypeDate},
			{ID: company.FieldIssueSize, Type: data.TypeInt64},
			{ID: company.FieldHasDefault, Type: data.TypeBool},
		}).
		Maybe()
	s.stockInfo.EXPECT().BundleID().Return("stock-info").Maybe()
	s.stockInfo.EXPECT().
		Fields().
		Return([]data.FieldDescriptor{
			{ID: company.FieldIssuerName, Type: data.TypeString},
			{ID: company.FieldSector, Type: data.TypeString},
			{ID: company.FieldIndustry, Type: data.TypeString},
			{ID: company.FieldCountry, Type: data.TypeString},
			{ID: company.FieldPrimaryReportTicker, Type: data.TypeString},
			{ID: company.FieldExchange, Type: data.TypeExchange},
			{ID: company.FieldCurrency, Type: data.TypeCurrency},
			{ID: company.FieldReportFrequency, Type: data.TypeReportFrequency},
		}).
		Maybe()

	s.fieldIDs = []data.Field{
		company.FieldTicker,
		company.FieldISIN,
		company.FieldName,
		company.FieldSecurityType,
		company.FieldListingLevel,
		company.FieldIssueDate,
		company.FieldRegistryDate,
		company.FieldIssueSize,
		company.FieldHasDefault,
		company.FieldIssuerName,
		company.FieldSector,
		company.FieldIndustry,
		company.FieldCountry,
		company.FieldPrimaryReportTicker,
		company.FieldExchange,
		company.FieldCurrency,
		company.FieldReportFrequency,
	}
	// Дефолтный профиль: одни и те же поля, любой непустой тикер.
	s.profiles.EXPECT().
		FindByTicker(mock.Anything, mock.Anything).
		Return(company.Profile{FieldIDs: s.fieldIDs}, nil).
		Maybe()

	registry := data.NewRegistry()
	s.Require().NoError(registry.Register("moex", s.securityDescription))
	s.Require().NoError(registry.Register("financemarker", s.stockInfo))

	srv := pconnect.NewServer(pconnect.ConfigServer{
		Companies: services.NewCompanyService(services.ConfigCompanyService{
			Profiles: s.profiles,
			Registry: registry,
		}),
		Registry: registry,
	})

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	s.server = httptest.NewServer(mux)
	s.client = companyv1connect.NewCompanyServiceClient(s.server.Client(), s.server.URL)
}

func (s *serverSuite) TearDownTest() {
	s.server.Close()
}

func (s *serverSuite) TestGetCompanyHappyPath() {
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "SBER").
		Return(data.FieldValues{
			company.FieldTicker:       "SBER",
			company.FieldISIN:         "RU0009029540",
			company.FieldName:         "Сбербанк России ПАО ао",
			company.FieldSecurityType: company.SecurityTypeCommonShare,
			company.FieldListingLevel: company.ListingLevelFirst,
		}, nil).
		Once()
	s.stockInfo.EXPECT().
		Fetch(mock.Anything, "SBER").
		Return(data.FieldValues{
			company.FieldIssuerName:          "Сбербанк",
			company.FieldSector:              "Финансы",
			company.FieldIndustry:            "Банковская деятельность",
			company.FieldCountry:             "Россия",
			company.FieldPrimaryReportTicker: "SBER",
			company.FieldExchange:            company.ExchangeMOEX,
			company.FieldCurrency:            company.CurrencyRUB,
		}, nil).
		Once()

	resp, err := s.call("SBER")
	s.Require().NoError(err)

	fields := resp.Msg.GetFields()
	s.Equal("SBER", fields[string(company.FieldTicker)].GetStringValue())
	s.Equal("RU0009029540", fields[string(company.FieldISIN)].GetStringValue())
	s.Equal("Сбербанк России ПАО ао", fields[string(company.FieldName)].GetStringValue())
	s.Equal("common_share", fields[string(company.FieldSecurityType)].GetStringValue())
	s.Equal("first", fields[string(company.FieldListingLevel)].GetStringValue())
	s.Equal("Сбербанк", fields[string(company.FieldIssuerName)].GetStringValue())
	s.Equal("moex", fields[string(company.FieldExchange)].GetStringValue())
	s.Equal("RUB", fields[string(company.FieldCurrency)].GetStringValue())
}

func (s *serverSuite) TestGetCompanyEnumCodes() {
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "SBER").
		Return(data.FieldValues{
			company.FieldSecurityType: company.SecurityTypePreferredShare,
			company.FieldListingLevel: company.ListingLevelSecond,
		}, nil).
		Once()
	s.stockInfo.EXPECT().
		Fetch(mock.Anything, "SBER").
		Return(data.FieldValues{
			company.FieldExchange:        company.ExchangeMOEX,
			company.FieldCurrency:        company.CurrencyUSD,
			company.FieldReportFrequency: company.ReportFrequencyQuarterly,
		}, nil).
		Once()

	resp, err := s.call("SBER")
	s.Require().NoError(err)

	fields := resp.Msg.GetFields()
	s.Equal("preferred_share", fields[string(company.FieldSecurityType)].GetStringValue())
	s.Equal("second", fields[string(company.FieldListingLevel)].GetStringValue())
	s.Equal("moex", fields[string(company.FieldExchange)].GetStringValue())
	s.Equal("USD", fields[string(company.FieldCurrency)].GetStringValue())
	s.Equal("quarterly", fields[string(company.FieldReportFrequency)].GetStringValue())
}

func (s *serverSuite) TestGetCompanyEnumUnspecifiedDroppedFromMap() {
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "any").
		Return(data.FieldValues{
			company.FieldSecurityType: company.SecurityTypeUnspecified,
		}, nil).
		Once()
	s.stockInfo.EXPECT().
		Fetch(mock.Anything, "any").
		Return(data.FieldValues{}, nil).
		Once()

	resp, err := s.call("any")
	s.Require().NoError(err)

	_, found := resp.Msg.GetFields()[string(company.FieldSecurityType)]
	s.False(found)
}

// TestGetCompanyZeroDateDroppedFromMap: нулевой time.Time у поля типа
// Date трактуется как «значения нет» — mapper не кладёт ключ в map,
// клиент видит отсутствие поля. Покрывает IsZero-ветку encodeDate.
func (s *serverSuite) TestGetCompanyZeroDateDroppedFromMap() {
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "any").
		Return(data.FieldValues{
			company.FieldIssueDate: time.Time{},
		}, nil).
		Once()
	s.stockInfo.EXPECT().
		Fetch(mock.Anything, "any").
		Return(data.FieldValues{}, nil).
		Once()

	resp, err := s.call("any")
	s.Require().NoError(err)

	_, found := resp.Msg.GetFields()[string(company.FieldIssueDate)]
	s.False(found)
}

func (s *serverSuite) TestGetCompanyDatesMappedToTimestamp() {
	issueDate := time.Date(2007, 7, 20, 0, 0, 0, 0, time.UTC)
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "SBER").
		Return(data.FieldValues{
			company.FieldIssueDate: issueDate,
		}, nil).
		Once()
	s.stockInfo.EXPECT().
		Fetch(mock.Anything, "SBER").
		Return(data.FieldValues{}, nil).
		Once()

	resp, err := s.call("SBER")
	s.Require().NoError(err)

	fields := resp.Msg.GetFields()
	s.Require().NotNil(fields[string(company.FieldIssueDate)].GetTimestampValue())
	s.Equal(issueDate.Unix(), fields[string(company.FieldIssueDate)].GetTimestampValue().GetSeconds())
	_, found := fields[string(company.FieldRegistryDate)]
	s.False(found)
}

func (s *serverSuite) TestGetCompanyIntAndBool() {
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "SBER").
		Return(data.FieldValues{
			company.FieldIssueSize:  int64(21586948000),
			company.FieldHasDefault: false,
		}, nil).
		Once()
	s.stockInfo.EXPECT().
		Fetch(mock.Anything, "SBER").
		Return(data.FieldValues{}, nil).
		Once()

	resp, err := s.call("SBER")
	s.Require().NoError(err)

	fields := resp.Msg.GetFields()
	s.Equal(int64(21586948000), fields[string(company.FieldIssueSize)].GetIntValue())
	s.False(fields[string(company.FieldHasDefault)].GetBoolValue())
}

func (s *serverSuite) TestGetCompanyProfileNotFoundIsCodeNotFound() {
	profiles := company_mock.NewProfileRepository(s.T())
	profiles.EXPECT().
		FindByTicker(mock.Anything, "missing").
		Return(company.Profile{}, company.ErrProfileNotFound).
		Once()
	s.replaceProfiles(profiles)

	_, err := s.call("missing")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeNotFound, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyUnknownFieldIsFailedPrecondition() {
	profiles := company_mock.NewProfileRepository(s.T())
	profiles.EXPECT().
		FindByTicker(mock.Anything, "SBER").
		Return(company.Profile{FieldIDs: []data.Field{"unknown"}}, nil).
		Once()
	s.replaceProfiles(profiles)

	_, err := s.call("SBER")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeFailedPrecondition, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyNotFoundFromSecurityDescription() {
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "missing").
		Return(nil, company.ErrNotFound).
		Once()
	s.stockInfo.EXPECT().
		Fetch(mock.Anything, "missing").
		Return(data.FieldValues{}, nil).
		Maybe()

	_, err := s.call("missing")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeNotFound, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyInvalidArgument() {
	_, err := s.call("")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInvalidArgument, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyEncodeErrorIsInternal() {
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "any").
		Return(data.FieldValues{company.FieldTicker: 42}, nil).
		Once()
	s.stockInfo.EXPECT().
		Fetch(mock.Anything, "any").
		Return(data.FieldValues{}, nil).
		Maybe()

	_, err := s.call("any")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInternal, connectErr.Code())
}

// TestGetCompanyEncodeMismatchAllBranches гоняет каждую ветку encode*
// в encodeFieldValue: bundle возвращает значение неверного Go-типа
// для каждого FieldType — mapper отдаёт ошибку, server заворачивает
// её в CodeInternal. Покрытие defensive type-cast'ов держится без
// доступа к приватному API.
func (s *serverSuite) TestGetCompanyEncodeMismatchAllBranches() {
	cases := []struct {
		raw       any
		name      string
		fieldID   data.Field
		fieldType data.FieldType
	}{
		{name: "string", fieldID: "synthetic-string", fieldType: data.TypeString, raw: 1},
		{name: "int64", fieldID: "synthetic-int64", fieldType: data.TypeInt64, raw: "not-int"},
		{name: "bool", fieldID: "synthetic-bool", fieldType: data.TypeBool, raw: "not-bool"},
		{name: "date", fieldID: "synthetic-date", fieldType: data.TypeDate, raw: "not-time"},
		{name: "security_type", fieldID: "synthetic-security-type", fieldType: data.TypeSecurityType, raw: "not-enum"},
		{name: "listing_level", fieldID: "synthetic-listing-level", fieldType: data.TypeListingLevel, raw: "not-enum"},
		{name: "currency", fieldID: "synthetic-currency", fieldType: data.TypeCurrency, raw: "not-enum"},
		{name: "exchange", fieldID: "synthetic-exchange", fieldType: data.TypeExchange, raw: "not-enum"},
		{name: "report_frequency", fieldID: "synthetic-report-frequency", fieldType: data.TypeReportFrequency, raw: "not-enum"},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			err := s.callWithSyntheticBundle(c.fieldID, c.fieldType, c.raw)
			var connectErr *connectrpc.Error
			s.Require().ErrorAs(err, &connectErr)
			s.Equal(connectrpc.CodeInternal, connectErr.Code())
		})
	}
}

// callWithSyntheticBundle строит изолированный server с единственным
// bundle, у которого ровно одно поле заданного типа. Возвращает
// ошибку вызова GetCompany — удобно для проверки defensive-веток
// mapper'а без накопления mockery-ожиданий между подкейсами.
func (s *serverSuite) callWithSyntheticBundle(fieldID data.Field, ft data.FieldType, raw any) error {
	s.T().Helper()

	profiles := company_mock.NewProfileRepository(s.T())
	profiles.EXPECT().
		FindByTicker(mock.Anything, "any").
		Return(company.Profile{FieldIDs: []data.Field{fieldID}}, nil).
		Once()

	bundle := data_mock.NewBundle(s.T())
	bundle.EXPECT().BundleID().Return("synthetic").Maybe()
	bundle.EXPECT().
		Fields().
		Return([]data.FieldDescriptor{{ID: fieldID, Type: ft}}).
		Maybe()
	bundle.EXPECT().
		Fetch(mock.Anything, "any").
		Return(data.FieldValues{fieldID: raw}, nil).
		Once()

	registry := data.NewRegistry()
	s.Require().NoError(registry.Register("synthetic", bundle))

	srv := pconnect.NewServer(pconnect.ConfigServer{
		Companies: services.NewCompanyService(services.ConfigCompanyService{
			Profiles: profiles,
			Registry: registry,
		}),
		Registry: registry,
	})

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := companyv1connect.NewCompanyServiceClient(server.Client(), server.URL)
	_, err := client.GetCompany(
		context.Background(),
		connectrpc.NewRequest(&companyv1.GetCompanyRequest{Ticker: "any"}),
	)
	return err
}

// TestGetCompanyUnknownFieldTypeIsInternal: bundle декларирует поле
// с FieldType вне известного набора (например, новая константа добавлена
// в data, но забыта в mapper). encodeFieldValue падает на default-ветке,
// server возвращает CodeInternal. Через публичный API закрывает defensive-
// ветку без приватных вспомогательных функций.
func (s *serverSuite) TestGetCompanyUnknownFieldTypeIsInternal() {
	err := s.callWithSyntheticBundle("synthetic-unknown-type", data.FieldType(9999), "value")
	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInternal, connectErr.Code())
}

func (s *serverSuite) TestGetCompanyInternal() {
	s.securityDescription.EXPECT().
		Fetch(mock.Anything, "any").
		Return(nil, errors.New("downstream boom")).
		Once()
	s.stockInfo.EXPECT().
		Fetch(mock.Anything, "any").
		Return(data.FieldValues{}, nil).
		Maybe()

	_, err := s.call("any")

	var connectErr *connectrpc.Error
	s.Require().ErrorAs(err, &connectErr)
	s.Equal(connectrpc.CodeInternal, connectErr.Code())
}

func (s *serverSuite) call(ticker string) (*connectrpc.Response[companyv1.GetCompanyResponse], error) {
	return s.client.GetCompany(
		context.Background(),
		connectrpc.NewRequest(&companyv1.GetCompanyRequest{Ticker: ticker}),
	)
}

// replaceProfiles перестраивает сервер на другой mock ProfileRepository.
// Используется в тестах с особым поведением профиля; SetupTest по
// дефолту вешает «один и тот же набор полей на любой тикер».
func (s *serverSuite) replaceProfiles(profiles *company_mock.ProfileRepository) {
	s.T().Helper()
	s.profiles = profiles
	registry := data.NewRegistry()
	s.Require().NoError(registry.Register("moex", s.securityDescription))
	s.Require().NoError(registry.Register("financemarker", s.stockInfo))

	srv := pconnect.NewServer(pconnect.ConfigServer{
		Companies: services.NewCompanyService(services.ConfigCompanyService{
			Profiles: profiles,
			Registry: registry,
		}),
		Registry: registry,
	})

	mux := http.NewServeMux()
	path, handler := companyv1connect.NewCompanyServiceHandler(srv)
	mux.Handle(path, handler)

	s.server.Close()
	s.server = httptest.NewServer(mux)
	s.client = companyv1connect.NewCompanyServiceClient(s.server.Client(), s.server.URL)
}
