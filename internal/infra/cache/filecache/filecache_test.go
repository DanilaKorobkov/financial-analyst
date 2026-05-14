package filecache_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/infra/cache/filecache"
)

type payload struct {
	IssuedAt  time.Time `json:"issued_at"`
	Name      string    `json:"name"`
	FaceValue int64     `json:"face_value"`
	Listed    bool      `json:"listed"`
}

type fileCacheSuite struct {
	suite.Suite

	dir string
}

func TestFileCacheSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(fileCacheSuite))
}

func (s *fileCacheSuite) SetupTest() {
	s.dir = s.T().TempDir()
}

func (s *fileCacheSuite) TestPutGetRoundTrip() {
	cache := filecache.New[payload](filecache.Config{Dir: s.dir})
	want := payload{Name: "SBER", FaceValue: 3, IssuedAt: time.Date(2007, 7, 11, 0, 0, 0, 0, time.UTC), Listed: true}

	s.Require().NoError(cache.Put("SBER", want, time.Hour))

	got, ok := cache.Get("SBER")
	s.Require().True(ok)
	s.Equal(want, got)
}

func (s *fileCacheSuite) TestGetMissReturnsFalse() {
	cache := filecache.New[payload](filecache.Config{Dir: s.dir})

	_, ok := cache.Get("missing")
	s.False(ok)
}

func (s *fileCacheSuite) TestExpiredEntryReportedAsMiss() {
	cache := filecache.New[payload](filecache.Config{Dir: s.dir})
	s.Require().NoError(cache.Put("SBER", payload{Name: "old"}, 10*time.Millisecond))

	time.Sleep(20 * time.Millisecond)

	_, ok := cache.Get("SBER")
	s.False(ok)
}

func (s *fileCacheSuite) TestZeroTTLNoExpiration() {
	cache := filecache.New[payload](filecache.Config{Dir: s.dir})
	want := payload{Name: "forever"}

	s.Require().NoError(cache.Put("SBER", want, 0))

	got, ok := cache.Get("SBER")
	s.Require().True(ok)
	s.Equal(want, got)
}

func (s *fileCacheSuite) TestCorruptedFileTreatedAsMiss() {
	cache := filecache.New[payload](filecache.Config{Dir: s.dir})

	s.Require().NoError(os.WriteFile(filepath.Join(s.dir, "SBER.json"), []byte("not json"), 0o600))

	_, ok := cache.Get("SBER")
	s.False(ok)
}

func (s *fileCacheSuite) TestKeyIsEscapedForFilesystem() {
	cache := filecache.New[payload](filecache.Config{Dir: s.dir})
	want := payload{Name: "exotic"}

	s.Require().NoError(cache.Put("WEIRD/KEY:1", want, time.Hour))

	entries, err := os.ReadDir(s.dir)
	s.Require().NoError(err)
	s.Require().Len(entries, 1)
	s.NotContains(entries[0].Name(), "/")

	got, ok := cache.Get("WEIRD/KEY:1")
	s.Require().True(ok)
	s.Equal(want, got)
}

func (s *fileCacheSuite) TestLoadOrFetchFetchesOnMiss() {
	cache := filecache.New[payload](filecache.Config{Dir: s.dir})
	want := payload{Name: "SBER"}

	calls := 0
	got, err := cache.LoadOrFetch(context.Background(), "SBER", time.Hour, func(_ context.Context) (payload, error) {
		calls++
		return want, nil
	})

	s.Require().NoError(err)
	s.Equal(want, got)
	s.Equal(1, calls)
}

func (s *fileCacheSuite) TestLoadOrFetchHitsCacheOnSecondCall() {
	cache := filecache.New[payload](filecache.Config{Dir: s.dir})
	want := payload{Name: "SBER"}

	_, err := cache.LoadOrFetch(context.Background(), "SBER", time.Hour, func(_ context.Context) (payload, error) {
		return want, nil
	})
	s.Require().NoError(err)

	calls := 0
	got, err := cache.LoadOrFetch(context.Background(), "SBER", time.Hour, func(_ context.Context) (payload, error) {
		calls++
		return payload{}, errors.New("should not be called")
	})

	s.Require().NoError(err)
	s.Equal(want, got)
	s.Equal(0, calls)
}

func (s *fileCacheSuite) TestLoadOrFetchPropagatesError() {
	cache := filecache.New[payload](filecache.Config{Dir: s.dir})
	boom := errors.New("boom")

	_, err := cache.LoadOrFetch(context.Background(), "any", time.Hour, func(_ context.Context) (payload, error) {
		return payload{}, boom
	})

	s.Require().ErrorIs(err, boom)
}
