package scheduler

import (
	"IMP/app/internal/core/tournaments"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextPollAt(t *testing.T) {
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	interval := 30 * time.Minute
	future := now.AddDate(0, 0, 1)
	past := now.AddDate(0, 0, -1)

	t.Run("running tournament gets next poll time", func(t *testing.T) {
		next := nextPollAt(tournaments.TournamentModel{EndAt: &future}, now, interval)

		require.NotNil(t, next)
		assert.Equal(t, now.Add(interval), *next)
	})

	t.Run("tournament without end date gets next poll time", func(t *testing.T) {
		next := nextPollAt(tournaments.TournamentModel{}, now, interval)

		require.NotNil(t, next)
		assert.Equal(t, now.Add(interval), *next)
	})

	t.Run("ended tournament gets no next poll time", func(t *testing.T) {
		assert.Nil(t, nextPollAt(tournaments.TournamentModel{EndAt: &past}, now, interval))
	})
}
