package games

import (
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameStatEntity_AssignTeams(t *testing.T) {
	homeStats := teams.GameTeamStatModel{Score: 81}
	awayStats := teams.GameTeamStatModel{Score: 77}
	homePlayers := []players.PlayerStatisticEntity{{PlayerExternalId: "1"}}
	awayPlayers := []players.PlayerStatisticEntity{{PlayerExternalId: "2"}}

	game := GameStatEntity{
		GameModel: GameModel{Title: "UNIT - BK \"Ta"},
		HomeTeamStat: teams.TeamStatEntity{
			TeamModel:         teams.TeamModel{Name: "UNICS"},
			GameTeamStatModel: homeStats,
			PlayerStats:       homePlayers,
		},
		AwayTeamStat: teams.TeamStatEntity{
			TeamModel:         teams.TeamModel{Name: "BK \"Tambov\""},
			GameTeamStatModel: awayStats,
			PlayerStats:       awayPlayers,
		},
	}
	home := teams.TeamModel{Id: 1, Alias: "UNIT"}
	away := teams.TeamModel{Id: 2, Alias: "TAMBOV"}

	game.AssignTeams(home, away)

	assert.Equal(t, "UNIT - TAMBOV", game.GameModel.Title)
	assert.Equal(t, home, game.HomeTeamStat.TeamModel)
	assert.Equal(t, away, game.AwayTeamStat.TeamModel)
	assert.Equal(t, homeStats, game.HomeTeamStat.GameTeamStatModel)
	assert.Equal(t, awayStats, game.AwayTeamStat.GameTeamStatModel)
	assert.Equal(t, homePlayers, game.HomeTeamStat.PlayerStats)
	assert.Equal(t, awayPlayers, game.AwayTeamStat.PlayerStats)
}
