package scheduler

import (
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/core/tournaments"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
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

func TestReloadTournament(t *testing.T) {
	oldEnd := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	newEnd := time.Date(2026, 7, 30, 0, 0, 0, 0, time.UTC)
	loaded := tournaments.TournamentModel{Id: 4, EndAt: &oldEnd}

	t.Run("takes fresh tournament dates from repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := tournaments_repo.NewMockTournamentsRepo(ctrl)
		repo.EXPECT().Get(uint(4)).Return(tournaments.TournamentModel{Id: 4, EndAt: &newEnd}, nil)

		assert.Equal(t, &newEnd, reloadTournament(repo, loaded).EndAt)
	})

	t.Run("keeps loaded copy when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := tournaments_repo.NewMockTournamentsRepo(ctrl)
		repo.EXPECT().Get(uint(4)).Return(tournaments.TournamentModel{}, errors.New("db error"))

		assert.Equal(t, loaded, reloadTournament(repo, loaded))
	})
}
