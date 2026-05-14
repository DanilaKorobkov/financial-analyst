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
	fieldIDs            []string
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

	s.fieldIDs = []string{
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
	s.Equal("SBER", fields[company.FieldTicker].GetStringValue())
	s.Equal("RU0009029540", fields[company.FieldISIN].GetStringValue())
	s.Equal("Сбербанк России ПАО ао", fields[company.FieldName].GetStringValue())
	s.Equal("common_share", fields[company.FieldSecurityType].GetStringValue())
	s.Equal("first", fields[company.FieldListingLevel].GetStringValue())
	s.Equal("Сбербанк", fields[company.FieldIssuerName].GetStringValue())
	s.Equal("moex", fields[company.FieldExchange].GetStringValue())
	s.Equal("RUB", fields[company.FieldCurrency].GetStringValue())
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
	s.Equal("preferred_share", fields[company.FieldSecurityType].GetStringValue())
	s.Equal("second", fields[company.FieldListingLevel].GetStringValue())
	s.Equal("moex", fields[company.FieldExchange].GetStringValue())
	s.Equal("USD", fields[company.FieldCurrency].GetStringValue())
	s.Equal("quarterly", fields[company.FieldReportFrequency].GetStringValue())
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

	_, found := resp.Msg.GetFields()[company.FieldSecurityType]
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
	s.Require().NotNil(fields[company.FieldIssueDate].GetTimestampValue())
	s.Equal(issueDate.Unix(), fields[company.FieldIssueDate].GetTimestampValue().GetSeconds())
	_, found := fields[company.FieldRegistryDate]
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
	s.Equal(int64(21586948000), fields[company.FieldIssueSize].GetIntValue())
	s.False(fields[company.FieldHasDefault].GetBoolValue())
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
		Return(company.Profile{FieldIDs: []string{"moex::unknown"}}, nil).
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
