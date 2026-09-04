package api_basketball

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/teams"
	"encoding/json"
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

func TestEntityTransformer_MapLeaguePeriod(t *testing.T) {
	transformer := EntityTransformer{}

	resp := LeaguesResponse{
		Response: []LeagueInfoEntity{
			{
				Id:   120,
				Name: "Euroleague",
				Seasons: []SeasonEntity{
					{Season: "2024", Start: "2024-10-03", End: "2025-05-25"},
					{Season: "2025", Start: "2025-09-30", End: "2026-05-24"},
				},
			},
		},
	}

	t.Run("returns parsed start and end of requested season", func(t *testing.T) {
		start, end, err := transformer.MapLeaguePeriod(resp, "2025")

		require.NoError(t, err)
		assert.Equal(t, time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC), start)
		assert.Equal(t, time.Date(2026, 5, 24, 0, 0, 0, 0, time.UTC), end)
	})

	t.Run("error when season not found", func(t *testing.T) {
		_, _, err := transformer.MapLeaguePeriod(resp, "2099")
		require.Error(t, err)
	})

	t.Run("error on empty response", func(t *testing.T) {
		_, _, err := transformer.MapLeaguePeriod(LeaguesResponse{}, "2025")
		require.Error(t, err)
	})

	t.Run("error on unparseable date", func(t *testing.T) {
		bad := LeaguesResponse{Response: []LeagueInfoEntity{{
			Seasons: []SeasonEntity{{Season: "2025", Start: "not-a-date", End: "2026-05-24"}},
		}}}
		_, _, err := transformer.MapLeaguePeriod(bad, "2025")
		require.Error(t, err)
	})
}

func TestSeasonValue_UnmarshalJSON(t *testing.T) {
	t.Run("numeric season decodes to string", func(t *testing.T) {
		var s SeasonEntity
		require.NoError(t, json.Unmarshal([]byte(`{"season": 2025, "start": "2025-09-30", "end": "2026-05-24"}`), &s))
		assert.Equal(t, SeasonValue("2025"), s.Season)
	})

	t.Run("string season decodes as-is", func(t *testing.T) {
		var s SeasonEntity
		require.NoError(t, json.Unmarshal([]byte(`{"season": "2023-2024", "start": "2023-10-05", "end": "2024-05-26"}`), &s))
		assert.Equal(t, SeasonValue("2023-2024"), s.Season)
	})
}
