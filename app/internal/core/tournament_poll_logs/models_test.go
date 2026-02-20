package tournament_poll_logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPollWatermarkModel_TableName(t *testing.T) {
	model := PollWatermarkModel{}
	assert.Equal(t, "poll_watermarks", model.TableName())
}
