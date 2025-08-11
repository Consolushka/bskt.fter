package service

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestNewPersistenceService verifies the behavior of the NewPersistenceService constructor
// in the PersistenceService under various conditions:
// - Verify that repositories are correctly assigned when creating service instance
// - Verify that the service instance is properly initialized with all dependencies
func TestNewPersistenceService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name       string
		setupRepos func() (*games_repo.MockGamesRepo, *teams_repo.MockTeamsRepo, *players_repo.MockPlayersRepo)
	}{
		{
			name: "Successfully creates PersistenceService with all repositories",
			setupRepos: func() (*games_repo.MockGamesRepo, *teams_repo.MockTeamsRepo, *players_repo.MockPlayersRepo) {
				return games_repo.NewMockGamesRepo(ctrl),
					teams_repo.NewMockTeamsRepo(ctrl),
					players_repo.NewMockPlayersRepo(ctrl)
			},
		},
		{
			name: "Creates PersistenceService with nil repositories",
			setupRepos: func() (*games_repo.MockGamesRepo, *teams_repo.MockTeamsRepo, *players_repo.MockPlayersRepo) {
				return nil, nil, nil
			},
		},
		{
			name: "Creates PersistenceService with mixed repository states",
			setupRepos: func() (*games_repo.MockGamesRepo, *teams_repo.MockTeamsRepo, *players_repo.MockPlayersRepo) {
				return games_repo.NewMockGamesRepo(ctrl),
					nil,
					players_repo.NewMockPlayersRepo(ctrl)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gamesRepo, teamsRepo, playersRepo := tc.setupRepos()

			result := NewPersistenceService(gamesRepo, teamsRepo, playersRepo)

			// Verify service instance is not nil
			assert.NotNil(t, result)

			// Verify repositories are correctly assigned
			assert.Equal(t, gamesRepo, result.gamesRepo)
			assert.Equal(t, teamsRepo, result.teamsRepo)
			assert.Equal(t, playersRepo, result.playersRepo)
		})
	}
}

// TestPersistenceService_SaveGame verifies the behavior of the SaveGame method
// in the PersistenceService under various conditions:
// - Verify successful saving of game with all related entities
// - Verify error handling when repositories fail
// - Verify that player operation errors don't stop processing
// - Verify that all repository calls are made in correct order
func TestPersistenceService_SaveGame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name        string
		inputGame   games.GameStatEntity
		setupMock   func(*games_repo.MockGamesRepo, *teams_repo.MockTeamsRepo, *players_repo.MockPlayersRepo)
		expectError error
	}{
		{
			name: "Successfully saves game with all entities",
			inputGame: games.GameStatEntity{
				GameModel: games.GameModel{Id: 1},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 1},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 1},
					PlayerStats: []players.PlayerStatisticEntity{
						{
							PlayerModel:             players.PlayerModel{Id: 1},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 1},
						},
					},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 2},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 2},
					PlayerStats: []players.PlayerStatisticEntity{
						{
							PlayerModel:             players.PlayerModel{Id: 2},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 2},
						},
					},
				},
			},
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				// Game repository calls
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 1}, nil)

				// Teams repository calls
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 1}, nil).Times(2)
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{Id: 1}, nil).Times(2)

				// Players repository calls
				mockPlayersRepo.EXPECT().FirstOrCreatePlayer(gomock.Any()).Return(players.PlayerModel{Id: 1}, nil).Times(2)
				mockPlayersRepo.EXPECT().FirstOrCreatePlayerStat(gomock.Any()).Return(players.GameTeamPlayerStatModel{Id: 1}, nil).Times(2)
			},
			expectError: nil,
		},
		{
			name: "Successfully saves game with multiple players per team",
			inputGame: games.GameStatEntity{
				GameModel: games.GameModel{Id: 2},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 10},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 10},
					PlayerStats: []players.PlayerStatisticEntity{
						{
							PlayerModel:             players.PlayerModel{Id: 101},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 101},
						},
						{
							PlayerModel:             players.PlayerModel{Id: 102},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 102},
						},
					},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 20},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 20},
					PlayerStats: []players.PlayerStatisticEntity{
						{
							PlayerModel:             players.PlayerModel{Id: 201},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 201},
						},
						{
							PlayerModel:             players.PlayerModel{Id: 202},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 202},
						},
						{
							PlayerModel:             players.PlayerModel{Id: 203},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 203},
						},
					},
				},
			},
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 2}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 10}, nil).Times(2)
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{Id: 10}, nil).Times(2)

				// 5 players total (2 home + 3 away)
				mockPlayersRepo.EXPECT().FirstOrCreatePlayer(gomock.Any()).Return(players.PlayerModel{Id: 1}, nil).Times(5)
				mockPlayersRepo.EXPECT().FirstOrCreatePlayerStat(gomock.Any()).Return(players.GameTeamPlayerStatModel{Id: 1}, nil).Times(5)
			},
			expectError: nil,
		},
		{
			name: "Successfully saves game with no players",
			inputGame: games.GameStatEntity{
				GameModel: games.GameModel{Id: 3},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 30},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 30},
					PlayerStats:       []players.PlayerStatisticEntity{},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 40},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 40},
					PlayerStats:       []players.PlayerStatisticEntity{},
				},
			},
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 3}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 30}, nil).Times(2)
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{Id: 30}, nil).Times(2)
				// No player repository calls expectedError
			},
			expectError: nil,
		},
		{
			name: "Returns error when game repository fails",
			inputGame: games.GameStatEntity{
				GameModel: games.GameModel{Id: 1},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 1},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 1},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 2},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 2},
				},
			},
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{}, errors.New("game repository error"))
			},
			expectError: errors.New("game repository error"),
		},
		{
			name: "Returns error when home team repository fails",
			inputGame: games.GameStatEntity{
				GameModel: games.GameModel{Id: 1},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 1},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 1},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 2},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 2},
				},
			},
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{}, errors.New("home team repository error"))
			},
			expectError: errors.New("home team repository error"),
		},
		{
			name: "Returns error when away team repository fails",
			inputGame: games.GameStatEntity{
				GameModel: games.GameModel{Id: 1},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 1},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 1},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 2},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 2},
				},
			},
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{}, errors.New("away team repository error"))
			},
			expectError: errors.New("away team repository error"),
		},
		{
			name: "Returns error when home team stats repository fails",
			inputGame: games.GameStatEntity{
				GameModel: games.GameModel{Id: 1},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 1},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 1},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 2},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 2},
				},
			},
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 1}, nil).Times(2)
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{}, errors.New("home team stats error"))
			},
			expectError: errors.New("home team stats error"),
		},
		{
			name: "Returns error when away team stats repository fails",
			inputGame: games.GameStatEntity{
				GameModel: games.GameModel{Id: 1},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 1},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 1},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 2},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 2},
				},
			},
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 1}, nil).Times(2)
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{}, errors.New("away team stats error"))
			},
			expectError: errors.New("away team stats error"),
		},
		{
			name: "Continues processing when player operations fail",
			inputGame: games.GameStatEntity{
				GameModel: games.GameModel{Id: 1},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 1},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 1},
					PlayerStats: []players.PlayerStatisticEntity{
						{
							PlayerModel:             players.PlayerModel{Id: 1},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 1},
						},
					},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 2},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 2},
					PlayerStats: []players.PlayerStatisticEntity{
						{
							PlayerModel:             players.PlayerModel{Id: 2},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 2},
						},
					},
				},
			},
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				// Game and teams succeed
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 1}, nil).Times(2)
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{Id: 1}, nil).Times(2)

				// Players fail but service continues
				mockPlayersRepo.EXPECT().FirstOrCreatePlayer(gomock.Any()).Return(players.PlayerModel{}, errors.New("player error")).Times(2)
			},
			expectError: nil,
		},
		{
			name: "Continues processing when player stats operations fail",
			inputGame: games.GameStatEntity{
				GameModel: games.GameModel{Id: 1},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 1},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 1},
					PlayerStats: []players.PlayerStatisticEntity{
						{
							PlayerModel:             players.PlayerModel{Id: 1},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 1},
						},
					},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel:         teams.TeamModel{Id: 2},
					GameTeamStatModel: teams.GameTeamStatModel{Id: 2},
					PlayerStats: []players.PlayerStatisticEntity{
						{
							PlayerModel:             players.PlayerModel{Id: 2},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 2},
						},
					},
				},
			},
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				// Game and teams succeed
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 1}, nil).Times(2)
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{Id: 1}, nil).Times(2)

				// Players succeed but stats fail
				mockPlayersRepo.EXPECT().FirstOrCreatePlayer(gomock.Any()).Return(players.PlayerModel{Id: 1}, nil).Times(2)
				mockPlayersRepo.EXPECT().FirstOrCreatePlayerStat(gomock.Any()).Return(players.GameTeamPlayerStatModel{}, errors.New("player stats error")).Times(2)
			},
			expectError: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockGamesRepo := games_repo.NewMockGamesRepo(ctrl)
			mockTeamsRepo := teams_repo.NewMockTeamsRepo(ctrl)
			mockPlayersRepo := players_repo.NewMockPlayersRepo(ctrl)

			tc.setupMock(mockGamesRepo, mockTeamsRepo, mockPlayersRepo)

			service := PersistenceService{
				gamesRepo:   mockGamesRepo,
				teamsRepo:   mockTeamsRepo,
				playersRepo: mockPlayersRepo,
			}

			err := service.SaveGame(tc.inputGame)

			assert.Equal(t, tc.expectError, err)
		})
	}
}

// TestPersistenceService_saveTeamModel tests the saveTeamModel method
// Verify successful team model saving and ID assignment
// Verify error handling when team repository fails
func TestPersistenceService_saveTeamModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name            string
		data            *teams.TeamStatEntity
		expectedError   error
		expectedModelId int
		setupMock       func(*teams_repo.MockTeamsRepo)
	}{
		{
			name: "successfully saves team model and assigns ID",
			data: &teams.TeamStatEntity{
				TeamModel:         teams.TeamModel{Id: 0},
				GameTeamStatModel: teams.GameTeamStatModel{},
			},
			expectedError:   nil,
			expectedModelId: 1,
			setupMock: func(mockTeamsRepo *teams_repo.MockTeamsRepo) {
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 1}, nil)
			},
		},
		{
			name: "returns error when team repository fails",
			data: &teams.TeamStatEntity{
				TeamModel:         teams.TeamModel{Id: 0},
				GameTeamStatModel: teams.GameTeamStatModel{},
			},
			expectedModelId: 0,
			expectedError:   errors.New("team repository error"),
			setupMock: func(mockTeamsRepo *teams_repo.MockTeamsRepo) {
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{}, errors.New("team repository error"))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockTeamsRepo := teams_repo.NewMockTeamsRepo(ctrl)
			tc.setupMock(mockTeamsRepo)

			service := PersistenceService{
				teamsRepo: mockTeamsRepo,
			}

			result := service.saveTeamModel(tc.data)

			if tc.expectedError != nil {
				assert.Equal(t, result, tc.expectedError)
				return
			}

			assert.NoError(t, result)
			assert.Equal(t, uint(tc.expectedModelId), tc.data.GameTeamStatModel.TeamId)
			assert.Equal(t, tc.expectedError, result)
		})
	}
}

// TestPersistenceService_saveTeamStatModel tests the saveTeamStatModel method
// Verify successful team stats saving and player stats ID assignment
// Verify error handling when team stats repository fails
func TestPersistenceService_saveTeamStatModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name      string
		data      *teams.TeamStatEntity
		expected  error
		setupMock func(*teams_repo.MockTeamsRepo)
	}{
		{
			name: "successfully saves team stats and assigns player GameTeamId",
			data: &teams.TeamStatEntity{
				TeamModel:         teams.TeamModel{},
				GameTeamStatModel: teams.GameTeamStatModel{Id: 0},
				PlayerStats: []players.PlayerStatisticEntity{
					{GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{}},
					{GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{}},
				},
			},
			expected: nil,
			setupMock: func(mockTeamsRepo *teams_repo.MockTeamsRepo) {
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{
					Id:     1,
					GameId: 2,
					TeamId: 3,
				}, nil)
			},
		},
		{
			name: "returns error when team stats repository fails",
			data: &teams.TeamStatEntity{
				GameTeamStatModel: teams.GameTeamStatModel{Id: 0},
			},
			expected: errors.New("team stats repository error"),
			setupMock: func(mockTeamsRepo *teams_repo.MockTeamsRepo) {
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{}, errors.New("team stats repository error"))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockTeamsRepo := teams_repo.NewMockTeamsRepo(ctrl)
			tc.setupMock(mockTeamsRepo)

			service := PersistenceService{
				teamsRepo: mockTeamsRepo,
			}

			result := service.saveTeamStatModel(tc.data)

			if tc.expected != nil {
				assert.Equal(t, result, tc.expected)
				return
			}

			assert.NoError(t, result)
			// Verify that GameTeamId was assigned to all player stats
			for _, playerStat := range tc.data.PlayerStats {
				assert.Equal(t, uint(2), playerStat.GameTeamPlayerStatModel.GameId)
				assert.Equal(t, uint(3), playerStat.GameTeamPlayerStatModel.TeamId)
			}
		})
	}
}

// TestPersistenceService_savePlayerModel tests the savePlayerModel method
// Verify successful player model saving and ID assignment
// Verify error handling when player repository fails
func TestPersistenceService_savePlayerModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name          string
		data          *players.PlayerStatisticEntity
		expectedError error
		setupMock     func(*players_repo.MockPlayersRepo)
	}{
		{
			name: "successfully saves player model and assigns ID",
			data: &players.PlayerStatisticEntity{
				PlayerModel:             players.PlayerModel{Id: 0},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{},
			},
			expectedError: nil,
			setupMock: func(mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockPlayersRepo.EXPECT().FirstOrCreatePlayer(gomock.Any()).Return(players.PlayerModel{Id: 1}, nil)
			},
		},
		{
			name: "returns error when player repository fails",
			data: &players.PlayerStatisticEntity{
				PlayerModel:             players.PlayerModel{Id: 0},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{},
			},
			expectedError: errors.New("player repository error"),
			setupMock: func(mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockPlayersRepo.EXPECT().FirstOrCreatePlayer(gomock.Any()).Return(players.PlayerModel{}, errors.New("player repository error"))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockPlayersRepo := players_repo.NewMockPlayersRepo(ctrl)
			tc.setupMock(mockPlayersRepo)

			service := PersistenceService{
				playersRepo: mockPlayersRepo,
			}

			result := service.savePlayerModel(tc.data)

			if tc.expectedError != nil {
				assert.Equal(t, result, tc.expectedError)
				return
			}

			assert.NoError(t, result)
			assert.Equal(t, uint(1), tc.data.GameTeamPlayerStatModel.PlayerId)
		})
	}
}

// TestPersistenceService_savePlayerStatModel tests the savePlayerStatModel method
// Verify successful player stats saving
// Verify error handling when player stats repository fails
func TestPersistenceService_savePlayerStatModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name          string
		data          *players.PlayerStatisticEntity
		expectedError error
		setupMock     func(*players_repo.MockPlayersRepo)
	}{
		{
			name: "successfully saves player stats",
			data: &players.PlayerStatisticEntity{
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 0},
			},
			expectedError: nil,
			setupMock: func(mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockPlayersRepo.EXPECT().FirstOrCreatePlayerStat(gomock.Any()).Return(players.GameTeamPlayerStatModel{Id: 1}, nil)
			},
		},
		{
			name: "returns error when player stats repository fails",
			data: &players.PlayerStatisticEntity{
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 0},
			},
			expectedError: errors.New("player stats repository error"),
			setupMock: func(mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockPlayersRepo.EXPECT().FirstOrCreatePlayerStat(gomock.Any()).Return(players.GameTeamPlayerStatModel{}, errors.New("player stats repository error"))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockPlayersRepo := players_repo.NewMockPlayersRepo(ctrl)
			tc.setupMock(mockPlayersRepo)

			service := PersistenceService{
				playersRepo: mockPlayersRepo,
			}

			result := service.savePlayerStatModel(tc.data)

			if tc.expectedError != nil {
				assert.Equal(t, result, tc.expectedError)
				return
			}

			assert.NoError(t, result)
		})
	}
}
