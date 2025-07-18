package leagues

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestLeagueModel_TableName tests the TableName method
// Verify that when TableName is called while getting table name - returns "leagues" string
func TestLeagueModel_TableName(t *testing.T) {
	cases := []struct {
		name     string
		data     struct{}
		expected string
		errorMsg string
	}{
		{
			name:     "returns correct table name",
			data:     struct{}{},
			expected: "leagues",
			errorMsg: "",
		},
	}

	model := LeagueModel{
		Id:        1,
		Name:      "Test League",
		Alias:     "TL",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result := model.TableName()

			assert.Equal(t, tc.expected, result)
		})
	}
}
