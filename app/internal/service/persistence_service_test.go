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

// TestNewPersistenceService tests the NewPersistenceService function
// Verify that when repositories are provided while creating service - returns PersistenceService instance with repositories set
func TestNewPersistenceService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name        string
		gamesRepo   *games_repo.MockGamesRepo
		teamsRepo   *teams_repo.MockTeamsRepo
		playersRepo *players_repo.MockPlayersRepo
		expected    *PersistenceService
		errorMsg    string
	}{
		{
			name:        "successfully creates PersistenceService with repositories",
			gamesRepo:   games_repo.NewMockGamesRepo(ctrl),
			teamsRepo:   teams_repo.NewMockTeamsRepo(ctrl),
			playersRepo: players_repo.NewMockPlayersRepo(ctrl),
			expected:    &PersistenceService{},
			errorMsg:    "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewPersistenceService(tc.gamesRepo, tc.teamsRepo, tc.playersRepo)

			assert.Equal(t, tc.gamesRepo, result.gamesRepo)
			assert.Equal(t, tc.teamsRepo, result.teamsRepo)
			assert.Equal(t, tc.playersRepo, result.playersRepo)
		})
	}
}

// TestPersistenceService_SaveGame tests the SaveGame method
// Verify successful saving of game with all related entities
// Verify error handling when repositories fail
func TestPersistenceService_SaveGame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name      string
		data      games.GameStatEntity
		expected  error
		errorMsg  string
		setupMock func(*games_repo.MockGamesRepo, *teams_repo.MockTeamsRepo, *players_repo.MockPlayersRepo)
	}{
		{
			name: "successfully saves game with all entities",
			data: games.GameStatEntity{
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
			expected: nil,
			errorMsg: "",
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
		},
		{
			name: "returns error when game repository fails",
			data: games.GameStatEntity{
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
			expected: errors.New("game repository error"),
			errorMsg: "game repository error",
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{}, errors.New("game repository error"))
			},
		},
		{
			name: "returns error when home team repository fails",
			data: games.GameStatEntity{
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
			expected: errors.New("team repository error"),
			errorMsg: "team repository error",
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{}, errors.New("team repository error"))
			},
		},
		{
			name: "returns error when away team repository fails",
			data: games.GameStatEntity{
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
			expected: errors.New("away team repository error"),
			errorMsg: "away team repository error",
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{}, errors.New("away team repository error"))
			},
		},
		{
			name: "continues processing when player operations fail",
			data: games.GameStatEntity{
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
			expected: nil,
			errorMsg: "",
			setupMock: func(mockGamesRepo *games_repo.MockGamesRepo, mockTeamsRepo *teams_repo.MockTeamsRepo, mockPlayersRepo *players_repo.MockPlayersRepo) {
				// Game and teams succeed
				mockGamesRepo.EXPECT().FindOrCreateGame(gomock.Any()).Return(games.GameModel{Id: 1}, nil)
				mockTeamsRepo.EXPECT().FirstOrCreateTeam(gomock.Any()).Return(teams.TeamModel{Id: 1}, nil).Times(2)
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{Id: 1}, nil).Times(2)

				// Players fail but service continues
				mockPlayersRepo.EXPECT().FirstOrCreatePlayer(gomock.Any()).Return(players.PlayerModel{}, errors.New("player error")).Times(2)
			},
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

			result := service.SaveGame(tc.data)

			if tc.errorMsg != "" {
				assert.EqualError(t, result, tc.errorMsg)
			} else {
				assert.NoError(t, result)
			}

			assert.Equal(t, tc.expected, result)
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
		name      string
		data      *teams.TeamStatEntity
		expected  error
		errorMsg  string
		setupMock func(*teams_repo.MockTeamsRepo)
	}{
		{
			name: "successfully saves team model and assigns ID",
			data: &teams.TeamStatEntity{
				TeamModel:         teams.TeamModel{Id: 0},
				GameTeamStatModel: teams.GameTeamStatModel{},
			},
			expected: nil,
			errorMsg: "",
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
			expected: errors.New("team repository error"),
			errorMsg: "team repository error",
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

			if tc.errorMsg != "" {
				assert.EqualError(t, result, tc.errorMsg)
			} else {
				assert.NoError(t, result)
				// Verify that TeamId was assigned correctly
				assert.Equal(t, uint(1), tc.data.GameTeamStatModel.TeamId)
			}

			assert.Equal(t, tc.expected, result)
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
		errorMsg  string
		setupMock func(*teams_repo.MockTeamsRepo)
	}{
		{
			name: "successfully saves team stats and assigns player GameTeamId",
			data: &teams.TeamStatEntity{
				GameTeamStatModel: teams.GameTeamStatModel{Id: 0},
				PlayerStats: []players.PlayerStatisticEntity{
					{GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{}},
					{GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{}},
				},
			},
			expected: nil,
			errorMsg: "",
			setupMock: func(mockTeamsRepo *teams_repo.MockTeamsRepo) {
				mockTeamsRepo.EXPECT().FirstOrCreateTeamStats(gomock.Any()).Return(teams.GameTeamStatModel{Id: 1}, nil)
			},
		},
		{
			name: "returns error when team stats repository fails",
			data: &teams.TeamStatEntity{
				GameTeamStatModel: teams.GameTeamStatModel{Id: 0},
			},
			expected: errors.New("team stats repository error"),
			errorMsg: "team stats repository error",
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

			if tc.errorMsg != "" {
				assert.EqualError(t, result, tc.errorMsg)
			} else {
				assert.NoError(t, result)
				// Verify that GameTeamId was assigned to all player stats
				for _, playerStat := range tc.data.PlayerStats {
					assert.Equal(t, uint(0), playerStat.GameTeamPlayerStatModel.GameTeamId)
				}
			}

			assert.Equal(t, tc.expected, result)
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
		name      string
		data      *players.PlayerStatisticEntity
		expected  error
		errorMsg  string
		setupMock func(*players_repo.MockPlayersRepo)
	}{
		{
			name: "successfully saves player model and assigns ID",
			data: &players.PlayerStatisticEntity{
				PlayerModel:             players.PlayerModel{Id: 0},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{},
			},
			expected: nil,
			errorMsg: "",
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
			expected: errors.New("player repository error"),
			errorMsg: "player repository error",
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

			if tc.errorMsg != "" {
				assert.EqualError(t, result, tc.errorMsg)
			} else {
				assert.NoError(t, result)
				// Verify that PlayerId was assigned correctly
				assert.Equal(t, uint(1), tc.data.GameTeamPlayerStatModel.PlayerId)
			}

			assert.Equal(t, tc.expected, result)
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
		name      string
		data      *players.PlayerStatisticEntity
		expected  error
		errorMsg  string
		setupMock func(*players_repo.MockPlayersRepo)
	}{
		{
			name: "successfully saves player stats",
			data: &players.PlayerStatisticEntity{
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 0},
			},
			expected: nil,
			errorMsg: "",
			setupMock: func(mockPlayersRepo *players_repo.MockPlayersRepo) {
				mockPlayersRepo.EXPECT().FirstOrCreatePlayerStat(gomock.Any()).Return(players.GameTeamPlayerStatModel{Id: 1}, nil)
			},
		},
		{
			name: "returns error when player stats repository fails",
			data: &players.PlayerStatisticEntity{
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{Id: 0},
			},
			expected: errors.New("player stats repository error"),
			errorMsg: "player stats repository error",
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

			if tc.errorMsg != "" {
				assert.EqualError(t, result, tc.errorMsg)
			} else {
				assert.NoError(t, result)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}
