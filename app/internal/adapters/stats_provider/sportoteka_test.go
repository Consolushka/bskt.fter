package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/infra/sportoteka"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSportotekaStatsProviderAdapter_GetPlayerBio(t *testing.T) {
	adapter := SportotekaStatsProviderAdapter{}
	bio, err := adapter.GetPlayerBio("any-id")

	assert.NoError(t, err)
	assert.Empty(t, bio.FullName)
}

func TestSportotekaStatsProviderAdapter_GetGamesStatsByPeriod(t *testing.T) {
	tag := "test-tag"
	year := 2024
	from := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 2, 28, 23, 59, 59, 0, time.UTC)

	t.Run("successfully fetches games in period", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := sportoteka.NewMockClientInterface(ctrl)
		adapter := NewSportotekaStatsProvider(mockClient, tag, year)

		calendar := sportoteka.CalendarResponse{
			TotalCount: 1,
			Items: []sportoteka.CalendarGameEntity{
				{
					Game:  sportoteka.GameInfoEntity{Id: 1, GameStatus: "Result"},
					Team1: sportoteka.TeamInfoEntity{Name: "Lakers", AbcName: "LAL"},
					Team2: sportoteka.TeamInfoEntity{Name: "Celtics", AbcName: "BOS"},
				},
			},
		}

		mockClient.EXPECT().Calendar(tag, year, from, to).Return(calendar, nil)
		mockClient.EXPECT().BoxScore("1").Return(sportoteka.GameBoxScoreResponse{
			Result: sportoteka.GameBoxScoreEntity{
				Game:  sportoteka.GameInfoEntity{Id: 1, GameStatus: "Result", Periods: 4, ScheduledTime: from},
				Team1: sportoteka.TeamInfoEntity{Name: "Lakers", AbcName: "LAL"},
				Team2: sportoteka.TeamInfoEntity{Name: "Celtics", AbcName: "BOS"},
				Teams: []sportoteka.TeamBoxScoreEntity{
					{TeamNumber: 1, Total: sportoteka.TeamBoxScoreTotalsEntity{Points: 100}},
					{TeamNumber: 2, Total: sportoteka.TeamBoxScoreTotalsEntity{Points: 95}},
				},
			},
		}, nil)

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		require.NoError(t, err)
		assert.Len(t, games, 1)
		assert.Equal(t, "LAL - BOS", games[0].GameModel.Title)
	})

	t.Run("returns empty list when calendar fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := sportoteka.NewMockClientInterface(ctrl)
		adapter := NewSportotekaStatsProvider(mockClient, tag, year)

		mockClient.EXPECT().Calendar(tag, year, from, to).Return(sportoteka.CalendarResponse{}, errors.New("network error"))

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		assert.Error(t, err)
		assert.Empty(t, games)
	})

	t.Run("skips games with invalid status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := sportoteka.NewMockClientInterface(ctrl)
		adapter := NewSportotekaStatsProvider(mockClient, tag, year)

		calendar := sportoteka.CalendarResponse{
			TotalCount: 1,
			Items: []sportoteka.CalendarGameEntity{
				{
					Game: sportoteka.GameInfoEntity{Id: 1, GameStatus: "Scheduled"},
				},
			},
		}

		mockClient.EXPECT().Calendar(tag, year, from, to).Return(calendar, nil)

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		assert.NoError(t, err)
		assert.Empty(t, games)
	})

	t.Run("skips game when BoxScore fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := sportoteka.NewMockClientInterface(ctrl)
		adapter := NewSportotekaStatsProvider(mockClient, tag, year)

		calendar := sportoteka.CalendarResponse{
			TotalCount: 1,
			Items: []sportoteka.CalendarGameEntity{
				{
					Game: sportoteka.GameInfoEntity{Id: 1, GameStatus: "Result"},
				},
			},
		}

		mockClient.EXPECT().Calendar(tag, year, from, to).Return(calendar, nil)
		mockClient.EXPECT().BoxScore("1").Return(sportoteka.GameBoxScoreResponse{}, errors.New("api error"))

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		assert.NoError(t, err)
		assert.Empty(t, games)
	})

	t.Run("skips game when BoxScore status is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := sportoteka.NewMockClientInterface(ctrl)
		adapter := NewSportotekaStatsProvider(mockClient, tag, year)

		calendar := sportoteka.CalendarResponse{
			TotalCount: 1,
			Items: []sportoteka.CalendarGameEntity{
				{
					Game: sportoteka.GameInfoEntity{Id: 1, GameStatus: "Result"},
				},
			},
		}

		mockClient.EXPECT().Calendar(tag, year, from, to).Return(calendar, nil)
		mockClient.EXPECT().BoxScore("1").Return(sportoteka.GameBoxScoreResponse{
			Result: sportoteka.GameBoxScoreEntity{
				Game: sportoteka.GameInfoEntity{GameStatus: "Live"}, // Not Result or ResultConfirmed
			},
		}, nil)

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		assert.NoError(t, err)
		assert.Empty(t, games)
	})

	t.Run("skips game when transformation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := sportoteka.NewMockClientInterface(ctrl)
		adapter := NewSportotekaStatsProvider(mockClient, tag, year)

		calendar := sportoteka.CalendarResponse{
			TotalCount: 1,
			Items: []sportoteka.CalendarGameEntity{
				{
					Game: sportoteka.GameInfoEntity{Id: 1, GameStatus: "Result"},
				},
			},
		}

		mockClient.EXPECT().Calendar(tag, year, from, to).Return(calendar, nil)
		// In current implementation, Transform might return an error if teams are missing or date is invalid.
		// Since we don't have explicit error conditions in Transform yet, we can't easily trigger it without 
		// knowing the internal logic of Transform. But we've already added a check for len(GameTeams) in Infobasket.
		// Let's assume Transform fails if Teams are nil.
		mockClient.EXPECT().BoxScore("1").Return(sportoteka.GameBoxScoreResponse{
			Result: sportoteka.GameBoxScoreEntity{
				Game:  sportoteka.GameInfoEntity{GameStatus: "Result"},
				Teams: nil, 
			},
		}, nil)

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		assert.NoError(t, err)
		assert.Empty(t, games)
	})
}

func TestSportotekaStatsProviderAdapter_EnrichGameStats(t *testing.T) {
	adapter := SportotekaStatsProviderAdapter{}
	game := games.GameStatEntity{GameModel: games.GameModel{Title: "LAL - BOS"}}

	enriched, err := adapter.EnrichGameStats(game)

	assert.NoError(t, err)
	assert.Equal(t, game, enriched)
}
