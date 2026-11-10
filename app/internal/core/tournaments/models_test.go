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

func TestTournamentModel_IsWithinPeriod(t *testing.T) {
	start := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		startAt  *time.Time
		endAt    *time.Time
		at       time.Time
		expected bool
	}{
		{name: "inside dates", startAt: &start, endAt: &end, at: start.AddDate(0, 1, 0), expected: true},
		{name: "exactly at start", startAt: &start, endAt: &end, at: start, expected: true},
		{name: "before start", startAt: &start, endAt: &end, at: start.Add(-time.Minute), expected: false},
		{name: "evening game of last day stored as next UTC day", startAt: &start, endAt: &end, at: end.Add(26 * time.Hour), expected: true},
		{name: "exactly at end of tolerance", startAt: &start, endAt: &end, at: end.Add(36 * time.Hour), expected: true},
		{name: "after end of tolerance", startAt: &start, endAt: &end, at: end.Add(36*time.Hour + time.Minute), expected: false},
		{name: "unknown start is not checked", endAt: &end, at: start.AddDate(-5, 0, 0), expected: true},
		{name: "unknown end is not checked", startAt: &start, at: end.AddDate(5, 0, 0), expected: true},
		{name: "known end is checked when start is unknown", endAt: &end, at: end.AddDate(0, 1, 0), expected: false},
		{name: "both dates unknown", at: end.AddDate(5, 0, 0), expected: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := TournamentModel{StartAt: tc.startAt, EndAt: tc.endAt}

			assert.Equal(t, tc.expected, model.IsWithinPeriod(tc.at))
		})
	}
}

func TestTournamentProvider_TableName(t *testing.T) {
	model := TournamentProvider{}
	assert.Equal(t, "tournament_providers", model.TableName())
}
