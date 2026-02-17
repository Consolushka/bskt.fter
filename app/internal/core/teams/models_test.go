package teams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamModel_TableName(t *testing.T) {
	model := TeamModel{}
	assert.Equal(t, "teams", model.TableName())
}

func TestGameTeamStatModel_TableName(t *testing.T) {
	model := GameTeamStatModel{}
	assert.Equal(t, "game_team_stats", model.TableName())
}
