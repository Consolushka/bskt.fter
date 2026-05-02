package service

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/core/teams"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPersistenceService_SaveTeamModel_AliasLogic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTeamsRepo := teams_repo.NewMockTeamsRepo(ctrl)
	service := NewPersistenceService(
		games_repo.NewMockGamesRepo(ctrl),
		mockTeamsRepo,
		players_repo.NewMockPlayersRepo(ctrl),
	)

	tests := []struct {
		name          string
		incomingTeam  teams.TeamModel
		dbReturnTeam  teams.TeamModel
		expectUpdate  bool
		expectedAlias string
	}{
		{
			name: "Case 1: New team with provider alias - no update needed",
			incomingTeam: teams.TeamModel{
				Name:  "Lakers",
				Alias: "LAL",
			},
			dbReturnTeam: teams.TeamModel{
				Id:    1,
				Name:  "Lakers",
				Alias: "LAL",
			},
			expectUpdate:  false,
			expectedAlias: "LAL",
		},
		{
			name: "Case 2: New team without provider alias - no update needed (saved with generated alias)",
			incomingTeam: teams.TeamModel{
				Name:  "Lakers",
				Alias: "",
			},
			dbReturnTeam: teams.TeamModel{
				Id:    1,
				Name:  "Lakers",
				Alias: "Lak", // Repo generated this during FirstOrCreate
			},
			expectUpdate:  false,
			expectedAlias: "Lak",
		},
		{
			name: "Case 3: Existing team (no alias in DB) and no provider alias - update with generated",
			incomingTeam: teams.TeamModel{
				Name:  "Lakers",
				Alias: "",
			},
			dbReturnTeam: teams.TeamModel{
				Id:    1,
				Name:  "Lakers",
				Alias: "", // Found in DB without alias
			},
			expectUpdate:  true,
			expectedAlias: "Lak",
		},
		{
			name: "Case 4: Existing team (no alias in DB) and has provider alias - update with provider's",
			incomingTeam: teams.TeamModel{
				Name:  "Lakers",
				Alias: "LAL",
			},
			dbReturnTeam: teams.TeamModel{
				Id:    1,
				Name:  "Lakers",
				Alias: "",
			},
			expectUpdate:  true,
			expectedAlias: "LAL",
		},
		{
			name: "Case 5: Existing team (with alias in DB) and no provider alias - no update",
			incomingTeam: teams.TeamModel{
				Name:  "Lakers",
				Alias: "",
			},
			dbReturnTeam: teams.TeamModel{
				Id:    1,
				Name:  "Lakers",
				Alias: "LAL",
			},
			expectUpdate:  false,
			expectedAlias: "LAL",
		},
		{
			name: "Case 6: Existing team (with alias in DB) and provider alias differs - update with provider's",
			incomingTeam: teams.TeamModel{
				Name:  "Lakers",
				Alias: "LAK",
			},
			dbReturnTeam: teams.TeamModel{
				Id:    1,
				Name:  "Lakers",
				Alias: "LAL",
			},
			expectUpdate:  true,
			expectedAlias: "LAK",
		},
		{
			name: "Case 7: Existing team (with alias in DB) and provider alias is the same - no update",
			incomingTeam: teams.TeamModel{
				Name:  "Lakers",
				Alias: "LAL",
			},
			dbReturnTeam: teams.TeamModel{
				Id:    1,
				Name:  "Lakers",
				Alias: "LAL",
			},
			expectUpdate:  false,
			expectedAlias: "LAL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teamStats := &teams.TeamStatEntity{
				TeamModel: tt.incomingTeam,
			}

			mockTeamsRepo.EXPECT().FirstOrCreate(gomock.Any()).Return(tt.dbReturnTeam, nil)

			if tt.expectUpdate {
				mockTeamsRepo.EXPECT().UpdateAlias(tt.dbReturnTeam.Id, tt.expectedAlias).Return(nil)
			}

			err := service.saveTeamModel(teamStats)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedAlias, teamStats.TeamModel.Alias)
		})
	}
}
