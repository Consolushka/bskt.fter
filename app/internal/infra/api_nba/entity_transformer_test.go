package api_nba

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEntityTransformer_Transform verifies the behavior of the Transform method
// in the EntityTransformer struct under various conditions:
// - Verify that basic game transformation works correctly
// - Verify that team data is properly mapped
// - Verify that date and time handling is correct
// - Verify that scores are properly assigned
func TestEntityTransformer_Transform(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name              string
		inputGame         GameEntity
		setupMockClient   func(mock *MockClientInterface)
		expectedGameTitle string
		expectedHomeTeam  string
		expectedAwayTeam  string
		expectedHomeScore int
		expectedAwayScore int
		expectError       error
	}{
		{
			name: "Transform basic game successfully",
			inputGame: GameEntity{
				Id:     12345,
				Season: 2023,
				Date: GameDateEntity{
					Start: time.Date(2023, 12, 25, 20, 0, 0, 0, time.UTC),
				},
				Teams: GameTeamsEntity{
					Home: TeamEntity{
						Id:       1,
						Name:     "Los Angeles Lakers",
						Nickname: "Lakers",
						Code:     "LAL",
					},
					Visitors: TeamEntity{
						Id:       2,
						Name:     "Boston Celtics",
						Nickname: "Celtics",
						Code:     "BOS",
					},
				},
				Scores: GameTeamsScoresEntity{
					Home: TeamScoresEntity{
						Points: 115,
					},
					Visitors: TeamScoresEntity{
						Points: 108,
					},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12345, 0, "").Return(PlayerStatisticResponse{
					Response: []PlayerStatisticEntity{
						{
							Player: PlayerStatisticsPlayerGeneralDataEntity{
								Id:        101,
								Firstname: "LeBron",
								Lastname:  "James",
							},
							Team:      TeamEntity{Id: 1},
							Min:       "35:24",
							PlusMinus: "12",
						},
						{
							Player: PlayerStatisticsPlayerGeneralDataEntity{
								Id:        102,
								Firstname: "Jayson",
								Lastname:  "Tatum",
							},
							Team:      TeamEntity{Id: 2},
							Min:       "38:15",
							PlusMinus: "-7",
						},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(101, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1984-12-30",
							},
						},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(102, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1990-01-01",
							},
						},
					},
				}, nil)
			},
			expectedGameTitle: "LAL - BOS",
			expectedHomeTeam:  "Lakers",
			expectedAwayTeam:  "Celtics",
			expectedHomeScore: 115,
			expectedAwayScore: 108,
			expectError:       nil,
		},
		{
			name: "Transform game with team name processing",
			inputGame: GameEntity{
				Id: 12346,
				Teams: GameTeamsEntity{
					Home: TeamEntity{
						Id:       3,
						Name:     "Golden State Warriors",
						Nickname: "Warriors",
						Code:     "GSW",
					},
					Visitors: TeamEntity{
						Id:       4,
						Name:     "Miami Heat",
						Nickname: "Heat",
						Code:     "MIA",
					},
				},
				Scores: GameTeamsScoresEntity{
					Home:     TeamScoresEntity{Points: 120},
					Visitors: TeamScoresEntity{Points: 95},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12346, 0, "").Return(PlayerStatisticResponse{}, nil)
			},
			expectedGameTitle: "GSW - MIA",
			expectedHomeTeam:  "Warriors",
			expectedAwayTeam:  "Heat",
			expectedHomeScore: 120,
			expectedAwayScore: 95,
			expectError:       nil,
		},
		{
			name: "Handle error from enrichGamePlayers",
			inputGame: GameEntity{
				Id:     12345,
				Season: 2023,
				Date: GameDateEntity{
					Start: time.Date(2023, 12, 25, 20, 0, 0, 0, time.UTC),
				},
				Teams: GameTeamsEntity{
					Home: TeamEntity{
						Id:       1,
						Name:     "Los Angeles Lakers",
						Nickname: "Lakers",
						Code:     "LAL",
					},
					Visitors: TeamEntity{
						Id:       2,
						Name:     "Boston Celtics",
						Nickname: "Celtics",
						Code:     "BOS",
					},
				},
				Scores: GameTeamsScoresEntity{
					Home: TeamScoresEntity{
						Points: 115,
					},
					Visitors: TeamScoresEntity{
						Points: 108,
					},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12345, 0, "").Return(PlayerStatisticResponse{}, errors.New("unexpected error"))
			},
			expectedGameTitle: "LAL - BOS",
			expectedHomeTeam:  "Lakers",
			expectedAwayTeam:  "Celtics",
			expectedHomeScore: 115,
			expectedAwayScore: 108,
			expectError:       errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := NewMockClientInterface(ctrl)
			transformer := EntityTransformer{client: mockClient}

			tc.setupMockClient(mockClient)

			result, err := transformer.Transform(tc.inputGame)

			if tc.expectError != nil {
				assert.Equal(t, tc.expectError, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedGameTitle, result.GameModel.Title)
			assert.Equal(t, tc.expectedHomeTeam, result.HomeTeamStat.TeamModel.Name)
			assert.Equal(t, tc.expectedAwayTeam, result.AwayTeamStat.TeamModel.Name)
			assert.Equal(t, tc.expectedHomeScore, result.HomeTeamStat.GameTeamStatModel.Score)
			assert.Equal(t, tc.expectedAwayScore, result.AwayTeamStat.GameTeamStatModel.Score)
		})
	}
}

// TestEntityTransformer_enrichGamePlayers verifies that player statistics
// are properly assigned to correct teams
func TestEntityTransformer_enrichGamePlayers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name                     string
		game                     GameEntity
		setupMockClient          func(mock *MockClientInterface)
		expectedHomePlayersCount int
		expectedAwayPlayersCount int
		expectError              error
	}{
		{
			name: "Players distributed correctly between teams",
			game: GameEntity{
				Id: 12345,
				Teams: GameTeamsEntity{
					Home:     TeamEntity{Id: 1},
					Visitors: TeamEntity{Id: 2},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12345, 0, "").Return(PlayerStatisticResponse{
					Response: []PlayerStatisticEntity{
						{
							Player:    PlayerStatisticsPlayerGeneralDataEntity{Id: 101},
							Team:      TeamEntity{Id: 1},
							Min:       "35:24",
							PlusMinus: "12",
						},
						{
							Player:    PlayerStatisticsPlayerGeneralDataEntity{Id: 102},
							Team:      TeamEntity{Id: 1},
							Min:       "28:15",
							PlusMinus: "8",
						},
						{
							Player:    PlayerStatisticsPlayerGeneralDataEntity{Id: 201},
							Team:      TeamEntity{Id: 2},
							Min:       "32:45",
							PlusMinus: "-5",
						},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(101, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{Birth: PlayerBirthEntity{Date: "1990-01-01"}},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(102, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{Birth: PlayerBirthEntity{Date: "1990-01-01"}},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(201, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{Birth: PlayerBirthEntity{Date: "1990-01-01"}},
					},
				}, nil)
			},
			expectedHomePlayersCount: 2,
			expectedAwayPlayersCount: 1,
			expectError:              nil,
		},
		{
			name: "Handle error while fetching players statistics",
			game: GameEntity{
				Id: 12345,
				Teams: GameTeamsEntity{
					Home:     TeamEntity{Id: 1},
					Visitors: TeamEntity{Id: 2},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12345, 0, "").Return(PlayerStatisticResponse{}, errors.New("unexpected error"))
			},
			expectedHomePlayersCount: 2,
			expectedAwayPlayersCount: 1,
			expectError:              errors.New("unexpected error"),
		},
		{
			name: "Skip player with error while fetching his bio",
			game: GameEntity{
				Id: 12345,
				Teams: GameTeamsEntity{
					Home:     TeamEntity{Id: 1},
					Visitors: TeamEntity{Id: 2},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12345, 0, "").Return(PlayerStatisticResponse{
					Response: []PlayerStatisticEntity{
						{
							Player:    PlayerStatisticsPlayerGeneralDataEntity{Id: 101},
							Team:      TeamEntity{Id: 1},
							Min:       "35:24",
							PlusMinus: "12",
						},
						{
							Player:    PlayerStatisticsPlayerGeneralDataEntity{Id: 102},
							Team:      TeamEntity{Id: 1},
							Min:       "28:15",
							PlusMinus: "8",
						},
						{
							Player:    PlayerStatisticsPlayerGeneralDataEntity{Id: 201},
							Team:      TeamEntity{Id: 2},
							Min:       "32:45",
							PlusMinus: "-5",
						},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(101, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{Birth: PlayerBirthEntity{Date: "1990-01-01"}},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(102, "", 0, 0, "", "").Return(PlayersResponse{}, errors.New("unexpected error"))

				mock.EXPECT().PlayerInfo(201, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{Birth: PlayerBirthEntity{Date: "1990-01-01"}},
					},
				}, nil)
			},
			expectedHomePlayersCount: 1,
			expectedAwayPlayersCount: 1,
			expectError:              nil,
		},
		{
			name: "Handle empty player statistics",
			game: GameEntity{
				Id: 12346,
				Teams: GameTeamsEntity{
					Home:     TeamEntity{Id: 1},
					Visitors: TeamEntity{Id: 2},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12346, 0, "").Return(PlayerStatisticResponse{}, nil)
			},
			expectedHomePlayersCount: 0,
			expectedAwayPlayersCount: 0,
			expectError:              nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := NewMockClientInterface(ctrl)
			transformer := EntityTransformer{client: mockClient}

			gameBusinessEntity := &games.GameStatEntity{
				HomeTeamStat: teams.TeamStatEntity{},
				AwayTeamStat: teams.TeamStatEntity{},
			}

			tc.setupMockClient(mockClient)

			err := transformer.enrichGamePlayers(tc.game, gameBusinessEntity)

			if tc.expectError != nil {
				assert.Equal(t, tc.expectError, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, gameBusinessEntity.HomeTeamStat.PlayerStats, tc.expectedHomePlayersCount)
			assert.Len(t, gameBusinessEntity.AwayTeamStat.PlayerStats, tc.expectedAwayPlayersCount)
		})
	}
}

// TestEntityTransformer_enrichPlayerStatistic verifies that individual player
// statistics are properly transformed and calculated
func TestEntityTransformer_enrichPlayerStatistic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name                  string
		inputPlayer           PlayerStatisticEntity
		setupMockClient       func(mock *MockClientInterface)
		expectedFullName      string
		expectedSecondsPlayed int
		expectedPlusMinus     int8
		expectedBirthDate     string
		expectError           error
	}{
		{
			name: "Transform player statistics correctly",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        101,
					Firstname: "LeBron",
					Lastname:  "James",
				},
				Min:       "35:24",
				PlusMinus: "12",
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(101, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1984-12-30",
							},
						},
					},
				}, nil)
			},
			expectedFullName:      "LeBron James",
			expectedSecondsPlayed: 35*60 + 24, // 2124 seconds
			expectedPlusMinus:     12,
			expectedBirthDate:     "1984-12-30",
			expectError:           nil,
		},
		{
			name: "Handle negative plus-minus",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        102,
					Firstname: "Jayson",
					Lastname:  "Tatum",
				},
				Min:       "38:15",
				PlusMinus: "-7",
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(102, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1998-03-03",
							},
						},
					},
				}, nil)
			},
			expectedFullName:      "Jayson Tatum",
			expectedSecondsPlayed: 38*60 + 15, // 2295 seconds
			expectedPlusMinus:     -7,
			expectedBirthDate:     "1998-03-03",
			expectError:           nil,
		},
		{
			name: "Handle invalid minutes format",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        103,
					Firstname: "Test",
					Lastname:  "Player",
				},
				Min:       "invalid:time",
				PlusMinus: "5",
			},
			expectError: errors.New("strconv.Atoi: parsing \"invalid\": invalid syntax"),
		},
		{
			name: "Handle invalid seconds format",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        103,
					Firstname: "Test",
					Lastname:  "Player",
				},
				Min:       "23:time",
				PlusMinus: "5",
			},
			expectError: errors.New("strconv.Atoi: parsing \"time\": invalid syntax"),
		},
		{
			name: "Handle invalid plus-minus format",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        104,
					Firstname: "Test",
					Lastname:  "Player",
				},
				Min:       "30:00",
				PlusMinus: "invalid",
			},
			expectError: errors.New("strconv.Atoi: parsing \"invalid\": invalid syntax"),
		},
		{
			name: "Handle error from player info",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        102,
					Firstname: "Jayson",
					Lastname:  "Tatum",
				},
				Min:       "38:15",
				PlusMinus: "-7",
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(102, "", 0, 0, "", "").Return(PlayersResponse{}, errors.New("unexpected error"))
			},
			expectedFullName:      "Jayson Tatum",
			expectedSecondsPlayed: 38*60 + 15, // 2295 seconds
			expectedPlusMinus:     -7,
			expectedBirthDate:     "1998-03-03",
			expectError:           errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := NewMockClientInterface(ctrl)
			if tc.setupMockClient != nil {
				tc.setupMockClient(mockClient)
			}
			transformer := EntityTransformer{client: mockClient}

			var playerBusinessEntity players.PlayerStatisticEntity
			err := transformer.enrichPlayerStatistic(tc.inputPlayer, &playerBusinessEntity)

			if tc.expectError != nil {
				assert.Equal(t, tc.expectError, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedFullName, playerBusinessEntity.PlayerModel.FullName)
			assert.Equal(t, tc.expectedSecondsPlayed, playerBusinessEntity.GameTeamPlayerStatModel.PlayedSeconds)
			assert.Equal(t, tc.expectedPlusMinus, playerBusinessEntity.GameTeamPlayerStatModel.PlsMin)

			expectedDate, _ := time.Parse("2006-01-02", tc.expectedBirthDate)
			assert.Equal(t, expectedDate, playerBusinessEntity.PlayerModel.BirthDate)
		})
	}
}

// TestEntityTransformer_enrichPlayerBio verifies that player biographical
// information is properly fetched and transformed
func TestEntityTransformer_enrichPlayerBio(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name              string
		playerId          int
		setupMockClient   func(mock *MockClientInterface)
		expectedBirthDate string
		expectError       error
	}{
		{
			name:     "Enrich player bio successfully",
			playerId: 101,
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(101, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1984-12-30",
							},
						},
					},
				}, nil)
			},
			expectedBirthDate: "1984-12-30",
			expectError:       nil,
		},
		{
			name:     "Handle error while fetching player bio",
			playerId: 101,
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(101, "", 0, 0, "", "").Return(PlayersResponse{}, errors.New("unexpected error"))
			},
			expectedBirthDate: "1984-12-30",
			expectError:       errors.New("unexpected error"),
		},
		{
			name:     "Handle invalid date format",
			playerId: 102,
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(102, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "invalid-date",
							},
						},
					},
				}, nil)
			},
			expectError: errors.New("parsing time \"invalid-date\" as \"2006-01-02\": cannot parse \"invalid-date\" as \"2006\""),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := NewMockClientInterface(ctrl)
			tc.setupMockClient(mockClient)
			transformer := EntityTransformer{client: mockClient}

			playerEntity := &players.PlayerStatisticEntity{}
			err := transformer.enrichPlayerBio(tc.playerId, playerEntity)

			if tc.expectError != nil {
				assert.Equal(t, tc.expectError, err)
				return
			}

			assert.NoError(t, err)
			expectedDate, _ := time.Parse("2006-01-02", tc.expectedBirthDate)
			assert.Equal(t, expectedDate, playerEntity.PlayerModel.BirthDate)
		})
	}
}

// TestNewEntityTransformer verifies that NewEntityTransformer returns
// a valid transformer instance with the provided client
func TestNewEntityTransformer(t *testing.T) {
	mockClient := new(MockClientInterface)
	transformer := NewEntityTransformer(mockClient)

	assert.NotNil(t, transformer.client)
	assert.Equal(t, mockClient, transformer.client)
}
