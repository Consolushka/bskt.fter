package tournament_poll_logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTournamentPollLogModel_TableName(t *testing.T) {
	model := TournamentPollLogModel{}
	assert.Equal(t, "tournament_poll_logs", model.TableName())
}
