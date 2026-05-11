//go:build integration

// Package integration — вертикальный срез приложения через сгенерированный
// Connect-клиент. Поднимает fake-MOEX, собирает app.New, прокидывает запросы
// клиента через всю цепочку.
package integration

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/app"
	"github.com/DanilaKorobkov/financial-analyst/gen/company/v1/companyv1connect"
)

type IntegrationSuite struct {
	suite.Suite

	moexHandler http.HandlerFunc
	moex        *httptest.Server
	api         *httptest.Server
	client      companyv1connect.CompanyServiceClient
}

func (s *IntegrationSuite) SetupSuite() {
	s.moex = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.moexHandler(w, r)
	}))

	handler := app.New(app.Config{
		Moex: app.MoexConfig{
			BaseURL: s.moex.URL + "/iss",
			Timeout: 5 * time.Second,
		},
		Server: app.ServerConfig{Port: 0}, // не используется внутри httptest
	})

	s.api = httptest.NewServer(handler)
	s.client = companyv1connect.NewCompanyServiceClient(s.api.Client(), s.api.URL)
}

func (s *IntegrationSuite) TearDownSuite() {
	s.api.Close()
	s.moex.Close()
}

func (s *IntegrationSuite) SetupTest() {
	s.moexHandler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// setMoexFixture настраивает fake-MOEX отдавать содержимое файла testdata
// для запросов c указанным ticker в URL path.
func (s *IntegrationSuite) setMoexFixture(ticker, fixture string) {
	body, err := os.ReadFile(filepath.Join("testdata", fixture)) //nolint:gosec // фикстура из testdata/
	s.Require().NoError(err)
	s.moexHandler = func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/"+ticker+".json") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write(body)
	}
}

func TestIntegrationSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IntegrationSuite))
}
