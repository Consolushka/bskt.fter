package api_basketball

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/teams"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntityTransformer_TransformWithoutPlayers(t *testing.T) {
	transformer := EntityTransformer{}

	inputGame := GameEntity{
		Id:   454569,
		Date: time.Date(2026, 4, 12, 19, 15, 0, 0, time.UTC),
		Teams: Teams{
			Home: Team{
				Id:   687,
				Name: "Maccabi Tel Aviv",
			},
			Away: Team{
				Id:   682,
				Name: "Hapoel Tel-Aviv",
			},
		},
		Scores: Scores{
			Home: ScoreDetails{
				Total: 88,
			},
			Away: ScoreDetails{
				Total: 99,
			},
		},
	}

	result := transformer.TransformWithoutPlayers(inputGame)

	assert.Equal(t, time.Date(2026, 4, 12, 19, 15, 0, 0, time.UTC), result.GameModel.ScheduledAt)
	assert.Equal(t, "Maccabi Tel Aviv - Hapoel Tel-Aviv", result.GameModel.Title)
	assert.Equal(t, 48, result.GameModel.Duration)
	assert.Equal(t, "Maccabi Tel Aviv", result.HomeTeamStat.TeamModel.Name)
	assert.Equal(t, 88, result.HomeTeamStat.GameTeamStatModel.Score)
	assert.Equal(t, -11, result.HomeTeamStat.GameTeamStatModel.FinalDiff)
}

func TestEntityTransformer_MapPlayerStatistics(t *testing.T) {
	transformer := EntityTransformer{}

	percentage := 80.0
	response := PlayerStatsResponse{
		Response: []PlayerStatsEntity{
			{
				Player: PlayerRef{
					Id:   3534,
					Name: "Hoard Jaylen",
				},
				Team:    TeamRef{Id: 687},
				Minutes: "34:41",
				Points:  12,
				Rebounds: PlayerRebounds{
					Total: 9,
				},
				Assists: 1,
				FieldGoals: StatsDetails{
					Percentage: &percentage,
				},
			},
		},
	}

	gameEntity := &games.GameStatEntity{
		HomeTeamStat: teams.TeamStatEntity{},
		AwayTeamStat: teams.TeamStatEntity{},
	}

	err := transformer.MapPlayerStatistics(response, 687, 682, gameEntity)

	require.NoError(t, err)
	require.Len(t, gameEntity.HomeTeamStat.PlayerStats, 1)
	require.Empty(t, gameEntity.AwayTeamStat.PlayerStats)

	homePlayer := gameEntity.HomeTeamStat.PlayerStats[0]
	assert.Equal(t, "3534", homePlayer.PlayerExternalId)
	assert.Equal(t, "Hoard Jaylen", homePlayer.PlayerModel.FullName)
	assert.InEpsilon(t, float32(0.8), homePlayer.GameTeamPlayerStatModel.FieldGoalsPercentage, 0.01)
	assert.Equal(t, uint8(12), homePlayer.GameTeamPlayerStatModel.Points)
}

func TestNewEntityTransformer(t *testing.T) {
	transformer := NewEntityTransformer()
	assert.NotNil(t, transformer)
}
