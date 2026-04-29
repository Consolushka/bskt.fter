package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/infra/api_basketball"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiBasketballStatsProviderAdapter_GetGamesStatsByPeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := api_basketball.NewMockClientInterface(ctrl)
	adapter := NewApiBasketballStatsProviderAdapter(mockClient, 12)

	from := time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 4, 12, 23, 59, 59, 0, time.UTC)

	mockClient.EXPECT().
		Games(0, "2026-04-12", "", "", "", "").
		Return(api_basketball.GamesResponse{
			Response: []api_basketball.GameEntity{
				{
					Id:   454569,
					Date: from,
					Status: api_basketball.Status{
						Short: "FT",
					},
					League: api_basketball.League{
						Id: 12,
					},
					Teams: api_basketball.Teams{
						Home: api_basketball.Team{Id: 687, Name: "Maccabi"},
						Away: api_basketball.Team{Id: 682, Name: "Hapoel"},
					},
					Scores: api_basketball.Scores{
						Home: api_basketball.ScoreDetails{Total: 88},
						Away: api_basketball.ScoreDetails{Total: 99},
					},
				},
				{
					Id: 111, // Wrong league
					League: api_basketball.League{
						Id: 999,
					},
				},
			},
		}, nil)

	result, err := adapter.GetGamesStatsByPeriod(from, to)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "454569", result[0].ExternalGameId)
}

func TestApiBasketballStatsProviderAdapter_EnrichGameStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := api_basketball.NewMockClientInterface(ctrl)
	adapter := NewApiBasketballStatsProviderAdapter(mockClient, 12)

	game := games.GameStatEntity{
		ExternalGameId:     "454569",
		HomeTeamExternalId: 687,
		AwayTeamExternalId: 682,
	}

	mockClient.EXPECT().
		PlayersStatistics(454569, 0, 0).
		Return(api_basketball.PlayerStatsResponse{
			Response: []api_basketball.PlayerStatsEntity{
				{
					Player:  api_basketball.PlayerRef{Id: 3534, Name: "Hoard Jaylen"},
					Team:    api_basketball.TeamRef{Id: 687},
					Points:  12,
					Minutes: "30:00",
				},
			},
		}, nil)

	result, err := adapter.EnrichGameStats(game)

	require.NoError(t, err)
	require.Len(t, result.HomeTeamStat.PlayerStats, 1)
	assert.Equal(t, "Hoard Jaylen", result.HomeTeamStat.PlayerStats[0].PlayerModel.FullName)
}

func TestApiBasketballStatsProviderAdapter_GetPlayerBio(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := api_basketball.NewMockClientInterface(ctrl)
	adapter := NewApiBasketballStatsProviderAdapter(mockClient, 12)

	mockClient.EXPECT().
		PlayerInfo(3534).
		Return(api_basketball.PlayerInfoResponse{
			Response: []api_basketball.PlayerInfoEntity{
				{
					Id:   3534,
					Name: "Hoard Jaylen",
				},
			},
		}, nil)

	result, err := adapter.GetPlayerBio("3534")

	require.NoError(t, err)
	assert.Equal(t, "Hoard Jaylen", result.FullName)
}
