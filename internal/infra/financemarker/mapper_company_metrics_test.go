// потребовал бы экспорта функций, нужных только тестам.
//
//nolint:testpackage // whitebox-тесты приватных мапперов: вынос в *_test
package financemarker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

func TestParseIdeaConsensus(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want entities.IdeaConsensus
	}{
		{"BUY", entities.IdeaConsensusBuy},
		{"HOLD", entities.IdeaConsensusHold},
		{"SELL", entities.IdeaConsensusSell},
		{"", entities.IdeaConsensusUnspecified},
		{"UNKNOWN", entities.IdeaConsensusUnspecified},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, parseIdeaConsensus(c.in), c.in)
	}
}

func TestParseInsiderConsensus(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want entities.InsiderConsensus
	}{
		{"BUYS", entities.InsiderConsensusBuys},
		{"SELLS", entities.InsiderConsensusSells},
		{"MIXED", entities.InsiderConsensusMixed},
		{"", entities.InsiderConsensusUnspecified},
		{"WTF", entities.InsiderConsensusUnspecified},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, parseInsiderConsensus(c.in), c.in)
	}
}

func TestParseChangedAt(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		got := parseChangedAt("2026-05-11T03:32:06")
		require.Equal(t, time.Date(2026, 5, 11, 3, 32, 6, 0, time.UTC), got)
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		assert.True(t, parseChangedAt("").IsZero())
	})
	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		assert.True(t, parseChangedAt("11.05.2026").IsZero())
	})
}

func TestMapCompanyMetricsMinimal(t *testing.T) {
	t.Parallel()

	// Только info — summary пустой, как при отсутствии данных у эмитента.
	dto := stockDTO{
		Info: infoDTO{
			Code:     "VTBR",
			Exchange: "MOEX",
			Name:     "ВТБ",
			Sector:   "Финансы",
		},
	}

	got := mapCompanyMetrics(&dto)

	assert.Equal(t, "VTBR", got.Card.Ticker)
	assert.Equal(t, "ВТБ", got.Card.Name)
	assert.Equal(t, entities.IdeaConsensusUnspecified, got.IdeaConsensus)
	assert.Equal(t, entities.InsiderConsensusUnspecified, got.InsiderConsensus)
	assert.True(t, got.ChangedAt.IsZero())
	assert.InDelta(t, 0.0, got.GrowthEbitda3y, 0.0001)
}
