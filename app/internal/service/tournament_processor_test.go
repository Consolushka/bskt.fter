package service

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/core/games"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewTournamentProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPersistence := NewMockPersistenceServiceInterface(ctrl)
	mockStatsProvider := stats_provider.NewMockStatsProvider(ctrl)

	processor := NewTournamentProcessor(mockStatsProvider, mockPersistence)

	assert.NotNil(t, processor)
	assert.Equal(t, mockPersistence, processor.persistenceService)
	assert.Equal(t, mockStatsProvider, processor.statsProvider)
}

func TestTournamentProcessor_ProcessTournamentGames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gameEntities := []games.GameStatEntity{
		{}, // Добавьте необходимые поля
		{},
	}

	cases := []struct {
		name       string
		setupMocks func(*MockPersistenceServiceInterface, *stats_provider.MockStatsProvider)
		expected   error
	}{
		{
			name: "successfully processes all games",
			setupMocks: func(mockPersistence *MockPersistenceServiceInterface, mockStats *stats_provider.MockStatsProvider) {
				mockStats.EXPECT().GetGamesStatsByDate(gomock.Any()).Return(gameEntities, nil)
				mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil).Times(2)
			},
			expected: nil,
		},
		{
			name: "returns error when stats provider fails",
			setupMocks: func(mockPersistence *MockPersistenceServiceInterface, mockStats *stats_provider.MockStatsProvider) {
				mockStats.EXPECT().GetGamesStatsByDate(gomock.Any()).Return(nil, errors.New("stats provider error"))
			},
			expected: errors.New("stats provider error"),
		},
		{
			name: "continues processing when some saves fail",
			setupMocks: func(mockPersistence *MockPersistenceServiceInterface, mockStats *stats_provider.MockStatsProvider) {
				mockStats.EXPECT().GetGamesStatsByDate(gomock.Any()).Return(gameEntities, nil)
				gomock.InOrder(
					mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(errors.New("save error")),
					mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil),
				)
			},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockPersistence := NewMockPersistenceServiceInterface(ctrl)
			mockStats := stats_provider.NewMockStatsProvider(ctrl)
			tc.setupMocks(mockPersistence, mockStats)

			processor := NewTournamentProcessor(mockStats, mockPersistence)
			result := processor.ProcessByPeriod()

			if tc.expected != nil {
				assert.EqualError(t, result, tc.expected.Error())
			} else {
				assert.NoError(t, result)
			}
		})
	}
}
