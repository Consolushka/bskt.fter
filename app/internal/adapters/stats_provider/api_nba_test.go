package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/infra/api_nba"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiNbaStatsProviderAdapter_GetPlayerBio(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := api_nba.NewMockClientInterface(ctrl)
	adapter := NewApiNbaStatsProviderAdapter(mockClient)

	t.Run("successfully fetches player bio", func(t *testing.T) {
		playerId := "123"
		mockClient.EXPECT().
			PlayerInfo(123, "", 0, 0, "", "").
			Return(api_nba.PlayersResponse{
				Response: []api_nba.PlayerEntity{
					{
						Firstname: "John",
						Lastname:  "Doe",
						Birth:     api_nba.PlayerBirthEntity{Date: "1990-01-01"},
					},
				},
			}, nil)

		expectedBio := players.PlayerBioEntity{
			FullName:  "John Doe",
			BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		bio, err := adapter.GetPlayerBio(playerId)

		require.NoError(t, err)
		assert.Equal(t, expectedBio, bio)
	})

	t.Run("returns error when id is not a number", func(t *testing.T) {
		_, err := adapter.GetPlayerBio("invalid")
		assert.ErrorContains(t, err, "atoi")
	})

	t.Run("returns error when PlayerInfo API call fails", func(t *testing.T) {
		mockClient.EXPECT().
			PlayerInfo(123, "", 0, 0, "", "").
			Return(api_nba.PlayersResponse{}, errors.New("network error"))

		_, err := adapter.GetPlayerBio("123")
		assert.ErrorContains(t, err, "playerInfo")
		assert.ErrorContains(t, err, "network error")
	})

	t.Run("returns error when API response is empty", func(t *testing.T) {
		mockClient.EXPECT().
			PlayerInfo(456, "", 0, 0, "", "").
			Return(api_nba.PlayersResponse{Response: []api_nba.PlayerEntity{}}, nil)

		_, err := adapter.GetPlayerBio("456")
		assert.EqualError(t, err, "empty player info response")
	})

	t.Run("returns error when birth date format is invalid", func(t *testing.T) {
		mockClient.EXPECT().
			PlayerInfo(789, "", 0, 0, "", "").
			Return(api_nba.PlayersResponse{
				Response: []api_nba.PlayerEntity{
					{
						Birth: api_nba.PlayerBirthEntity{Date: "invalid-date"},
					},
				},
			}, nil)

		_, err := adapter.GetPlayerBio("789")
		assert.ErrorContains(t, err, "time.Parse")
	})
}

func TestApiNbaStatsProviderAdapter_GetGamesStatsByPeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := api_nba.NewMockClientInterface(ctrl)
	adapter := NewApiNbaStatsProviderAdapter(mockClient)

	t.Run("successfully fetches games for a single day", func(t *testing.T) {
		date := time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC)
		dateStr := "2024-02-14"

		mockClient.EXPECT().
			Games(0, dateStr, "1", "", "", "").
			Return(api_nba.GamesResponse{
				Response: []api_nba.GameEntity{
					{
						Id:     1001,
						Status: api_nba.GameStatusEntity{Short: 3},
						Teams: api_nba.GameTeamsEntity{
							Home:     api_nba.TeamEntity{Id: 1, Code: "LAL"},
							Visitors: api_nba.TeamEntity{Id: 2, Code: "BOS"},
						},
						Date: api_nba.GameDateEntity{Start: date},
					},
					{
						Id:     1002,
						Status: api_nba.GameStatusEntity{Short: 1}, // Scheduled, should be skipped
					},
				},
			}, nil)

		games, err := adapter.GetGamesStatsByPeriod(date, date)

		require.NoError(t, err)
		assert.Len(t, games, 1)
		assert.Equal(t, "1001", games[0].ExternalGameId)
		assert.Equal(t, 1, games[0].HomeTeamExternalId)
		assert.Equal(t, 2, games[0].AwayTeamExternalId)
	})

	t.Run("successfully fetches games for a date range (2 calls)", func(t *testing.T) {
		from := time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

		mockClient.EXPECT().
			Games(0, "2024-02-14", "1", "", "", "").
			Return(api_nba.GamesResponse{
				Response: []api_nba.GameEntity{
					{Id: 1001, Status: api_nba.GameStatusEntity{Short: 3}},
				},
			}, nil)

		mockClient.EXPECT().
			Games(0, "2024-02-15", "1", "", "", "").
			Return(api_nba.GamesResponse{
				Response: []api_nba.GameEntity{
					{Id: 1002, Status: api_nba.GameStatusEntity{Short: 3}},
				},
			}, nil)

		games, err := adapter.GetGamesStatsByPeriod(from, to)

		require.NoError(t, err)
		assert.Len(t, games, 2)
	})

	t.Run("returns error when first API call fails (from date)", func(t *testing.T) {
		from := time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

		mockClient.EXPECT().
			Games(0, "2024-02-14", "1", "", "", "").
			Return(api_nba.GamesResponse{}, errors.New("from call failed"))

		_, err := adapter.GetGamesStatsByPeriod(from, to)
		assert.ErrorContains(t, err, "from call failed")
	})

	t.Run("returns error when second API call fails (to date)", func(t *testing.T) {
		from := time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

		mockClient.EXPECT().
			Games(0, "2024-02-14", "1", "", "", "").
			Return(api_nba.GamesResponse{}, nil)

		mockClient.EXPECT().
			Games(0, "2024-02-15", "1", "", "", "").
			Return(api_nba.GamesResponse{}, errors.New("to call failed"))

		_, err := adapter.GetGamesStatsByPeriod(from, to)
		assert.ErrorContains(t, err, "to call failed")
	})

	t.Run("returns error when single API call fails", func(t *testing.T) {
		date := time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC)
		mockClient.EXPECT().
			Games(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(api_nba.GamesResponse{}, errors.New("api error"))

		_, err := adapter.GetGamesStatsByPeriod(date, date)
		assert.ErrorContains(t, err, "api error")
	})
}

func TestApiNbaStatsProviderAdapter_EnrichGameStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := api_nba.NewMockClientInterface(ctrl)
	adapter := NewApiNbaStatsProviderAdapter(mockClient)

	t.Run("returns game as is if ExternalGameId is empty", func(t *testing.T) {
		inputGame := games.GameStatEntity{ExternalGameId: ""}
		result, err := adapter.EnrichGameStats(inputGame)

		require.NoError(t, err)
		assert.Equal(t, inputGame, result)
	})

	t.Run("returns error when ExternalGameId is not a number", func(t *testing.T) {
		inputGame := games.GameStatEntity{ExternalGameId: "not-a-number"}
		_, err := adapter.EnrichGameStats(inputGame)
		assert.ErrorContains(t, err, "atoi")
	})

	t.Run("successfully enriches game stats", func(t *testing.T) {
		game := games.GameStatEntity{
			ExternalGameId:     "1001",
			HomeTeamExternalId: 1,
			AwayTeamExternalId: 2,
		}

		mockClient.EXPECT().
			PlayersStatistics(0, 1001, 0, "").
			Return(api_nba.PlayerStatisticResponse{
				Response: []api_nba.PlayerStatisticEntity{
					{
						Player:    api_nba.PlayerStatisticsPlayerGeneralDataEntity{Id: 101, Firstname: "LeBron", Lastname: "James"},
						Team:      api_nba.TeamEntity{Id: 1},
						Points:    30,
						Min:       "30:00",
						Fgp:       "50.0",
						PlusMinus: "10",
					},
				},
			}, nil)

		enriched, err := adapter.EnrichGameStats(game)

		require.NoError(t, err)
		assert.Len(t, enriched.HomeTeamStat.PlayerStats, 1)
		assert.Equal(t, "LeBron James", enriched.HomeTeamStat.PlayerStats[0].PlayerModel.FullName)
	})

	t.Run("returns error when MapPlayerStatistics fails", func(t *testing.T) {
		game := games.GameStatEntity{ExternalGameId: "1001"}

		mockClient.EXPECT().
			PlayersStatistics(0, 1001, 0, "").
			Return(api_nba.PlayerStatisticResponse{
				Response: []api_nba.PlayerStatisticEntity{
					{Min: "invalid-format"}, // Causes error in mapPlayerStatistic
				},
			}, nil)

		_, err := adapter.EnrichGameStats(game)
		assert.ErrorContains(t, err, "MapPlayerStatistics returned error")
	})

	t.Run("returns error when PlayersStatistics API call fails", func(t *testing.T) {
		game := games.GameStatEntity{ExternalGameId: "1001"}
		mockClient.EXPECT().
			PlayersStatistics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(api_nba.PlayerStatisticResponse{}, errors.New("api error"))

		_, err := adapter.EnrichGameStats(game)
		assert.ErrorContains(t, err, "api error")
	})
}
