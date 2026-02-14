package api_nba

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

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
		name                   string
		inputGame              GameEntity
		setupMockClient        func(mock *MockClientInterface)
		expectedGameStatEntity games.GameStatEntity
		expectError            error
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
			expectedGameStatEntity: games.GameStatEntity{
				GameModel: games.GameModel{
					ScheduledAt: time.Date(2023, 12, 25, 20, 0, 0, 0, time.UTC),
					Title:       "LAL - BOS",
				},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel: teams.TeamModel{
						Name:     "Lakers",
						HomeTown: "Los Angeles", // "Los Angeles Lakers" -> remove "Lakers" -> trim space
					},
					GameTeamStatModel: teams.GameTeamStatModel{
						Score:     115,
						FinalDiff: 7,
					},
					PlayerStats: []players.PlayerStatisticEntity{
						{
							PlayerExternalId: "101",
							PlayerModel: players.PlayerModel{
								FullName:  "LeBron James",
								BirthDate: time.Date(1984, 12, 30, 0, 0, 0, 0, time.UTC),
							},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
								PlayedSeconds: 35*60 + 24, // 2124 seconds
								PlsMin:        12,
							},
						},
					},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel: teams.TeamModel{
						Name:     "Celtics",
						HomeTown: "Boston", // "Boston Celtics" -> remove "Celtics" -> trim space
					},
					GameTeamStatModel: teams.GameTeamStatModel{
						Score:     108,
						FinalDiff: -7,
					},
					PlayerStats: []players.PlayerStatisticEntity{
						{
							PlayerExternalId: "102",
							PlayerModel: players.PlayerModel{
								FullName:  "Jayson Tatum",
								BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
							},
							GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
								PlayedSeconds: 38*60 + 15, // 2295 seconds
								PlsMin:        -7,
							},
						},
					},
				},
			},
			expectError: nil,
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
			expectedGameStatEntity: games.GameStatEntity{
				GameModel: games.GameModel{
					Title: "GSW - MIA",
				},
				HomeTeamStat: teams.TeamStatEntity{
					TeamModel: teams.TeamModel{
						Name:     "Warriors",
						HomeTown: "Golden State", // "Golden State Warriors" -> remove "Warriors" -> trim space
					},
					GameTeamStatModel: teams.GameTeamStatModel{
						Score:     120,
						FinalDiff: 25,
					},
					PlayerStats: []players.PlayerStatisticEntity{},
				},
				AwayTeamStat: teams.TeamStatEntity{
					TeamModel: teams.TeamModel{
						Name:     "Heat",
						HomeTown: "Miami", // "Miami Heat" -> remove "Heat" -> trim space
					},
					GameTeamStatModel: teams.GameTeamStatModel{
						Score:     95,
						FinalDiff: -25,
					},
					PlayerStats: []players.PlayerStatisticEntity{},
				},
			},
			expectError: nil,
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
			expectedGameStatEntity: games.GameStatEntity{},
			expectError:            errors.New("unexpected error"),
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

			require.NoError(t, err)
			assert.Equal(t, tc.expectedGameStatEntity, result)
		})
	}
}

// TestEntityTransformer_TransformWithoutPlayers_Duration verifies
// duration calculation for regular periods and overtimes.
func TestEntityTransformer_TransformWithoutPlayers_Duration(t *testing.T) {
	cases := []struct {
		name             string
		periodsTotal     int
		expectedDuration int
	}{
		{
			name:             "regular_game_four_periods",
			periodsTotal:     4,
			expectedDuration: 48,
		},
		{
			name:             "five_periods_uses_regular_duration_for_first_five",
			periodsTotal:     5,
			expectedDuration: 60,
		},
		{
			name:             "overtime_adds_five_minutes_after_fifth_period",
			periodsTotal:     6,
			expectedDuration: 65,
		},
	}

	transformer := EntityTransformer{}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := transformer.TransformWithoutPlayers(GameEntity{
				Periods: GamePeriodsEntity{
					Total: tc.periodsTotal,
				},
				Teams: GameTeamsEntity{
					Home: TeamEntity{
						Name:     "Los Angeles Lakers",
						Nickname: "Lakers",
						Code:     "LAL",
					},
					Visitors: TeamEntity{
						Name:     "Boston Celtics",
						Nickname: "Celtics",
						Code:     "BOS",
					},
				},
			})

			assert.Equal(t, tc.expectedDuration, result.GameModel.Duration)
		})
	}
}

// TestEntityTransformer_enrichGamePlayers verifies the behavior of the enrichGamePlayers method
// in the EntityTransformer struct under various conditions:
// - Verify that player statistics are correctly fetched from the API
// - Verify that players are properly assigned to home and away teams
// - Verify that empty player statistics responses are handled correctly
// - Verify that API errors are properly propagated
// - Verify that individual player enrichment errors don't stop processing
// - Verify that players from other teams are ignored
// - Verify that the game business entity is correctly mutated with player data
func TestEntityTransformer_enrichGamePlayers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type expectedResult struct {
		homePlayerStats []players.PlayerStatisticEntity
		awayPlayerStats []players.PlayerStatisticEntity
		expectError     error
	}

	cases := []struct {
		name            string
		inputGame       GameEntity
		setupMockClient func(mock *MockClientInterface)
		expectedResult  expectedResult
	}{
		{
			name: "Successfully enriches game with player statistics",
			inputGame: GameEntity{
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
							Player: PlayerStatisticsPlayerGeneralDataEntity{
								Id:        101,
								Firstname: "LeBron",
								Lastname:  "James",
							},
							Team:      TeamEntity{Id: 1}, // Home team
							Min:       "35:24",
							PlusMinus: "12",
						},
						{
							Player: PlayerStatisticsPlayerGeneralDataEntity{
								Id:        102,
								Firstname: "Jayson",
								Lastname:  "Tatum",
							},
							Team:      TeamEntity{Id: 2}, // Away team
							Min:       "38:15",
							PlusMinus: "-7",
						},
						{
							Player: PlayerStatisticsPlayerGeneralDataEntity{
								Id:        103,
								Firstname: "Anthony",
								Lastname:  "Davis",
							},
							Team:      TeamEntity{Id: 1}, // Home team
							Min:       "32:45",
							PlusMinus: "8",
						},
					},
				}, nil)

				// Mock PlayerInfo calls for each player
				mock.EXPECT().PlayerInfo(101, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{Date: "1984-12-30"},
						},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(102, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{Date: "1998-03-03"},
						},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(103, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{Date: "1993-03-11"},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				homePlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "101",
						PlayerModel: players.PlayerModel{
							FullName:  "LeBron James",
							BirthDate: time.Date(1984, 12, 30, 0, 0, 0, 0, time.UTC),
						},
						GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
							PlayedSeconds: 35*60 + 24, // 2124 seconds
							PlsMin:        12,
						},
					},
					{
						PlayerExternalId: "103",
						PlayerModel: players.PlayerModel{
							FullName:  "Anthony Davis",
							BirthDate: time.Date(1993, 3, 11, 0, 0, 0, 0, time.UTC),
						},
						GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
							PlayedSeconds: 32*60 + 45, // 1965 seconds
							PlsMin:        8,
						},
					},
				},
				awayPlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "102",
						PlayerModel: players.PlayerModel{
							FullName:  "Jayson Tatum",
							BirthDate: time.Date(1998, 3, 3, 0, 0, 0, 0, time.UTC),
						},
						GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
							PlayedSeconds: 38*60 + 15, // 2295 seconds
							PlsMin:        -7,
						},
					},
				},
				expectError: nil,
			},
		},
		{
			name: "Handles empty player statistics response",
			inputGame: GameEntity{
				Id: 12346,
				Teams: GameTeamsEntity{
					Home:     TeamEntity{Id: 3},
					Visitors: TeamEntity{Id: 4},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12346, 0, "").Return(PlayerStatisticResponse{
					Response: []PlayerStatisticEntity{},
				}, nil)
			},
			expectedResult: expectedResult{
				homePlayerStats: []players.PlayerStatisticEntity{},
				awayPlayerStats: []players.PlayerStatisticEntity{},
				expectError:     nil,
			},
		},
		{
			name: "Handles error from PlayersStatistics API",
			inputGame: GameEntity{
				Id: 12347,
				Teams: GameTeamsEntity{
					Home:     TeamEntity{Id: 5},
					Visitors: TeamEntity{Id: 6},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12347, 0, "").Return(PlayerStatisticResponse{}, errors.New("API error"))
			},
			expectedResult: expectedResult{
				homePlayerStats: []players.PlayerStatisticEntity{},
				awayPlayerStats: []players.PlayerStatisticEntity{},
				expectError:     errors.New("API error"),
			},
		},
		{
			name: "Skips players with enrichment errors but continues processing",
			inputGame: GameEntity{
				Id: 12348,
				Teams: GameTeamsEntity{
					Home:     TeamEntity{Id: 7},
					Visitors: TeamEntity{Id: 8},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12348, 0, "").Return(PlayerStatisticResponse{
					Response: []PlayerStatisticEntity{
						{
							Player: PlayerStatisticsPlayerGeneralDataEntity{
								Id:        201,
								Firstname: "Good",
								Lastname:  "Player",
							},
							Team:      TeamEntity{Id: 7},
							Min:       "30:00",
							PlusMinus: "5",
						},
						{
							Player: PlayerStatisticsPlayerGeneralDataEntity{
								Id:        202,
								Firstname: "Bad",
								Lastname:  "Player",
							},
							Team:      TeamEntity{Id: 7},
							Min:       "invalid:time", // This will cause parsing error
							PlusMinus: "3",
						},
					},
				}, nil)

				// Only good player gets PlayerInfo call since bad player fails earlier
				mock.EXPECT().PlayerInfo(201, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{Date: "1995-01-01"},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				homePlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "201",
						PlayerModel: players.PlayerModel{
							FullName:  "Good Player",
							BirthDate: time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
							PlayedSeconds: 30 * 60, // 1800 seconds
							PlsMin:        5,
						},
					},
				},
				awayPlayerStats: []players.PlayerStatisticEntity{},
				expectError:     nil,
			},
		},
		{
			name: "Correctly assigns players to teams based on team ID",
			inputGame: GameEntity{
				Id: 12349,
				Teams: GameTeamsEntity{
					Home:     TeamEntity{Id: 10},
					Visitors: TeamEntity{Id: 20},
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayersStatistics(0, 12349, 0, "").Return(PlayerStatisticResponse{
					Response: []PlayerStatisticEntity{
						{
							Player: PlayerStatisticsPlayerGeneralDataEntity{
								Id:        301,
								Firstname: "Home",
								Lastname:  "Player1",
							},
							Team:      TeamEntity{Id: 10}, // Home team
							Min:       "25:30",
							PlusMinus: "2",
						},
						{
							Player: PlayerStatisticsPlayerGeneralDataEntity{
								Id:        302,
								Firstname: "Away",
								Lastname:  "Player1",
							},
							Team:      TeamEntity{Id: 20}, // Away team
							Min:       "28:45",
							PlusMinus: "-3",
						},
						{
							Player: PlayerStatisticsPlayerGeneralDataEntity{
								Id:        303,
								Firstname: "Other",
								Lastname:  "Player",
							},
							Team:      TeamEntity{Id: 99}, // Different team (should be ignored)
							Min:       "30:00",
							PlusMinus: "10",
						},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(301, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{Date: "1990-05-15"},
						},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(302, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{Date: "1992-08-20"},
						},
					},
				}, nil)

				mock.EXPECT().PlayerInfo(303, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{Date: "1988-12-01"},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				homePlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "301",
						PlayerModel: players.PlayerModel{
							FullName:  "Home Player1",
							BirthDate: time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
						},
						GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
							PlayedSeconds: 25*60 + 30, // 1530 seconds
							PlsMin:        2,
						},
					},
				},
				awayPlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "302",
						PlayerModel: players.PlayerModel{
							FullName:  "Away Player1",
							BirthDate: time.Date(1992, 8, 20, 0, 0, 0, 0, time.UTC),
						},
						GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
							PlayedSeconds: 28*60 + 45, // 1725 seconds
							PlsMin:        -3,
						},
					},
				},
				expectError: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := NewMockClientInterface(ctrl)
			transformer := EntityTransformer{client: mockClient}

			tc.setupMockClient(mockClient)

			// Create initial game business entity
			gameBusinessEntity := &games.GameStatEntity{
				HomeTeamStat: teams.TeamStatEntity{},
				AwayTeamStat: teams.TeamStatEntity{},
			}

			err := transformer.enrichGamePlayers(tc.inputGame, gameBusinessEntity)

			if tc.expectedResult.expectError != nil {
				assert.Equal(t, tc.expectedResult.expectError, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult.homePlayerStats, gameBusinessEntity.HomeTeamStat.PlayerStats)
			assert.Equal(t, tc.expectedResult.awayPlayerStats, gameBusinessEntity.AwayTeamStat.PlayerStats)
		})
	}
}

// TestEntityTransformer_enrichPlayerStatistic verifies the behavior of the enrichPlayerStatistic method
// in the EntityTransformer struct under various conditions:
// - Verify that player statistics are correctly parsed and transformed
// - Verify that biographical data is properly fetched and assigned
// - Verify that time format parsing works correctly (minutes:seconds)
// - Verify that plus/minus values are handled correctly (positive, negative, zero)
// - Verify that error handling works for invalid input formats
// - Verify that API errors are properly propagated
// - Verify that date parsing errors are handled appropriately
func TestEntityTransformer_enrichPlayerStatistic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type expectedResult struct {
		playerStatisticEntity players.PlayerStatisticEntity
		expectError           error
	}

	cases := []struct {
		name            string
		inputPlayer     PlayerStatisticEntity
		setupMockClient func(mock *MockClientInterface)
		expectedResult  expectedResult
	}{
		{
			name: "Successfully enriches player statistic with valid data",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        123,
					Firstname: "Stephen",
					Lastname:  "Curry",
				},
				Min:       "36:45",
				PlusMinus: "15",
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(123, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1988-03-14",
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{
					PlayerExternalId: "123",
					PlayerModel: players.PlayerModel{
						FullName:  "Stephen Curry",
						BirthDate: time.Date(1988, 3, 14, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 36*60 + 45, // 2205 seconds
						PlsMin:        15,
					},
				},
				expectError: nil,
			},
		},
		{
			name: "Handles negative plus/minus values",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        456,
					Firstname: "Russell",
					Lastname:  "Westbrook",
				},
				Min:       "28:12",
				PlusMinus: "-8",
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(456, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1988-11-12",
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{
					PlayerExternalId: "456",
					PlayerModel: players.PlayerModel{
						FullName:  "Russell Westbrook",
						BirthDate: time.Date(1988, 11, 12, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 28*60 + 12, // 1692 seconds
						PlsMin:        -8,
					},
				},
				expectError: nil,
			},
		},
		{
			name: "Handles zero playing time",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        789,
					Firstname: "Bench",
					Lastname:  "Player",
				},
				Min:       "0:00",
				PlusMinus: "0",
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(789, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1995-06-25",
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{
					PlayerExternalId: "789",
					PlayerModel: players.PlayerModel{
						FullName:  "Bench Player",
						BirthDate: time.Date(1995, 6, 25, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 0,
						PlsMin:        0,
					},
				},
				expectError: nil,
			},
		},
		{
			name: "Handles minutes without seconds",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        333,
					Firstname: "NoSeconds",
					Lastname:  "Player",
				},
				Min:       "25",
				PlusMinus: "4",
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(333, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1994-01-10",
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{
					PlayerExternalId: "333",
					PlayerModel: players.PlayerModel{
						FullName:  "NoSeconds Player",
						BirthDate: time.Date(1994, 1, 10, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 25 * 60,
						PlsMin:        4,
					},
				},
				expectError: nil,
			},
		},
		{
			name: "Returns error when minutes format is invalid",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        321,
					Firstname: "Invalid",
					Lastname:  "Minutes",
				},
				Min:       "invalid:format",
				PlusMinus: "5",
			},
			setupMockClient: func(mock *MockClientInterface) {
				// No PlayerInfo call expected since parsing should fail first
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{},
				expectError:           errors.New(`strconv.Atoi: parsing "invalid": invalid syntax`),
			},
		},
		{
			name: "Returns error when seconds format is invalid",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        654,
					Firstname: "Invalid",
					Lastname:  "Seconds",
				},
				Min:       "25:invalid",
				PlusMinus: "3",
			},
			setupMockClient: func(mock *MockClientInterface) {
				// No PlayerInfo call expected since parsing should fail first
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{},
				expectError:           errors.New(`strconv.Atoi: parsing "invalid": invalid syntax`),
			},
		},
		{
			name: "Returns error when plus/minus format is invalid",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        987,
					Firstname: "Invalid",
					Lastname:  "PlusMinus",
				},
				Min:       "30:15",
				PlusMinus: "not_a_number",
			},
			setupMockClient: func(mock *MockClientInterface) {
				// No PlayerInfo call expected since parsing should fail first
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{},
				expectError:           errors.New(`strconv.Atoi: parsing "not_a_number": invalid syntax`),
			},
		},
		{
			name: "Returns error when PlayerInfo API fails",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        555,
					Firstname: "API",
					Lastname:  "Error",
				},
				Min:       "25:30",
				PlusMinus: "7",
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(555, "", 0, 0, "", "").Return(PlayersResponse{}, errors.New("API connection failed"))
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{},
				expectError:           errors.New("API connection failed"),
			},
		},
		{
			name: "Returns error when birth date parsing fails",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        777,
					Firstname: "Invalid",
					Lastname:  "BirthDate",
				},
				Min:       "20:00",
				PlusMinus: "2",
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(777, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "invalid-date-format",
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{},
				expectError:           errors.New(`parsing time "invalid-date-format" as "2006-01-02": cannot parse "invalid-date-format" as "2006"`),
			},
		},
		{
			name: "Handles single digit minutes and seconds",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        111,
					Firstname: "Short",
					Lastname:  "Time",
				},
				Min:       "5:09",
				PlusMinus: "1",
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(111, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "2000-01-01",
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{
					PlayerExternalId: "111",
					PlayerModel: players.PlayerModel{
						FullName:  "Short Time",
						BirthDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 5*60 + 9, // 309 seconds
						PlsMin:        1,
					},
				},
				expectError: nil,
			},
		},
		{
			name: "Handles large plus/minus values within int8 range",
			inputPlayer: PlayerStatisticEntity{
				Player: PlayerStatisticsPlayerGeneralDataEntity{
					Id:        222,
					Firstname: "High",
					Lastname:  "Impact",
				},
				Min:       "42:30",
				PlusMinus: "127", // Max int8 value
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(222, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1985-12-25",
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerStatisticEntity: players.PlayerStatisticEntity{
					PlayerExternalId: "222",
					PlayerModel: players.PlayerModel{
						FullName:  "High Impact",
						BirthDate: time.Date(1985, 12, 25, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 42*60 + 30, // 2550 seconds
						PlsMin:        127,
					},
				},
				expectError: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := NewMockClientInterface(ctrl)
			transformer := EntityTransformer{client: mockClient}

			tc.setupMockClient(mockClient)

			var actualPlayerEntity players.PlayerStatisticEntity
			err := transformer.enrichPlayerStatistic(tc.inputPlayer, &actualPlayerEntity)

			if tc.expectedResult.expectError != nil {
				assert.Equal(t, tc.expectedResult.expectError.Error(), err.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult.playerStatisticEntity, actualPlayerEntity)
		})
	}
}

// TestEntityTransformer_enrichPlayerBio verifies the behavior of the enrichPlayerBio method
// in the EntityTransformer struct under various conditions:
// - Verify that player biographical data is correctly fetched and assigned
// - Verify that birth date parsing works correctly
// - Verify that API errors are properly handled
// - Verify that invalid date formats are handled appropriately
// - Verify that empty player responses are handled correctly
func TestEntityTransformer_enrichPlayerBio(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type expectedResult struct {
		playerEntity players.PlayerStatisticEntity
		expectError  error
	}

	cases := []struct {
		name            string
		playerId        int
		initialPlayer   players.PlayerStatisticEntity
		setupMockClient func(mock *MockClientInterface)
		expectedResult  expectedResult
	}{
		{
			name:     "Successfully enriches player with biographical data",
			playerId: 123,
			initialPlayer: players.PlayerStatisticEntity{
				PlayerModel: players.PlayerModel{
					FullName: "LeBron James",
				},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
					PlayedSeconds: 2100,
					PlsMin:        15,
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(123, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1984-12-30",
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerEntity: players.PlayerStatisticEntity{
					PlayerModel: players.PlayerModel{
						FullName:  "LeBron James",
						BirthDate: time.Date(1984, 12, 30, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 2100,
						PlsMin:        15,
					},
				},
				expectError: nil,
			},
		},
		{
			name:     "Handles leap year birth dates correctly",
			playerId: 456,
			initialPlayer: players.PlayerStatisticEntity{
				PlayerModel: players.PlayerModel{
					FullName: "Leap Year Player",
				},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
					PlayedSeconds: 1800,
					PlsMin:        -5,
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(456, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "2000-02-29", // Leap year date
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerEntity: players.PlayerStatisticEntity{
					PlayerModel: players.PlayerModel{
						FullName:  "Leap Year Player",
						BirthDate: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 1800,
						PlsMin:        -5,
					},
				},
				expectError: nil,
			},
		},
		{
			name:     "Preserves existing player data when enriching",
			playerId: 789,
			initialPlayer: players.PlayerStatisticEntity{
				PlayerModel: players.PlayerModel{
					FullName: "Existing Player",
				},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
					PlayedSeconds: 2400,
					PlsMin:        8,
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(789, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1992-07-15",
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerEntity: players.PlayerStatisticEntity{
					PlayerModel: players.PlayerModel{
						FullName:  "Existing Player",
						BirthDate: time.Date(1992, 7, 15, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 2400,
						PlsMin:        8,
					},
				},
				expectError: nil,
			},
		},
		{
			name:     "Returns error when PlayerInfo API fails",
			playerId: 404,
			initialPlayer: players.PlayerStatisticEntity{
				PlayerModel: players.PlayerModel{
					FullName: "API Error Player",
				},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
					PlayedSeconds: 1200,
					PlsMin:        0,
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(404, "", 0, 0, "", "").Return(PlayersResponse{}, errors.New("player not found"))
			},
			expectedResult: expectedResult{
				playerEntity: players.PlayerStatisticEntity{
					PlayerModel: players.PlayerModel{
						FullName: "API Error Player",
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 1200,
						PlsMin:        0,
					},
				},
				expectError: errors.New("player not found"),
			},
		},
		{
			name:     "Returns error when PlayerInfo response is empty",
			playerId: 405,
			initialPlayer: players.PlayerStatisticEntity{
				PlayerModel: players.PlayerModel{
					FullName: "Empty Response Player",
				},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
					PlayedSeconds: 1000,
					PlsMin:        1,
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(405, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{},
				}, nil)
			},
			expectedResult: expectedResult{
				playerEntity: players.PlayerStatisticEntity{
					PlayerModel: players.PlayerModel{
						FullName: "Empty Response Player",
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 1000,
						PlsMin:        1,
					},
				},
				expectError: errors.New("empty player info response"),
			},
		},
		{
			name:     "Returns error when birth date format is invalid",
			playerId: 555,
			initialPlayer: players.PlayerStatisticEntity{
				PlayerModel: players.PlayerModel{
					FullName: "Invalid Date Player",
				},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
					PlayedSeconds: 3000,
					PlsMin:        12,
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(555, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "invalid-date-format",
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerEntity: players.PlayerStatisticEntity{
					PlayerModel: players.PlayerModel{
						FullName: "Invalid Date Player",
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 3000,
						PlsMin:        12,
					},
				},
				expectError: errors.New(`parsing time "invalid-date-format" as "2006-01-02": cannot parse "invalid-date-format" as "2006"`),
			},
		},
		{
			name:     "Returns error when birth date has wrong format pattern",
			playerId: 666,
			initialPlayer: players.PlayerStatisticEntity{
				PlayerModel: players.PlayerModel{
					FullName: "Wrong Format Player",
				},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
					PlayedSeconds: 1500,
					PlsMin:        -3,
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(666, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "12-30-1984", // Wrong format (MM-DD-YYYY instead of YYYY-MM-DD)
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerEntity: players.PlayerStatisticEntity{
					PlayerModel: players.PlayerModel{
						FullName: "Wrong Format Player",
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 1500,
						PlsMin:        -3,
					},
				},
				expectError: errors.New(`parsing time "12-30-1984" as "2006-01-02": cannot parse "12-30-1984" as "2006"`),
			},
		},
		{
			name:     "Handles edge case birth dates (very old and very recent)",
			playerId: 888,
			initialPlayer: players.PlayerStatisticEntity{
				PlayerModel: players.PlayerModel{
					FullName: "Old Player",
				},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
					PlayedSeconds: 600,
					PlsMin:        1,
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(888, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "1950-01-01", // Very old date
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerEntity: players.PlayerStatisticEntity{
					PlayerModel: players.PlayerModel{
						FullName:  "Old Player",
						BirthDate: time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 600,
						PlsMin:        1,
					},
				},
				expectError: nil,
			},
		},
		{
			name:     "Handles recent birth dates",
			playerId: 999,
			initialPlayer: players.PlayerStatisticEntity{
				PlayerModel: players.PlayerModel{
					FullName: "Young Player",
				},
				GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
					PlayedSeconds: 2700,
					PlsMin:        20,
				},
			},
			setupMockClient: func(mock *MockClientInterface) {
				mock.EXPECT().PlayerInfo(999, "", 0, 0, "", "").Return(PlayersResponse{
					Response: []PlayerEntity{
						{
							Birth: PlayerBirthEntity{
								Date: "2005-12-31", // Recent date
							},
						},
					},
				}, nil)
			},
			expectedResult: expectedResult{
				playerEntity: players.PlayerStatisticEntity{
					PlayerModel: players.PlayerModel{
						FullName:  "Young Player",
						BirthDate: time.Date(2005, 12, 31, 0, 0, 0, 0, time.UTC),
					},
					GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
						PlayedSeconds: 2700,
						PlsMin:        20,
					},
				},
				expectError: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := NewMockClientInterface(ctrl)
			transformer := EntityTransformer{client: mockClient}

			tc.setupMockClient(mockClient)

			// Create a copy of initial player to avoid mutation
			actualPlayerEntity := tc.initialPlayer
			err := transformer.enrichPlayerBio(tc.playerId, &actualPlayerEntity)

			if tc.expectedResult.expectError != nil {
				assert.Equal(t, tc.expectedResult.expectError.Error(), err.Error())
				// Verify that the original data is preserved when there's an error
				assert.Equal(t, tc.expectedResult.playerEntity, actualPlayerEntity)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult.playerEntity, actualPlayerEntity)
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
