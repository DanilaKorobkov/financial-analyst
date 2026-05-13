package company_test

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex"
	moexcompany "github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/company"
)

//go:embed testdata/*.json
var testdataFS embed.FS

type repositorySuite struct {
	suite.Suite

	handler func(http.ResponseWriter, *http.Request)
	server  *httptest.Server
	repo    *moexcompany.Repository
}

func TestRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(repositorySuite))
}

func (s *repositorySuite) SetupTest() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handler(w, r)
	}))
	client := moex.NewClient(moex.ConfigClient{
		BaseURL: s.server.URL + "/iss",
		Timeout: 5 * time.Second,
	})
	s.repo = moexcompany.NewRepository(client)
}

func (s *repositorySuite) TearDownTest() {
	s.server.Close()
}

func (s *repositorySuite) TestFindByTickerHappyPath() {
	body := s.readFixture("sber.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/iss/securities/SBER.json", r.URL.Path)
		s.Equal("extended", r.URL.Query().Get("iss.json"))
		s.Equal("off", r.URL.Query().Get("iss.meta"))
		s.Equal("description", r.URL.Query().Get("iss.only"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	found, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	expected := domaincompany.Company{
		Ticker:       "SBER",
		ISIN:         "RU0009029540",
		Name:         "Сбербанк России ПАО ао",
		SecurityType: domaincompany.SecurityTypeCommonShare,
		ListingLevel: domaincompany.ListingLevelFirst,
	}
	s.Equal(expected, found)
}

func (s *repositorySuite) TestFindByTickerInvalidJSON() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}

	_, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "decode extended payload")
}

func (s *repositorySuite) TestFindByTickerNotFound() {
	body := s.readFixture("not_found.json")
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.repo.FindByTicker(context.Background(), "ZZZZ")

	s.Require().ErrorIs(err, domaincompany.ErrNotFound)
}

func (s *repositorySuite) TestFindByTickerServerError() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err := s.repo.FindByTicker(context.Background(), "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "moex http status 500")
}

func (s *repositorySuite) TestFindByTickerContextCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.repo.FindByTicker(ctx, "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "context canceled")
	s.ErrorContains(err, "moex request")
}

// TestFindByTickerTypeAndLevelMatrix проходит по всем ожидаемым значениям полей
// TYPE и LISTLEVEL блока description — каждое значение в маппинг-таблицах
// и ветка-fallback для неизвестного TYPE проверяются хотя бы один раз.
func (s *repositorySuite) TestFindByTickerTypeAndLevelMatrix() {
	cases := []struct {
		name      string
		typeValue string
		level     string
		wantType  domaincompany.SecurityType
		wantLevel domaincompany.ListingLevel
	}{
		{
			name:      "preferred share with second level",
			typeValue: "preferred_share",
			level:     "2",
			wantType:  domaincompany.SecurityTypePreferredShare,
			wantLevel: domaincompany.ListingLevelSecond,
		},
		{
			name:      "depositary receipt with third level",
			typeValue: "depositary_receipt",
			level:     "3",
			wantType:  domaincompany.SecurityTypeDepositaryReceipt,
			wantLevel: domaincompany.ListingLevelThird,
		},
		{
			name:      "неизвестный TYPE падает в Unspecified",
			typeValue: "exotic_new_type",
			level:     "1",
			wantType:  domaincompany.SecurityTypeUnspecified,
			wantLevel: domaincompany.ListingLevelFirst,
		},
		{
			name:      "пустой LISTLEVEL трактуется как Unspecified",
			typeValue: "common_share",
			level:     "",
			wantType:  domaincompany.SecurityTypeCommonShare,
			wantLevel: domaincompany.ListingLevelUnspecified,
		},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			body := []byte(fmt.Sprintf(
				`[{"description":[{"name":"SECID","value":"X"},{"name":"TYPE","value":%q},{"name":"LISTLEVEL","value":%q}]}]`,
				c.typeValue, c.level,
			))
			s.handler = func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write(body)
			}
			found, err := s.repo.FindByTicker(context.Background(), "X")
			s.Require().NoError(err)
			s.Equal(c.wantType, found.SecurityType)
			s.Equal(c.wantLevel, found.ListingLevel)
		})
	}
}

func (s *repositorySuite) TestFindByTickerInvalidListLevel() {
	body := []byte(`[{"description":[{"name":"SECID","value":"X"},{"name":"LISTLEVEL","value":"9"}]}]`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.repo.FindByTicker(context.Background(), "X")

	s.Require().Error(err)
	s.ErrorContains(err, "LISTLEVEL")
}

func (s *repositorySuite) TestFindByTickerMissingDescriptionBlock() {
	body := []byte(`[{"charsetinfo":{"name":"utf-8"}},{"securities":[]}]`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.repo.FindByTicker(context.Background(), "X")

	s.Require().Error(err)
	s.ErrorContains(err, "description block missing")
}

func (s *repositorySuite) TestFindByTickerInvalidDescriptionBlock() {
	// description присутствует, но это не массив объектов {name, value}.
	body := []byte(`[{"description": 123}]`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.repo.FindByTicker(context.Background(), "X")

	s.Require().Error(err)
	s.ErrorContains(err, "decode description block")
}

func (s *repositorySuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}
