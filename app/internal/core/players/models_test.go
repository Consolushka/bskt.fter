package players

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerModel_TableName(t *testing.T) {
	model := PlayerModel{}
	assert.Equal(t, "players", model.TableName())
}

func TestGameTeamPlayerStatModel_TableName(t *testing.T) {
	model := GameTeamPlayerStatModel{}
	assert.Equal(t, "game_team_player_stats", model.TableName())
}
