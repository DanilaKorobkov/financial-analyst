package httpcache_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/infra/cache/filecache"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/cache/httpcache"
)

// countingTransport — обёртка над http.RoundTripper, считающая вызовы.
type countingTransport struct {
	inner http.RoundTripper
	calls int
}

// erroringTransport — всегда возвращает заданную ошибку, в сеть не ходит.
type erroringTransport struct {
	err error
}

// fetched — материализованный ответ после чтения и закрытия body.
// Тесты работают с этой структурой, чтобы *http.Response не утекал
// наружу из doGet и не порождал предупреждения bodyclose.
type fetched struct {
	headers    http.Header
	body       string
	statusCode int
}

type transportSuite struct {
	suite.Suite

	server  *httptest.Server
	handler func(http.ResponseWriter, *http.Request)
	dir     string
}

func (c *countingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	c.calls++
	return c.inner.RoundTrip(req)
}

func (e *erroringTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, e.err
}

func TestTransportSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(transportSuite))
}

func (s *transportSuite) SetupTest() {
	s.dir = s.T().TempDir()
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handler(w, r)
	}))
}

func (s *transportSuite) TearDownTest() {
	s.server.Close()
}

func (s *transportSuite) TestNoTTLBypassesCache() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("payload"))
	}
	transport, counter := s.newTransport()
	client := &http.Client{Transport: transport}

	for range 3 {
		got := s.fetch(client, s.server.URL+"/stocks", 0)
		s.Equal("payload", got.body)
	}

	s.Equal(3, counter.calls)
}

func (s *transportSuite) TestSecondCallReturnsFromCache() {
	calls := 0
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"v":1}`))
	}
	transport, counter := s.newTransport()
	client := &http.Client{Transport: transport}

	first := s.fetch(client, s.server.URL+"/stocks?include=info", time.Hour)
	s.Equal(`{"v":1}`, first.body)
	s.Equal("application/json", first.headers.Get("Content-Type"))

	second := s.fetch(client, s.server.URL+"/stocks?include=info", time.Hour)
	s.Equal(`{"v":1}`, second.body)
	s.Equal("application/json", second.headers.Get("Content-Type"))

	s.Equal(1, calls, "upstream должен быть позван один раз")
	s.Equal(1, counter.calls, "wrapped transport должен быть позван один раз")
}

func (s *transportSuite) TestNon2xxNotCached() {
	calls := 0
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusNotFound)
	}
	transport, _ := s.newTransport()
	client := &http.Client{Transport: transport}

	first := s.fetch(client, s.server.URL+"/stocks", time.Hour)
	s.Equal(http.StatusNotFound, first.statusCode)

	second := s.fetch(client, s.server.URL+"/stocks", time.Hour)
	s.Equal(http.StatusNotFound, second.statusCode)

	s.Equal(2, calls, "404 не должен кешироваться")
}

func (s *transportSuite) TestExpiredEntryRefetched() {
	calls := 0
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		calls++
		_, _ = w.Write([]byte("payload"))
	}
	transport, _ := s.newTransport()
	client := &http.Client{Transport: transport}

	s.fetch(client, s.server.URL+"/stocks", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	s.fetch(client, s.server.URL+"/stocks", 10*time.Millisecond)

	s.Equal(2, calls)
}

func (s *transportSuite) TestIgnoredQueryParamCollapsesKey() {
	calls := 0
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		calls++
		_, _ = w.Write([]byte("payload"))
	}
	transport, _ := s.newTransport()
	client := &http.Client{Transport: transport}

	s.fetch(client, s.server.URL+"/stocks?api_token=A&include=info", time.Hour)
	s.fetch(client, s.server.URL+"/stocks?api_token=B&include=info", time.Hour)

	s.Equal(1, calls, "разный api_token не должен порождать новую запись в кеше")
}

func (s *transportSuite) TestDifferentMeaningfulQueryParamsCacheSeparately() {
	calls := 0
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		calls++
		_, _ = w.Write([]byte("payload"))
	}
	transport, _ := s.newTransport()
	client := &http.Client{Transport: transport}

	s.fetch(client, s.server.URL+"/stocks?include=info", time.Hour)
	s.fetch(client, s.server.URL+"/stocks?include=ratios", time.Hour)

	s.Equal(2, calls)
}

func (s *transportSuite) TestQueryParamsOrderIndependent() {
	calls := 0
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		calls++
		_, _ = w.Write([]byte("payload"))
	}
	transport, _ := s.newTransport()
	client := &http.Client{Transport: transport}

	s.fetch(client, s.server.URL+"/stocks?a=1&b=2", time.Hour)
	s.fetch(client, s.server.URL+"/stocks?b=2&a=1", time.Hour)

	s.Equal(1, calls, "порядок query-параметров не должен влиять на ключ")
}

func (s *transportSuite) TestRequestKeyIncludesNoQueryWhenAbsent() {
	calls := 0
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		calls++
		_, _ = w.Write([]byte("payload"))
	}
	transport, _ := s.newTransport()
	client := &http.Client{Transport: transport}

	s.fetch(client, s.server.URL+"/stocks", time.Hour)
	s.fetch(client, s.server.URL+"/stocks", time.Hour)

	s.Equal(1, calls)
}

func (s *transportSuite) TestNetworkErrorPropagated() {
	boom := errors.New("network down")
	store := filecache.New[httpcache.Entry](filecache.Config{Dir: s.dir})
	transport := httpcache.NewTransport(&erroringTransport{err: boom}, store)

	req, err := http.NewRequestWithContext(
		httpcache.WithTTL(context.Background(), time.Hour),
		http.MethodGet,
		"http://example/x",
		strings.NewReader(""),
	)
	s.Require().NoError(err)

	resp, err := transport.RoundTrip(req)
	if resp != nil {
		_ = resp.Body.Close()
	}
	s.Require().ErrorIs(err, boom)
}

func (s *transportSuite) TestWithTTLZeroOrNegativeIsNoOp() {
	ctx := context.Background()
	s.Equal(ctx, httpcache.WithTTL(ctx, 0))
	s.Equal(ctx, httpcache.WithTTL(ctx, -time.Second))
}

func (s *transportSuite) TestSpecialQueryValuesAreKeyed() {
	calls := 0
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		calls++
		_, _ = w.Write([]byte("payload"))
	}
	transport, _ := s.newTransport()
	client := &http.Client{Transport: transport}

	v := url.QueryEscape("MOEX:SBER ZZ")
	s.fetch(client, s.server.URL+"/stocks?symbol="+v, time.Hour)
	s.fetch(client, s.server.URL+"/stocks?symbol="+v, time.Hour)

	s.Equal(1, calls)
}

// helpers — ниже всех тестов

func (s *transportSuite) newTransport() (*httpcache.Transport, *countingTransport) {
	s.T().Helper()
	store := filecache.New[httpcache.Entry](filecache.Config{Dir: s.dir})
	counter := &countingTransport{inner: http.DefaultTransport}
	return httpcache.NewTransport(counter, store), counter
}

// fetch делает GET, сразу читает и закрывает body, возвращает уже
// материализованную fetched. Так *http.Response не утекает из helper'а
// в тесты — это и читаемее, и нет ложных предупреждений bodyclose.
func (s *transportSuite) fetch(client *http.Client, urlStr string, ttl time.Duration) fetched {
	s.T().Helper()

	ctx := context.Background()
	if ttl > 0 {
		ctx = httpcache.WithTTL(ctx, ttl)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, http.NoBody)
	s.Require().NoError(err)

	resp, err := client.Do(req)
	s.Require().NoError(err)
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	return fetched{
		statusCode: resp.StatusCode,
		headers:    resp.Header.Clone(),
		body:       string(body),
	}
}
