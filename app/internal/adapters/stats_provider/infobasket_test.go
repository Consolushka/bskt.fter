package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/infra/infobasket"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfobasketStatsProviderAdapter_GetPlayerBio(t *testing.T) {
	adapter := InfobasketStatsProviderAdapter{}
	bio, err := adapter.GetPlayerBio("any-id")

	require.NoError(t, err)
	assert.Empty(t, bio.FullName)
}

func TestInfobasketStatsProviderAdapter_GetGamesStatsByPeriod(t *testing.T) {
	compId := 123
	from := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC)

	t.Run("successfully fetches games in period", func(t *testing.T) {
		_ = os.RemoveAll("tmp/cache/")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := infobasket.NewMockClientInterface(ctrl)
		adapter := NewInfobasketStatsProviderAdapter(mockClient, compId)

		schedule := []infobasket.GameScheduleDto{
			{GameID: 1, GameDate: "10.02.2026", GameTime: "18:00", GameTimeMsk: "18:00"}, // In period
			{GameID: 2, GameDate: "05.03.2026", GameTime: "19:00", GameTimeMsk: "19:00"}, // Out of period
			{GameID: 3, GameDate: "15.02.2026", GameTime: "--:--", GameTimeMsk: "--:--"}, // Not scheduled
		}

		mockClient.EXPECT().ScheduledGames(compId).Return(schedule, nil)
		mockClient.EXPECT().BoxScore("1").Return(infobasket.GameBoxScoreResponse{
			GameStatus:  1, // Finished
			GameDate:    "10.02.2026",
			GameTimeMsk: "18.00",
			GameTeams: []infobasket.TeamBoxScoreDto{
				{
					Score: 80,
					TeamName: infobasket.TeamNameBoxScoreDto{
						CompTeamAbcNameEn: "LAL",
						CompTeamNameEn:    "Lakers",
					},
				},
				{
					Score: 75,
					TeamName: infobasket.TeamNameBoxScoreDto{
						CompTeamAbcNameEn: "BOS",
						CompTeamNameEn:    "Celtics",
					},
				},
			},
		}, nil)

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		require.NoError(t, err)
		assert.Len(t, games, 1)
		assert.Equal(t, "LAL - BOS", games[0].GameModel.Title)
	})

	t.Run("returns error when schedule loading fails", func(t *testing.T) {
		_ = os.RemoveAll("tmp/cache/")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := infobasket.NewMockClientInterface(ctrl)
		adapter := NewInfobasketStatsProviderAdapter(mockClient, compId)

		mockClient.EXPECT().ScheduledGames(compId).Return(nil, errors.New("network error"))

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		require.Error(t, err)
		assert.Nil(t, games)
	})

	t.Run("skips game when BoxScore fails", func(t *testing.T) {
		_ = os.RemoveAll("tmp/cache/")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := infobasket.NewMockClientInterface(ctrl)
		adapter := NewInfobasketStatsProviderAdapter(mockClient, compId)

		schedule := []infobasket.GameScheduleDto{
			{GameID: 1, GameDate: "10.02.2026", GameTime: "18:00", GameTimeMsk: "18:00"},
		}

		mockClient.EXPECT().ScheduledGames(compId).Return(schedule, nil)
		mockClient.EXPECT().BoxScore("1").Return(infobasket.GameBoxScoreResponse{}, errors.New("box score error"))

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		require.NoError(t, err)
		assert.Empty(t, games)
	})

	t.Run("skips game when not finished", func(t *testing.T) {
		_ = os.RemoveAll("tmp/cache/")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := infobasket.NewMockClientInterface(ctrl)
		adapter := NewInfobasketStatsProviderAdapter(mockClient, compId)

		schedule := []infobasket.GameScheduleDto{
			{GameID: 1, GameDate: "10.02.2026", GameTime: "18:00", GameTimeMsk: "18:00"},
		}

		mockClient.EXPECT().ScheduledGames(compId).Return(schedule, nil)
		mockClient.EXPECT().BoxScore("1").Return(infobasket.GameBoxScoreResponse{
			GameStatus: 0, // Scheduled/Live, but not finished
		}, nil)

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		require.NoError(t, err)
		assert.Empty(t, games)
	})

	t.Run("skips game when schedule date is invalid", func(t *testing.T) {
		_ = os.RemoveAll("tmp/cache/")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := infobasket.NewMockClientInterface(ctrl)
		adapter := NewInfobasketStatsProviderAdapter(mockClient, compId)

		schedule := []infobasket.GameScheduleDto{
			{GameID: 1, GameDate: "invalid-date", GameTime: "18:00"},
		}

		mockClient.EXPECT().ScheduledGames(compId).Return(schedule, nil)

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		require.NoError(t, err)
		assert.Empty(t, games)
	})

	t.Run("skips game when transformation fails", func(t *testing.T) {
		_ = os.RemoveAll("tmp/cache/")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := infobasket.NewMockClientInterface(ctrl)
		adapter := NewInfobasketStatsProviderAdapter(mockClient, compId)

		schedule := []infobasket.GameScheduleDto{
			{GameID: 1, GameDate: "10.02.2026", GameTime: "18:00", GameTimeMsk: "18:00"},
		}

		mockClient.EXPECT().ScheduledGames(compId).Return(schedule, nil)
		// Provide response with no teams to trigger transformation error
		mockClient.EXPECT().BoxScore("1").Return(infobasket.GameBoxScoreResponse{
			GameStatus:  1,
			GameDate:    "10.02.2026",
			GameTimeMsk: "18.00",
			GameTeams:   nil,
		}, nil)

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		require.NoError(t, err)
		assert.Empty(t, games)
	})
}

func TestInfobasketStatsProviderAdapter_EnrichGameStats(t *testing.T) {
	adapter := InfobasketStatsProviderAdapter{}
	game := games.GameStatEntity{GameModel: games.GameModel{Title: "LAL - BOS"}}

	enriched, err := adapter.EnrichGameStats(game)

	require.NoError(t, err)
	assert.Equal(t, game, enriched)
}
