package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/infra/financemarker/client"
)

type clientSuite struct {
	suite.Suite
}

func TestClientSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(clientSuite))
}

// TestCacheDirWrapsTransport: при непустом CacheDir клиент оборачивает
// transport в httpcache и продолжает обслуживать запросы.
func (s *clientSuite) TestCacheDirWrapsTransport() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := client.New(client.Config{
		BaseURL:  srv.URL,
		Token:    "t",
		Timeout:  time.Second,
		CacheDir: s.T().TempDir(),
	})

	resp, err := c.R().Get("/")
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode())
}
