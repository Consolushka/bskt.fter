package teams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamModel_TableName(t *testing.T) {
	model := TeamModel{}
	assert.Equal(t, "teams", model.TableName())
}

func TestTeamModel_BeforeSave(t *testing.T) {
	tests := []struct {
		name     string
		input    TeamModel
		expected string
	}{
		{
			name: "empty alias, long name",
			input: TeamModel{
				Name: "Lakers",
			},
			expected: "Lak",
		},
		{
			name: "empty alias, short name",
			input: TeamModel{
				Name: "OKC",
			},
			expected: "OKC",
		},
		{
			name: "empty alias, 2-char name",
			input: TeamModel{
				Name: "CS",
			},
			expected: "CS",
		},
		{
			name: "existing alias, should not change",
			input: TeamModel{
				Name:  "Lakers",
				Alias: "LAL",
			},
			expected: "LAL",
		},
		{
			name: "multibyte characters",
			input: TeamModel{
				Name: "ЦСКА",
			},
			expected: "ЦСК",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.input
			_ = model.BeforeSave(nil)
			assert.Equal(t, tt.expected, model.Alias)
		})
	}
}
