package teams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamModel_TableName(t *testing.T) {
	model := TeamModel{}
	assert.Equal(t, "teams", model.TableName())
}

func TestTeamModel_AutoGenerateAlias(t *testing.T) {
	tests := []struct {
		name     string
		input    TeamModel
		expected string
	}{
		{
			name: "long name",
			input: TeamModel{
				Name: "Lakers",
			},
			expected: "Lak",
		},
		{
			name: "short name",
			input: TeamModel{
				Name: "OKC",
			},
			expected: "OKC",
		},
		{
			name: "2-char name",
			input: TeamModel{
				Name: "CS",
			},
			expected: "CS",
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
			assert.Equal(t, tt.expected, model.AutoGenerateAlias())
		})
	}
}
