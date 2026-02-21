package service

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/tournament_poll_logs_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/core/tournaments"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewTournamentsOrchestrator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	persistence := NewMockPersistenceServiceInterface(ctrl)
	tournamentsRepo := tournaments_repo.NewMockTournamentsRepo(ctrl)
	playersRepo := players_repo.NewMockPlayersRepo(ctrl)
	gamesRepo := games_repo.NewMockGamesRepo(ctrl)
	pollLogRepo := tournament_poll_logs_repo.NewMockTournamentPollLogsRepo(ctrl)

	result := NewTournamentsOrchestrator(persistence, tournamentsRepo, playersRepo, gamesRepo, pollLogRepo)

	assert.NotNil(t, result)
	assert.Equal(t, persistence, result.persistenceService)
	assert.Equal(t, tournamentsRepo, result.tournamentsRepo)
	assert.Equal(t, playersRepo, result.playersRepo)
	assert.Equal(t, gamesRepo, result.gamesRepo)
	assert.Equal(t, pollLogRepo, result.pollLogRepo)
}

func TestTournamentsOrchestrator_ProcessAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	persistence := NewMockPersistenceServiceInterface(ctrl)
	tournamentsRepo := tournaments_repo.NewMockTournamentsRepo(ctrl)
	playersRepo := players_repo.NewMockPlayersRepo(ctrl)
	gamesRepo := games_repo.NewMockGamesRepo(ctrl)
	pollLogRepo := tournament_poll_logs_repo.NewMockTournamentPollLogsRepo(ctrl)

	orchestrator := NewTournamentsOrchestrator(persistence, tournamentsRepo, playersRepo, gamesRepo, pollLogRepo)
	from := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 2, 14, 23, 59, 59, 0, time.UTC)

	t.Run("returns repo error", func(t *testing.T) {
		tournamentsRepo.EXPECT().ListActive().Return(nil, errors.New("repository error"))
		err := orchestrator.ProcessAll(from, to)
		assert.ErrorContains(t, err, "repository error")
	})

	t.Run("success with empty tournaments", func(t *testing.T) {
		tournamentsRepo.EXPECT().ListActive().Return([]tournaments.TournamentModel{}, nil)
		err := orchestrator.ProcessAll(from, to)
		assert.NoError(t, err)
	})

	t.Run("processes non-empty tournaments and enters goroutine loop branches", func(t *testing.T) {
		tournamentsRepo.EXPECT().ListActive().Return([]tournaments.TournamentModel{
			{
				Id: 1,
				Provider: tournaments.TournamentProvider{
					ProviderName: "API_NBA",
					Params:       []byte("{invalid-json}"),
				},
			},
			{
				Id: 2,
				Provider: tournaments.TournamentProvider{
					ProviderName: "UNKNOWN_PROVIDER",
					Params:       nil,
				},
			},
		}, nil)

		err := orchestrator.ProcessAll(from, to)
		assert.NoError(t, err)
	})
}
