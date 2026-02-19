package service

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewTournamentProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPersistence := NewMockPersistenceServiceInterface(ctrl)
	mockStatsProvider := stats_provider.NewMockStatsProvider(ctrl)
	mockPlayersRepo := players_repo.NewMockPlayersRepo(ctrl)
	mockGamesRepo := games_repo.NewMockGamesRepo(ctrl)

	processor := NewTournamentProcessor(mockStatsProvider, mockPersistence, mockPlayersRepo, mockGamesRepo, 12)

	assert.NotNil(t, processor)
	assert.Equal(t, uint(12), processor.tournamentId)
	assert.Equal(t, mockPersistence, processor.persistenceService)
	assert.Equal(t, mockStatsProvider, processor.statsProvider)
	assert.Equal(t, mockPlayersRepo, processor.playersRepo)
	assert.Equal(t, mockGamesRepo, processor.gamesRepo)
}

func TestTournamentProcessor_ProcessByPeriod(t *testing.T) {
	from := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 2, 14, 23, 59, 59, 0, time.UTC)

	t.Run("returns error when stats provider fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return(nil, errors.New("stats provider error"))

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, 1)
		err := processor.ProcessByPeriod(from, to)

		assert.ErrorContains(t, err, "stats provider error")
	})

	t.Run("skips existing games without enrichment and saving", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		gameEntities := []games.GameStatEntity{
			{GameModel: games.GameModel{Title: "LAL - BOS", ScheduledAt: from}},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return(gameEntities, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(true, nil)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, 99)
		err := processor.ProcessByPeriod(from, to)

		assert.NoError(t, err)
	})

	t.Run("enriches and saves new game", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		gameEntity := games.GameStatEntity{GameModel: games.GameModel{Title: "MIA - NYK", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{gameEntity}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(gameEntity, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, 77)
		err := processor.ProcessByPeriod(from, to)

		assert.NoError(t, err)
	})

	t.Run("returns error when game existence check fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		gameEntity := games.GameStatEntity{GameModel: games.GameModel{Title: "DAL - PHX", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{gameEntity}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, errors.New("db error"))

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, 11)
		err := processor.ProcessByPeriod(from, to)

		assert.ErrorContains(t, err, "db error")
	})

	t.Run("continues when enrich game stats fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game1 := games.GameStatEntity{GameModel: games.GameModel{Title: "GAME-1", ScheduledAt: from}}
		game2 := games.GameStatEntity{GameModel: games.GameModel{Title: "GAME-2", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game1, game2}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil).Times(2)
		gomock.InOrder(
			mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(games.GameStatEntity{}, errors.New("enrich error")),
			mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game2, nil),
		)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, 20)
		err := processor.ProcessByPeriod(from, to)

		assert.NoError(t, err)
	})

	t.Run("continues when save game fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game1 := games.GameStatEntity{GameModel: games.GameModel{Title: "SAVE-1", ScheduledAt: from}}
		game2 := games.GameStatEntity{GameModel: games.GameModel{Title: "SAVE-2", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game1, game2}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil).Times(2)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game1, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game2, nil)
		gomock.InOrder(
			mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(errors.New("save error")),
			mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil),
		)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, 21)
		err := processor.ProcessByPeriod(from, to)

		assert.NoError(t, err)
	})

	t.Run("continues when ListByFullName fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "PLAYERS-ERR", ScheduledAt: from}}
		enriched := games.GameStatEntity{
			GameModel: game.GameModel,
			HomeTeamStat: teams.TeamStatEntity{
				PlayerStats: []players.PlayerStatisticEntity{
					{PlayerModel: players.PlayerModel{FullName: "John Doe"}},
				},
			},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(enriched, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockPlayers.EXPECT().ListByFullName("John Doe").Return(nil, errors.New("players repo error"))

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, 22)
		err := processor.ProcessByPeriod(from, to)

		assert.NoError(t, err)
	})

	t.Run("continues when player not found and GetPlayerBio fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "BIO-ERR", ScheduledAt: from}}
		enriched := games.GameStatEntity{
			GameModel: game.GameModel,
			HomeTeamStat: teams.TeamStatEntity{
				PlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "123",
						PlayerModel:      players.PlayerModel{FullName: ""},
					},
				},
			},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(enriched, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockPlayers.EXPECT().ListByFullName("").Return([]players.PlayerModel{}, nil)
		mockStats.EXPECT().GetPlayerBio("123").Return(players.PlayerBioEntity{}, errors.New("bio error"))

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, 23)
		err := processor.ProcessByPeriod(from, to)

		assert.NoError(t, err)
	})

	t.Run("continues when player not found and GetPlayerBio succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "BIO-OK", ScheduledAt: from}}
		enriched := games.GameStatEntity{
			GameModel: game.GameModel,
			HomeTeamStat: teams.TeamStatEntity{
				PlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "555",
						PlayerModel:      players.PlayerModel{FullName: "", BirthDate: time.Time{}},
					},
				},
			},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(enriched, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockPlayers.EXPECT().ListByFullName("").Return([]players.PlayerModel{}, nil)
		mockStats.EXPECT().GetPlayerBio("555").Return(players.PlayerBioEntity{
			FullName:  "Recovered Player",
			BirthDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, 24)
		err := processor.ProcessByPeriod(from, to)

		assert.NoError(t, err)
	})
}
