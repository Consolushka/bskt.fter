package games

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameModel_TableName(t *testing.T) {
	model := GameModel{}
	assert.Equal(t, "games", model.TableName())
}
