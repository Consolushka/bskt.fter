package tournaments

import (
	"IMP/app/internal/core/leagues"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestTournamentModel_TableName tests the TableName method
// Verify that when TableName is called while getting table name - returns "tournaments" string
func TestTournamentModel_TableName(t *testing.T) {
	cases := []struct {
		name     string
		data     struct{}
		expected string
		errorMsg string
	}{
		{
			name:     "returns correct table name",
			data:     struct{}{},
			expected: "tournaments",
			errorMsg: "",
		},
	}

	model := TournamentModel{
		Id:        1,
		LeagueId:  1,
		Name:      "Test Tournament",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
		League: leagues.LeagueModel{
			Id:        1,
			Name:      "Test League",
			Alias:     "TL",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result := model.TableName()

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTournamentProvider_TableName(t *testing.T) {
	model := TournamentProvider{}
	assert.Equal(t, "tournament_providers", model.TableName())
}
