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

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/aggregates/company"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/bundles/securitydescription"
	"github.com/DanilaKorobkov/financial-analyst/internal/infra/moex/client"
)

//go:embed testdata/*.json
var testdataFS embed.FS

type sourceSuite struct {
	suite.Suite

	handler func(http.ResponseWriter, *http.Request)
	server  *httptest.Server
	source  *securitydescription.Source
}

func TestSourceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(sourceSuite))
}

func (s *sourceSuite) SetupTest() {
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
	s.source = securitydescription.New(c)
}

func (s *sourceSuite) TearDownTest() {
	s.server.Close()
}

func (s *sourceSuite) TestFindByTickerHappyPath() {
	body := s.readFixture("sber.json")
	s.handler = func(w http.ResponseWriter, r *http.Request) {
		s.Equal("/iss/securities/SBER.json", r.URL.Path)
		s.Equal("extended", r.URL.Query().Get("iss.json"))
		s.Equal("off", r.URL.Query().Get("iss.meta"))
		s.Equal("description", r.URL.Query().Get("iss.only"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}

	got, err := s.source.FindByTicker(context.Background(), "SBER")

	s.Require().NoError(err)
	s.Require().NotNil(got)
	s.Equal("SBER", got.Ticker)
	s.Equal("RU0009029540", got.ISIN)
	s.Equal("Сбербанк России ПАО ао", got.Name)
	s.Equal("Сбербанк", got.ShortName)
	s.Equal("Акции обыкновенные", got.IssueName)
	s.Equal("Sberbank", got.LatName)
	s.Equal("10301481B", got.RegNumber)
	s.Equal("Акция обыкновенная", got.SecurityTypeName)
	s.Equal("stock_shares", got.SecurityGroup)
	s.Equal("Акции", got.SecurityGroupName)
	s.Equal(company.SecurityTypeCommonShare, got.SecurityType)
	s.Equal(company.ListingLevelFirst, got.ListingLevel)
	s.Equal("3", got.FaceValue)
	s.Equal(company.CurrencyRUB, got.FaceUnit)
	s.Equal(int64(21586948000), got.IssueSize)
	s.Equal(time.Date(2007, 7, 20, 0, 0, 0, 0, time.UTC), got.IssueDate)
	s.Equal(time.Date(2007, 7, 11, 0, 0, 0, 0, time.UTC), got.RegistryDate)
	s.Equal("484", got.EmitterID)
	s.False(got.HasProspectus)
	s.False(got.HasDefault)
	s.False(got.HasTechnicalDefault)
	s.False(got.EmitentMismatchCurrent)
	s.False(got.IsQualifiedInvestors)
	s.True(got.MorningSession)
	s.True(got.EveningSession)
	s.True(got.WeekendSession)
}

func (s *sourceSuite) TestFindByTickerInvalidJSON() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}

	_, err := s.source.FindByTicker(context.Background(), "SBER")

	s.Require().Error(err)
	s.ErrorContains(err, "decode extended payload")
}

func (s *sourceSuite) TestFindByTickerNotFound() {
	body := s.readFixture("not_found.json")
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	got, err := s.source.FindByTicker(context.Background(), "missing")

	s.Require().ErrorIs(err, company.ErrNotFound)
	s.Nil(got)
}

func (s *sourceSuite) TestFindByTickerServerError() {
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err := s.source.FindByTicker(context.Background(), "any")

	s.Require().Error(err)
	s.ErrorContains(err, "moex http status 500")
}

func (s *sourceSuite) TestFindByTickerContextCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.source.FindByTicker(ctx, "any")

	s.Require().Error(err)
	s.Require().ErrorContains(err, "context canceled")
	s.Require().ErrorContains(err, "moex request")
}

// TestFindByTickerTypeAndLevelMatrix проходит по всем ожидаемым значениям
// полей TYPE и LISTLEVEL блока description.
func (s *sourceSuite) TestFindByTickerTypeAndLevelMatrix() {
	cases := []struct {
		name      string
		typeValue string
		level     string
		wantType  company.SecurityType
		wantLevel company.ListingLevel
	}{
		{
			name:      "preferred share with second level",
			typeValue: "preferred_share",
			level:     "2",
			wantType:  company.SecurityTypePreferredShare,
			wantLevel: company.ListingLevelSecond,
		},
		{
			name:      "depositary receipt with third level",
			typeValue: "depositary_receipt",
			level:     "3",
			wantType:  company.SecurityTypeDepositaryReceipt,
			wantLevel: company.ListingLevelThird,
		},
		{
			name:      "неизвестный TYPE падает в Unspecified",
			typeValue: "exotic_new_type",
			level:     "1",
			wantType:  company.SecurityTypeUnspecified,
			wantLevel: company.ListingLevelFirst,
		},
		{
			name:      "пустой LISTLEVEL трактуется как Unspecified",
			typeValue: "common_share",
			level:     "",
			wantType:  company.SecurityTypeCommonShare,
			wantLevel: company.ListingLevelUnspecified,
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
			got, err := s.source.FindByTicker(context.Background(), "X")
			s.Require().NoError(err)
			s.Equal(c.wantType, got.SecurityType)
			s.Equal(c.wantLevel, got.ListingLevel)
		})
	}
}

func (s *sourceSuite) TestFindByTickerInvalidListLevel() {
	body := []byte(`[{"description":[{"name":"SECID","value":"X"},{"name":"LISTLEVEL","value":"9"}]}]`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.source.FindByTicker(context.Background(), "any")

	s.Require().Error(err)
	s.ErrorContains(err, "LISTLEVEL")
}

func (s *sourceSuite) TestFindByTickerMissingDescriptionBlock() {
	body := []byte(`[{"charsetinfo":{"name":"utf-8"}},{"securities":[]}]`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.source.FindByTicker(context.Background(), "any")

	s.Require().Error(err)
	s.ErrorContains(err, "description block missing")
}

func (s *sourceSuite) TestFindByTickerInvalidDescriptionBlock() {
	body := []byte(`[{"description": 123}]`)
	s.handler = func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	_, err := s.source.FindByTicker(context.Background(), "any")

	s.Require().Error(err)
	s.ErrorContains(err, "decode description block")
}

func (s *sourceSuite) readFixture(name string) []byte {
	s.T().Helper()
	raw, err := testdataFS.ReadFile("testdata/" + name)
	s.Require().NoError(err)
	return raw
}
