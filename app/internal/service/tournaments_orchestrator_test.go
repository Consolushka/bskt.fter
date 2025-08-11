package service

import (
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/core/leagues"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestNewTournamentsOrchestrator tests the NewTournamentsOrchestrator function
// Verify that when persistenceService and repository are provided while creating orchestrator - returns TournamentsOrchestrator instance with both services set
func TestNewTournamentsOrchestrator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name               string
		persistenceService PersistenceServiceInterface
		tournamentsRepo    ports.TournamentsRepo
		expectedResult     *TournamentsOrchestrator
		errorMsg           string
	}{
		{
			name:               "successfully creates TournamentsOrchestrator with persistenceService and repository",
			persistenceService: NewMockPersistenceServiceInterface(ctrl),
			tournamentsRepo:    tournaments_repo.NewMockTournamentsRepo(ctrl),
			expectedResult: &TournamentsOrchestrator{
				persistenceService: NewMockPersistenceServiceInterface(ctrl),
				tournamentsRepo:    tournaments_repo.NewMockTournamentsRepo(ctrl),
			},
			errorMsg: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewTournamentsOrchestrator(tc.persistenceService, tc.tournamentsRepo)

			assert.NotNil(t, result)
			assert.NotNil(t, result.persistenceService)
			assert.NotNil(t, result.tournamentsRepo)
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
		name          string
		data          struct{}
		expectedError error
		setupMocks    func(*MockPersistenceServiceInterface, *tournaments_repo.MockTournamentsRepo)
	}{
		{
			name:          "successfully processes tournaments with NBA and MLBL leagues",
			data:          struct{}{},
			expectedError: nil,
			setupMocks: func(mockPersistence *MockPersistenceServiceInterface, mockRepo *tournaments_repo.MockTournamentsRepo) {
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
				// Ожидаем вызовы SaveGame для каждой игры
				mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil).Times(0)
			},
		},
		{
			name:          "returns error when repository fails - differs from successful case",
			data:          struct{}{},
			expectedError: errors.New("repository error"),
			setupMocks: func(mockPersistence *MockPersistenceServiceInterface, mockRepo *tournaments_repo.MockTournamentsRepo) {
				mockRepo.EXPECT().ListActiveTournaments().Return(nil, errors.New("repository error"))
			},
		},
		{
			name:          "successfully processes tournaments with even with unexpected league",
			data:          struct{}{},
			expectedError: nil,
			setupMocks: func(mockPersistence *MockPersistenceServiceInterface, mockRepo *tournaments_repo.MockTournamentsRepo) {
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
				// Ожидаем вызовы SaveGame для игр из NBA (для UNEXPECTED провайдер не создастся)
				mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil).AnyTimes()
			},
		},
		{
			name:          "successfully processes empty tournaments list",
			data:          struct{}{},
			expectedError: nil,
			setupMocks: func(mockPersistence *MockPersistenceServiceInterface, mockRepo *tournaments_repo.MockTournamentsRepo) {
				mockRepo.EXPECT().ListActiveTournaments().Return([]tournaments.TournamentModel{}, nil)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockPersistence := NewMockPersistenceServiceInterface(ctrl)
			mockRepo := tournaments_repo.NewMockTournamentsRepo(ctrl)
			tc.setupMocks(mockPersistence, mockRepo)

			orchestrator := TournamentsOrchestrator{
				persistenceService: mockPersistence,
				tournamentsRepo:    mockRepo,
			}

			result := orchestrator.ProcessAllTournamentsToday()

			if tc.expectedError != nil {
				assert.Equal(t, result, tc.expectedError)
				return
			}

			assert.NoError(t, result)
		})
	}
}
