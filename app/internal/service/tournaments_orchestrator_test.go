package service

import (
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/core/leagues"
	"IMP/app/internal/core/tournaments"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestNewTournamentsOrchestrator tests the NewTournamentsOrchestrator function
// Verify that when repository is provided while creating orchestrator - returns TournamentsOrchestrator instance with repository set
func TestNewTournamentsOrchestrator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name     string
		data     *tournaments_repo.MockTournamentsRepo
		expected *TournamentsOrchestrator
		errorMsg string
	}{
		{
			name:     "successfully creates TournamentsOrchestrator with repository",
			data:     tournaments_repo.NewMockTournamentsRepo(ctrl),
			expected: &TournamentsOrchestrator{repo: tournaments_repo.NewMockTournamentsRepo(ctrl)},
			errorMsg: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result := NewTournamentsOrchestrator(tc.data)

			assert.Equal(t, tc.expected.repo, result.repo)
		})
	}
}

// TestTournamentsOrchestrator_ProcessAllTournamentsToday tests the ProcessAllTournamentsToday method
// Verify that when repository returns tournaments successfully while processing tournaments - returns no error
// Verify that when repository returns error while fetching tournaments - returns the same error
func TestTournamentsOrchestrator_ProcessAllTournamentsToday(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name      string
		data      struct{}
		expected  error
		errorMsg  string
		setupMock func(*tournaments_repo.MockTournamentsRepo)
	}{
		{
			name:     "successfully processes tournaments with NBA and MLBL leagues",
			data:     struct{}{},
			expected: nil,
			errorMsg: "",
			setupMock: func(mockRepo *tournaments_repo.MockTournamentsRepo) {
				mockRepo.EXPECT().ListActiveTournaments().Return([]tournaments.TournamentModel{
					{
						League: leagues.LeagueModel{
							Alias: "NBA",
						},
					},
					{
						League: leagues.LeagueModel{
							Alias: "MLBL",
						},
					},
				}, nil)
			},
		},
		{
			name:     "returns error when repository fails - differs from successful case",
			data:     struct{}{},
			expected: errors.New("repository error"),
			errorMsg: "repository error",
			setupMock: func(mockRepo *tournaments_repo.MockTournamentsRepo) {
				mockRepo.EXPECT().ListActiveTournaments().Return(nil, errors.New("repository error"))
			},
		},
		{
			name:     "successfully processes tournaments with even with unexpected league",
			data:     struct{}{},
			expected: nil,
			errorMsg: "",
			setupMock: func(mockRepo *tournaments_repo.MockTournamentsRepo) {
				mockRepo.EXPECT().ListActiveTournaments().Return([]tournaments.TournamentModel{
					{
						League: leagues.LeagueModel{
							Alias: "NBA",
						},
					},
					{
						League: leagues.LeagueModel{
							Alias: "UNEXPECTED",
						},
					},
				}, nil)
			},
		},
		{
			name:     "successfully processes empty tournaments list",
			data:     struct{}{},
			expected: nil,
			errorMsg: "",
			setupMock: func(mockRepo *tournaments_repo.MockTournamentsRepo) {
				mockRepo.EXPECT().ListActiveTournaments().Return([]tournaments.TournamentModel{}, nil)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := tournaments_repo.NewMockTournamentsRepo(ctrl)
			tc.setupMock(mockRepo)

			orchestrator := TournamentsOrchestrator{repo: mockRepo}

			result := orchestrator.ProcessAllTournamentsToday()

			if tc.errorMsg != "" {
				assert.EqualError(t, result, tc.errorMsg)
			} else {
				assert.NoError(t, result)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}
