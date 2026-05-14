package securitydescription_test

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
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/bundles/securitydescription"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/client"
)

//go:embed testdata/*.json
var testdataFS embed.FS

type bundleSuite struct {
	suite.Suite

	handler func(http.ResponseWriter, *http.Request)
	server  *httptest.Server
	bundle  *securitydescription.Bundle
}

func TestBundleSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(bundleSuite))
}

func (s *bundleSuite) SetupTest() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handler(w, r)
	}))
	c := client.New(client.Config{
		BaseURL: s.server.URL + "/iss",
		Timeout: 5 * time.Second,
	})
	s.bundle = securitydescription.New(c)
}

func (s *bundleSuite) TearDownTest() {
	s.server.Close()
}

func (s *bundleSuite) TestMetadata() {
	s.Equal(securitydescription.ID, s.bundle.BundleID())
	s.Require().Len(s.bundle.Fields(), 26)
}

func (s *bundleSuite) TestFetchHappyPath() {
	body := s.readFixture("sber.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/iss/securities/SBER.json", r.URL.Path)
		s.Equal("extended", r.URL.Query().Get("iss.json"))
		s.Equal("off", r.URL.Query().Get("iss.meta"))
		s.Equal("description", r.URL.Query().Get("iss.only"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	values, err := s.bundle.Fetch(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Equal("SBER", values[domaincompany.FieldTicker])
	s.Equal("RU0009029540", values[domaincompany.FieldISIN])
	s.Equal("Сбербанк России ПАО ао", values[domaincompany.FieldName])
	s.Equal("Сбербанк", values[domaincompany.FieldShortName])
	s.Equal("Акции обыкновенные", values[domaincompany.FieldIssueName])
	s.Equal("Sberbank", values[domaincompany.FieldLatName])
	s.Equal("10301481B", values[domaincompany.FieldRegNumber])
	s.Equal("Акция обыкновенная", values[domaincompany.FieldSecurityTypeName])
	s.Equal("stock_shares", values[domaincompany.FieldSecurityGroup])
	s.Equal("Акции", values[domaincompany.FieldSecurityGroupName])
	s.Equal(domaincompany.SecurityTypeCommonShare, values[domaincompany.FieldSecurityType])
	s.Equal(domaincompany.ListingLevelFirst, values[domaincompany.FieldListingLevel])
	s.Equal("3", values[domaincompany.FieldFaceValue])
	s.Equal(domaincompany.CurrencyRUB, values[domaincompany.FieldFaceUnit])
	s.Equal(int64(21586948000), values[domaincompany.FieldIssueSize])
	s.Equal(time.Date(2007, 7, 20, 0, 0, 0, 0, time.UTC), values[domaincompany.FieldIssueDate])
	s.Equal(time.Date(2007, 7, 11, 0, 0, 0, 0, time.UTC), values[domaincompany.FieldRegistryDate])
	s.Equal("484", values[domaincompany.FieldEmitterID])
	s.Equal(false, values[domaincompany.FieldHasProspectus])
	s.Equal(false, values[domaincompany.FieldHasDefault])
	s.Equal(false, values[domaincompany.FieldHasTechnicalDefault])
	s.Equal(false, values[domaincompany.FieldEmitentMismatchCurrent])
	s.Equal(false, values[domaincompany.FieldIsQualifiedInvestors])
	s.Equal(true, values[domaincompany.FieldMorningSession])
	s.Equal(true, values[domaincompany.FieldEveningSession])
	s.Equal(true, values[domaincompany.FieldWeekendSession])
}

func (s *bundleSuite) TestFetchInvalidJSON() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}

	_, err := s.bundle.Fetch(context.Background(), "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "decode extended payload")
}

func (s *bundleSuite) TestFetchNotFound() {
	body := s.readFixture("not_found.json")
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.bundle.Fetch(context.Background(), "missing")

	s.Require().ErrorIs(err, domaincompany.ErrNotFound)
}

func (s *bundleSuite) TestFetchServerError() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err := s.bundle.Fetch(context.Background(), "any")

	s.Require().Error(err)
	s.ErrorContains(err, "moex http status 500")
}

func (s *bundleSuite) TestFetchContextCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.bundle.Fetch(ctx, "any")

	s.Require().Error(err)
	s.ErrorContains(err, "context canceled")
	s.ErrorContains(err, "moex request")
}

// TestFetchTypeAndLevelMatrix проходит по всем ожидаемым значениям полей
// TYPE и LISTLEVEL блока description — каждое значение в маппинг-таблицах
// и ветка-fallback для неизвестного TYPE проверяются хотя бы один раз.
func (s *bundleSuite) TestFetchTypeAndLevelMatrix() {
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
			values, err := s.bundle.Fetch(context.Background(), "X")
			s.Require().NoError(err)
			s.Equal(c.wantType, values[domaincompany.FieldSecurityType])
			s.Equal(c.wantLevel, values[domaincompany.FieldListingLevel])
		})
	}
}

func (s *bundleSuite) TestFetchInvalidListLevel() {
	body := []byte(`[{"description":[{"name":"SECID","value":"X"},{"name":"LISTLEVEL","value":"9"}]}]`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.bundle.Fetch(context.Background(), "any")

	s.Require().Error(err)
	s.ErrorContains(err, "LISTLEVEL")
}

func (s *bundleSuite) TestFetchMissingDescriptionBlock() {
	body := []byte(`[{"charsetinfo":{"name":"utf-8"}},{"securities":[]}]`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.bundle.Fetch(context.Background(), "any")

	s.Require().Error(err)
	s.ErrorContains(err, "description block missing")
}

func (s *bundleSuite) TestFetchInvalidDescriptionBlock() {
	// description присутствует, но это не массив объектов {name, value}.
	body := []byte(`[{"description": 123}]`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.bundle.Fetch(context.Background(), "any")

	s.Require().Error(err)
	s.ErrorContains(err, "decode description block")
}

func (s *bundleSuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}
