package api_nba

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntityTransformer_TransformWithoutPlayers(t *testing.T) {
	transformer := EntityTransformer{}

	inputGame := GameEntity{
		Id:     12345,
		Season: 2023,
		Date: GameDateEntity{
			Start: time.Date(2023, 12, 25, 20, 0, 0, 0, time.UTC),
		},
		Teams: GameTeamsEntity{
			Home: TeamEntity{
				Id:       1,
				Name:     "Los Angeles Lakers",
				Nickname: "Lakers",
				Code:     "LAL",
			},
			Visitors: TeamEntity{
				Id:       2,
				Name:     "Boston Celtics",
				Nickname: "Celtics",
				Code:     "BOS",
			},
		},
		Scores: GameTeamsScoresEntity{
			Home: TeamScoresEntity{
				Points: 115,
			},
			Visitors: TeamScoresEntity{
				Points: 108,
			},
		},
		Periods: GamePeriodsEntity{
			Total: 4,
		},
	}

	result := transformer.TransformWithoutPlayers(inputGame)

	assert.Equal(t, time.Date(2023, 12, 25, 20, 0, 0, 0, time.UTC), result.GameModel.ScheduledAt)
	assert.Equal(t, "LAL - BOS", result.GameModel.Title)
	assert.Equal(t, 48, result.GameModel.Duration)
	assert.Equal(t, "Lakers", result.HomeTeamStat.TeamModel.Name)
	assert.Equal(t, "Los Angeles", result.HomeTeamStat.TeamModel.HomeTown)
	assert.Equal(t, 115, result.HomeTeamStat.GameTeamStatModel.Score)
	assert.Equal(t, 7, result.HomeTeamStat.GameTeamStatModel.FinalDiff)
}

func TestEntityTransformer_MapPlayerStatistics(t *testing.T) {
	transformer := EntityTransformer{}

	response := PlayerStatisticResponse{
		Response: []PlayerStatisticEntity{
			{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        101,
					Firstname: "LeBron",
					Lastname:  "James",
				},
				Team:      TeamEntity{Id: 1},
				Min:       "35:24",
				PlusMinus: "12",
				Points:    27,
				TotReb:    8,
				Assists:   10,
				Steals:    2,
				Blocks:    1,
				Fgp:       "50.5",
				Turnovers: 3,
			},
			{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        102,
					Firstname: "Jayson",
					Lastname:  "Tatum",
				},
				Team:      TeamEntity{Id: 2},
				Min:       "38:15",
				PlusMinus: "-7",
				Points:    32,
				TotReb:    12,
				Assists:   5,
				Steals:    1,
				Blocks:    0,
				Fgp:       "45.2",
				Turnovers: 4,
			},
		},
	}

	gameEntity := &games.GameStatEntity{
		HomeTeamStat: teams.TeamStatEntity{},
		AwayTeamStat: teams.TeamStatEntity{},
	}

	err := transformer.MapPlayerStatistics(response, 1, 2, gameEntity)

	require.NoError(t, err)
	require.Len(t, gameEntity.HomeTeamStat.PlayerStats, 1)
	require.Len(t, gameEntity.AwayTeamStat.PlayerStats, 1)

	homePlayer := gameEntity.HomeTeamStat.PlayerStats[0]
	assert.Equal(t, "101", homePlayer.PlayerExternalId)
	assert.Equal(t, "LeBron James", homePlayer.PlayerModel.FullName)
	assert.InEpsilon(t, float32(0.505), homePlayer.GameTeamPlayerStatModel.FieldGoalsPercentage, 0.01)
	assert.Equal(t, uint8(27), homePlayer.GameTeamPlayerStatModel.Points)

	awayPlayer := gameEntity.AwayTeamStat.PlayerStats[0]
	assert.Equal(t, "102", awayPlayer.PlayerExternalId)
	assert.Equal(t, "Jayson Tatum", awayPlayer.PlayerModel.FullName)
	assert.InEpsilon(t, float32(0.452), awayPlayer.GameTeamPlayerStatModel.FieldGoalsPercentage, 0.01)
}

func TestEntityTransformer_mapPlayerStatistic_Errors(t *testing.T) {
	transformer := EntityTransformer{}

	t.Run("invalid_minutes", func(t *testing.T) {
		player := PlayerStatisticEntity{Min: "invalid:time"}
		var result players.PlayerStatisticEntity
		err := transformer.mapPlayerStatistic(player, &result)
		assert.Error(t, err)
	})

	t.Run("invalid_plus_minus", func(t *testing.T) {
		player := PlayerStatisticEntity{Min: "30:00", PlusMinus: "not_a_number"}
		var result players.PlayerStatisticEntity
		err := transformer.mapPlayerStatistic(player, &result)
		assert.Error(t, err)
	})

	t.Run("invalid_fgp", func(t *testing.T) {
		player := PlayerStatisticEntity{Min: "30:00", PlusMinus: "5", Fgp: "not_a_float"}
		var result players.PlayerStatisticEntity
		err := transformer.mapPlayerStatistic(player, &result)
		assert.Error(t, err)
	})
}

func TestNewEntityTransformer(t *testing.T) {
	transformer := NewEntityTransformer()
	assert.NotNil(t, transformer)
}
