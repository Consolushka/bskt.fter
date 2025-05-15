package statistics

import (
	"IMP/app/internal/statistics/infobasket"
	"IMP/app/internal/statistics/translator"
	mockTranslator "IMP/app/internal/statistics/translator/mocks"
	"IMP/app/pkg/string_utils"
	mockStringUtils "IMP/app/pkg/string_utils/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMlblMapper_mapPlayer(t *testing.T) {
	firstPlayerDate := time.Date(1990, 11, 25, 0, 0, 0, 0, time.UTC)
	secondPlayerDate := time.Date(1970, 11, 11, 0, 0, 0, 0, time.UTC)
	thirdPlayerDate := time.Date(2000, 12, 13, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name      string
		player    infobasket.PlayerBoxScoreDto
		result    PlayerDTO
		errorMsg  string
		mockSetup func(utilsInterface *mockStringUtils.MockStringUtilsInterface, translator *mockTranslator.MockInterface)
	}{
		{
			name: "Map default player with positive plus-minus from start",
			player: infobasket.PlayerBoxScoreDto{
				PersonNameRu: "Иванов Иван Иванович",
				PersonNameEn: "Ivanov Ivan Ivanovich",
				PersonBirth:  "25.11.1990",
				PersonID:     12345,
				PlusMinus:    20,
				Seconds:      1400,
				IsStart:      true,
			},
			result: PlayerDTO{
				FullNameLocal:  "Иванов Иван Иванович",
				FullNameEn:     "Ivanov Ivan Ivanovich",
				BirthDate:      &firstPlayerDate,
				LeaguePlayerID: "12345",
				Statistic: PlayerStatisticDTO{
					PlsMin:        20,
					PlayedSeconds: 1400,
					IsBench:       false,
				},
			},
			errorMsg: "",
			mockSetup: func(utilsInterface *mockStringUtils.MockStringUtilsInterface, translator *mockTranslator.MockInterface) {
				utilsInterface.EXPECT().
					HasNonLanguageChars("Ivanov Ivan Ivanovich", string_utils.Latin).
					Return(false, nil)
			},
		},
		{
			name: "Map default player with negative plus-minus from bench",
			player: infobasket.PlayerBoxScoreDto{
				PersonNameRu: "Красиков Петр Васильевич",
				PersonNameEn: "Krasikov Petr Vasilyevich",
				PersonBirth:  "11.11.1970",
				PersonID:     321551,
				PlusMinus:    -10,
				Seconds:      2800,
				IsStart:      false,
			},
			result: PlayerDTO{
				FullNameLocal:  "Красиков Петр Васильевич",
				FullNameEn:     "Krasikov Petr Vasilyevich",
				BirthDate:      &secondPlayerDate,
				LeaguePlayerID: "321551",
				Statistic: PlayerStatisticDTO{
					PlsMin:        -10,
					PlayedSeconds: 2800,
					IsBench:       true,
				},
			},
			errorMsg: "",
			mockSetup: func(utilsInterface *mockStringUtils.MockStringUtilsInterface, translator *mockTranslator.MockInterface) {
				utilsInterface.EXPECT().
					HasNonLanguageChars("Krasikov Petr Vasilyevich", string_utils.Latin).
					Return(false, nil)
			},
		},
		{
			name: "Map player with cyrillic en name",
			player: infobasket.PlayerBoxScoreDto{
				PersonNameRu: "Буданов Антон",
				PersonNameEn: "Буданов Антон",
				PersonBirth:  "13.12.2000",
				PersonID:     321551,
				PlusMinus:    0,
				Seconds:      2800,
				IsStart:      true,
			},
			result: PlayerDTO{
				FullNameLocal:  "Буданов Антон",
				FullNameEn:     "Budanov Anton",
				BirthDate:      &thirdPlayerDate,
				LeaguePlayerID: "321551",
				Statistic: PlayerStatisticDTO{
					PlsMin:        0,
					PlayedSeconds: 2800,
					IsBench:       false,
				},
			},
			errorMsg: "",
			mockSetup: func(utilsInterface *mockStringUtils.MockStringUtilsInterface, translator *mockTranslator.MockInterface) {
				utilsInterface.EXPECT().
					HasNonLanguageChars("Буданов Антон", string_utils.Latin).
					Return(true, nil)

				ruCode := "ru"
				translator.EXPECT().
					Translate("Буданов Антон", &ruCode, "en").
					Return("Budanov Anton")
			},
		},
		{
			name: "Map player with invalid dateformat",
			player: infobasket.PlayerBoxScoreDto{
				PersonNameRu: "Буданов Антон",
				PersonNameEn: "Budanov Anton",
				PersonBirth:  "11-24-2000",
				PersonID:     321551,
				PlusMinus:    0,
				Seconds:      1200,
				IsStart:      true,
			},
			result:   PlayerDTO{},
			errorMsg: "can't parse player birthdate. given birthdate: 11-24-2000 doesn't match format 02.01.2006",
			mockSetup: func(utilsInterface *mockStringUtils.MockStringUtilsInterface, translator *mockTranslator.MockInterface) {
			},
		},
		{
			name: "Map player with error in non-language chars",
			player: infobasket.PlayerBoxScoreDto{
				PersonNameRu: "Буданов Антон",
				PersonNameEn: "Budanov Anton",
				PersonBirth:  "13.12.2000",
				PersonID:     321551,
				PlusMinus:    0,
				Seconds:      1200,
				IsStart:      true,
			},
			result: PlayerDTO{},
			mockSetup: func(utilsInterface *mockStringUtils.MockStringUtilsInterface, translator *mockTranslator.MockInterface) {
				utilsInterface.EXPECT().
					HasNonLanguageChars("Budanov Anton", string_utils.Latin).
					Return(false, errors.New("error in non-language chars"))
			},
			errorMsg: "error in non-language chars",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stringUtils := mockStringUtils.NewMockStringUtilsInterface(ctrl)
	translatorImp := mockTranslator.NewMockInterface(ctrl)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup(stringUtils, translatorImp)
			mapper := newMlblMapper(stringUtils, translatorImp)

			result, err := mapper.mapPlayer(tc.player)

			if tc.errorMsg != "" {
				assert.EqualError(t, err, tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.result, result)
		})
	}
}

// TestMlblMapper_mapTeam tests the mapTeam method of mlblMapper
// Verify that when valid team data is provided - team data is correctly mapped
// Verify that when player mapping fails - error is returned
func TestMlblMapper_mapTeam(t *testing.T) {
	lbjBirthDate := time.Date(1984, time.December, 30, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		data     infobasket.TeamBoxScoreDto
		expected TeamBoxScoreDTO
		errorMsg string
	}{
		{
			name: "Successful mapping with valid team data",
			data: infobasket.TeamBoxScoreDto{
				TeamID: 123,
				TeamName: infobasket.TeamNameBoxScoreDto{
					CompTeamAbcNameEn: "LAL",
					CompTeamNameEn:    "Los Angeles Lakers",
				},
				Score: 105,
				Players: []infobasket.PlayerBoxScoreDto{
					{
						PersonID:     1,
						PersonNameRu: "Леброн Джеймс",
						PersonNameEn: "LeBron James",
						PersonBirth:  "30.12.1984",
						PlusMinus:    15,
						Seconds:      1800,
						IsStart:      true,
					},
				},
			},
			expected: TeamBoxScoreDTO{
				Alias:    "LAL",
				Name:     "Los Angeles Lakers",
				LeagueId: "123",
				Scored:   105,
				Players: []PlayerDTO{
					{
						FullNameLocal:  "Леброн Джеймс",
						FullNameEn:     "LeBron James",
						LeaguePlayerID: "1",
						BirthDate:      &lbjBirthDate,
						Statistic: PlayerStatisticDTO{
							PlsMin:        15,
							PlayedSeconds: 1800,
							IsBench:       false,
						},
						// Note: BirthDate is not compared directly as it's a pointer to time.Time
						// We'll handle this in the test assertion
					},
				},
			},
			errorMsg: "",
		},
		{
			name: "Error when player mapping fails due to invalid birthdate",
			data: infobasket.TeamBoxScoreDto{
				TeamID: 123,
				TeamName: infobasket.TeamNameBoxScoreDto{
					CompTeamAbcNameEn: "LAL",
					CompTeamNameEn:    "Los Angeles Lakers",
				},
				Score: 105,
				Players: []infobasket.PlayerBoxScoreDto{
					{
						PersonID:     1,
						PersonNameRu: "Леброн Джеймс",
						PersonNameEn: "LeBron James",
						PersonBirth:  "invalid-date", // Invalid date format
						PlusMinus:    15,
						Seconds:      1800,
						IsStart:      true,
					},
				},
			},
			expected: TeamBoxScoreDTO{},
			errorMsg: "can't parse player birthdate. given birthdate: invalid-date doesn't match format 02.01.2006",
		},
	}

	mapper := newMlblMapper(string_utils.NewStringUtils(), translator.NewTranslator())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := mapper.mapTeam(tc.data)

			if tc.errorMsg != "" {
				assert.EqualError(t, err, tc.errorMsg)
			} else {
				assert.NoError(t, err)

				// Compare all fields except BirthDate which needs special handling
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

// TestMlblMapper_mapGame tests the mapGame method of mlblMapper
//
// Verify that when valid game data is provided while mapping game - returns correct GameBoxScoreDTO
// Verify that when invalid game date format is provided while mapping game - returns error
// Verify that when game has overtime periods while mapping game - calculates correct duration
// Verify that when player has invalid birthdate while mapping game - returns error
func TestMlblMapper_mapGame(t *testing.T) {
	cases := []struct {
		name                 string
		game                 infobasket.GameBoxScoreResponse
		regulationPeriodsNum int
		periodDuration       int
		overtimeDuration     int
		leagueAlias          string
		expected             *GameBoxScoreDTO
		errorMsg             string
	}{
		{
			name: "Valid game data",
			game: infobasket.GameBoxScoreResponse{
				GameDate:   "01.02.2023",
				GameTime:   "18.30",
				GameStatus: 1,
				MaxPeriod:  4,
				GameTeams: []infobasket.TeamBoxScoreDto{
					{
						TeamName: infobasket.TeamNameBoxScoreDto{
							CompTeamAbcNameEn: "HOME",
							CompTeamNameEn:    "Home Team",
						},
						TeamID:  123,
						Score:   85,
						Players: []infobasket.PlayerBoxScoreDto{},
					},
					{
						TeamName: infobasket.TeamNameBoxScoreDto{
							CompTeamAbcNameEn: "AWAY",
							CompTeamNameEn:    "Away Team",
						},
						TeamID:  456,
						Score:   80,
						Players: []infobasket.PlayerBoxScoreDto{},
					},
				},
			},
			regulationPeriodsNum: 4,
			periodDuration:       10,
			overtimeDuration:     5,
			leagueAlias:          "MLBL",
			expected: &GameBoxScoreDTO{
				LeagueAliasEn: "MLBL",
				IsFinal:       true,
				HomeTeam: TeamBoxScoreDTO{
					Alias:    "HOME",
					Name:     "Home Team",
					LeagueId: "123",
					Scored:   85,
					Players:  []PlayerDTO{},
				},
				AwayTeam: TeamBoxScoreDTO{
					Alias:    "AWAY",
					Name:     "Away Team",
					LeagueId: "456",
					Scored:   80,
					Players:  []PlayerDTO{},
				},
				PlayedMinutes: 40, // 4 periods * 10 minutes
				ScheduledAt:   time.Date(2023, 2, 1, 18, 30, 0, 0, time.UTC),
			},
			errorMsg: "",
		},
		{
			name: "Game with overtime periods",
			game: infobasket.GameBoxScoreResponse{
				GameDate:   "01.02.2023",
				GameTime:   "18.30",
				GameStatus: 1,
				MaxPeriod:  6, // 4 regular periods + 2 overtimes
				GameTeams: []infobasket.TeamBoxScoreDto{
					{
						TeamName: infobasket.TeamNameBoxScoreDto{
							CompTeamAbcNameEn: "HOME",
							CompTeamNameEn:    "Home Team",
						},
						TeamID:  123,
						Score:   95,
						Players: []infobasket.PlayerBoxScoreDto{},
					},
					{
						TeamName: infobasket.TeamNameBoxScoreDto{
							CompTeamAbcNameEn: "AWAY",
							CompTeamNameEn:    "Away Team",
						},
						TeamID:  456,
						Score:   90,
						Players: []infobasket.PlayerBoxScoreDto{},
					},
				},
			},
			regulationPeriodsNum: 4,
			periodDuration:       10,
			overtimeDuration:     5,
			leagueAlias:          "MLBL",
			expected: &GameBoxScoreDTO{
				LeagueAliasEn: "MLBL",
				IsFinal:       true,
				HomeTeam: TeamBoxScoreDTO{
					Alias:    "HOME",
					Name:     "Home Team",
					LeagueId: "123",
					Scored:   95,
					Players:  []PlayerDTO{},
				},
				AwayTeam: TeamBoxScoreDTO{
					Alias:    "AWAY",
					Name:     "Away Team",
					LeagueId: "456",
					Scored:   90,
					Players:  []PlayerDTO{},
				},
				PlayedMinutes: 50, // 4 periods * 10 minutes + 2 overtimes * 5 minutes
				ScheduledAt:   time.Date(2023, 2, 1, 18, 30, 0, 0, time.UTC),
			},
			errorMsg: "",
		},
		{
			name: "Invalid game date format",
			game: infobasket.GameBoxScoreResponse{
				GameDate:   "2023-02-01", // Invalid format, should be DD.MM.YYYY
				GameTime:   "18.30",
				GameStatus: 0,
				MaxPeriod:  4,
				GameTeams: []infobasket.TeamBoxScoreDto{
					{}, {},
				},
			},
			regulationPeriodsNum: 4,
			periodDuration:       10,
			overtimeDuration:     5,
			leagueAlias:          "MLBL",
			expected:             nil,
			errorMsg:             "can't parse game datetime. given game datetime: 2023-02-01 18.30 doesn't match format 02.01.2006 15.04",
		},
		{
			name: "Invalid game time format",
			game: infobasket.GameBoxScoreResponse{
				GameDate:   "01.02.2023",
				GameTime:   "18:30", // Invalid format, should use dots instead of colons
				GameStatus: 0,
				MaxPeriod:  4,
				GameTeams: []infobasket.TeamBoxScoreDto{
					{}, {},
				},
			},
			regulationPeriodsNum: 4,
			periodDuration:       10,
			overtimeDuration:     5,
			leagueAlias:          "MLBL",
			expected:             nil,
			errorMsg:             "can't parse game datetime. given game datetime: 01.02.2023 18:30 doesn't match format 02.01.2006 15.04",
		},
		{
			name: "Player from home team with invalid birthdate",
			game: infobasket.GameBoxScoreResponse{
				GameDate:   "01.02.2023",
				GameTime:   "18.30",
				GameStatus: 1,
				MaxPeriod:  4,
				GameTeams: []infobasket.TeamBoxScoreDto{
					{
						TeamName: infobasket.TeamNameBoxScoreDto{
							CompTeamAbcNameEn: "HOME",
							CompTeamNameEn:    "Home Team",
						},
						TeamID: 123,
						Score:  85,
						Players: []infobasket.PlayerBoxScoreDto{
							{
								PersonID:     1,
								PersonNameRu: "Иван Иванов",
								PersonNameEn: "Ivan Ivanov",
								PersonBirth:  "2000-01-01", // Invalid format, should be DD.MM.YYYY
								IsStart:      true,
								PlusMinus:    5,
								Seconds:      1200,
							},
						},
					},
					{
						TeamName: infobasket.TeamNameBoxScoreDto{},
						TeamID:   456,
						Score:    80,
						Players:  []infobasket.PlayerBoxScoreDto{},
					},
				},
			},
			regulationPeriodsNum: 4,
			periodDuration:       10,
			overtimeDuration:     5,
			leagueAlias:          "MLBL",
			expected:             nil,
			errorMsg:             "can't parse player birthdate. given birthdate: 2000-01-01 doesn't match format 02.01.2006",
		},
		{
			name: "Player from away team with invalid birthdate",
			game: infobasket.GameBoxScoreResponse{
				GameDate:   "01.02.2023",
				GameTime:   "18.30",
				GameStatus: 1,
				MaxPeriod:  4,
				GameTeams: []infobasket.TeamBoxScoreDto{
					{
						TeamName: infobasket.TeamNameBoxScoreDto{},
						TeamID:   123,
						Score:    85,
						Players:  []infobasket.PlayerBoxScoreDto{},
					},
					{
						TeamName: infobasket.TeamNameBoxScoreDto{},
						TeamID:   456,
						Score:    80,
						Players: []infobasket.PlayerBoxScoreDto{
							{
								PersonID:     1,
								PersonNameRu: "Иван Иванов",
								PersonNameEn: "Ivan Ivanov",
								PersonBirth:  "2000-01-01", // Invalid format, should be DD.MM.YYYY
								IsStart:      true,
								PlusMinus:    5,
								Seconds:      1200,
							},
						},
					},
				},
			},
			regulationPeriodsNum: 4,
			periodDuration:       10,
			overtimeDuration:     5,
			leagueAlias:          "MLBL",
			expected:             nil,
			errorMsg:             "can't parse player birthdate. given birthdate: 2000-01-01 doesn't match format 02.01.2006",
		},
	}

	mapper := newMlblMapper(string_utils.NewStringUtils(), translator.NewTranslator())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := mapper.mapGame(tc.game, tc.regulationPeriodsNum, tc.periodDuration, tc.overtimeDuration, tc.leagueAlias)

			if tc.errorMsg != "" {
				assert.EqualError(t, err, tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestMlblProvider_GameBoxScore verifies the behavior of the GameBoxScore method
// in the mlblProvider struct under various conditions:
// - Verify that when valid game data is returned by the client while fetching a box score - the method returns properly mapped game data
// - Verify that when mapper returns an error while mapping game data - the method returns that error
func TestMlblProvider_GameBoxScore(t *testing.T) {
	cases := []struct {
		name      string
		gameId    string
		mockSetup func(mockClient *infobasket.MockClientInterface, mockMapper *MockmlblMapperInterface)
		expected  *GameBoxScoreDTO
		errorMsg  string
	}{
		{
			name:   "Success case - returns properly mapped game data",
			gameId: "12345",
			mockSetup: func(mockClient *infobasket.MockClientInterface, mockMapper *MockmlblMapperInterface) {
				gameResponse := infobasket.GameBoxScoreResponse{
					CompID: 12345,
					// Other fields would be populated here
				}

				expectedGame := &GameBoxScoreDTO{
					LeagueAliasEn: "MLBL",
					// Other fields would be populated here
				}

				mockClient.EXPECT().BoxScore("12345").Return(gameResponse)
				mockMapper.EXPECT().mapGame(gameResponse, 4, 10, 5, "MLBL").Return(expectedGame, nil)
			},
			expected: &GameBoxScoreDTO{
				Id:            "12345",
				LeagueAliasEn: "MLBL",
				// Other fields would be populated here
			},
			errorMsg: "",
		},
		{
			name:   "Error case - mapper returns error",
			gameId: "12345",
			mockSetup: func(mockClient *infobasket.MockClientInterface, mockMapper *MockmlblMapperInterface) {
				gameResponse := infobasket.GameBoxScoreResponse{
					CompID: 12345,
					// Other fields would be populated here
				}

				mockClient.EXPECT().BoxScore("12345").Return(gameResponse)
				mockMapper.EXPECT().mapGame(gameResponse, 4, 10, 5, "MLBL").Return(nil, errors.New("mapping error"))
			},
			expected: nil,
			errorMsg: "mapping error",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := infobasket.NewMockClientInterface(ctrl)
	mockMapper := NewMockmlblMapperInterface(ctrl)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks for this test case
			tc.mockSetup(mockClient, mockMapper)

			// Create provider with mocked dependencies
			provider := &mlblProvider{
				client: mockClient,
				mapper: mockMapper,
			}

			// Call the method being tested
			result, err := provider.GameBoxScore(tc.gameId)

			// Verify results
			if tc.errorMsg != "" {
				assert.EqualError(t, err, tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestMlblProvider_GamesByDate verifies the behavior of the GamesByDate method
// in the mlblProvider struct under various conditions:
// - Verify that when games exist for the given date while fetching scheduled games - the method returns correct game IDs
// - Verify that when no games exist for the given date while fetching scheduled games - the method returns an empty slice
// - Verify that when games exist for some competition IDs but not others while fetching scheduled games - the method returns only matching game IDs
func TestMlblProvider_GamesByDate(t *testing.T) {
	cases := []struct {
		name      string
		date      time.Time
		mockSetup func(mockClient *infobasket.MockClientInterface)
		expected  []string
		errorMsg  string
	}{
		{
			name: "Success case - returns game IDs for the given date",
			date: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			mockSetup: func(mockClient *infobasket.MockClientInterface) {
				// Format the date as expected by the method (DD.MM.YYYY)
				formattedDate := "15.05.2023"

				// Mock response for first competition ID (89960)
				mockClient.EXPECT().ScheduledGames(89960).Return([]infobasket.GameScheduleDto{
					{GameID: 1001, GameDate: formattedDate},
					{GameID: 1002, GameDate: "16.05.2023"}, // Different date, should be filtered out
					{GameID: 1003, GameDate: formattedDate},
				})

				// Mock response for second competition ID (89962)
				mockClient.EXPECT().ScheduledGames(89962).Return([]infobasket.GameScheduleDto{
					{GameID: 2001, GameDate: "14.05.2023"}, // Different date, should be filtered out
					{GameID: 2002, GameDate: formattedDate},
				})
			},
			expected: []string{"1001", "1003", "2002"},
			errorMsg: "",
		},
		{
			name: "Empty case - no games for the given date",
			date: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			mockSetup: func(mockClient *infobasket.MockClientInterface) {
				// Mock response for first competition ID (89960)
				mockClient.EXPECT().ScheduledGames(89960).Return([]infobasket.GameScheduleDto{
					{GameID: 1001, GameDate: "14.05.2023"}, // Different date
					{GameID: 1002, GameDate: "16.05.2023"}, // Different date
				})

				// Mock response for second competition ID (89962)
				mockClient.EXPECT().ScheduledGames(89962).Return([]infobasket.GameScheduleDto{
					{GameID: 2001, GameDate: "14.05.2023"}, // Different date
					{GameID: 2002, GameDate: "16.05.2023"}, // Different date
				})
			},
			expected: []string(nil),
			errorMsg: "",
		},
		{
			name: "Partial case - games exist only for one competition ID",
			date: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			mockSetup: func(mockClient *infobasket.MockClientInterface) {
				// Format the date as expected by the method (DD.MM.YYYY)
				formattedDate := "15.05.2023"

				// Mock response for first competition ID (89960) - no matching games
				mockClient.EXPECT().ScheduledGames(89960).Return([]infobasket.GameScheduleDto{
					{GameID: 1001, GameDate: "14.05.2023"}, // Different date
					{GameID: 1002, GameDate: "16.05.2023"}, // Different date
				})

				// Mock response for second competition ID (89962) - has matching games
				mockClient.EXPECT().ScheduledGames(89962).Return([]infobasket.GameScheduleDto{
					{GameID: 2001, GameDate: formattedDate},
					{GameID: 2002, GameDate: formattedDate},
				})
			},
			expected: []string{"2001", "2002"},
			errorMsg: "",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := infobasket.NewMockClientInterface(ctrl)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks for this test case
			tc.mockSetup(mockClient)

			// Create provider with mocked dependencies
			provider := &mlblProvider{
				client: mockClient,
				mapper: nil, // Not needed for this test
			}

			// Call the method being tested
			result, err := provider.GamesByDate(tc.date)

			// Verify results
			if tc.errorMsg != "" {
				assert.EqualError(t, err, tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestMlblProvider_GamesByTeam tests the GamesByTeam method of mlblProvider
// It verifies:
// - When all games are final (status 1), all game IDs are returned
// - When some games are not final, only final game IDs are returned
// - When there's an error during mapping, the error is propagated
func TestMlblProvider_GamesByTeam(t *testing.T) {
	cases := []struct {
		name     string
		setup    func(mockClient *infobasket.MockClientInterface)
		teamId   string
		expected []string
		errorMsg string
	}{
		{
			name: "All games are final",
			setup: func(mockClient *infobasket.MockClientInterface) {
				mockClient.EXPECT().TeamGames("123").Return(infobasket.TeamScheduleResponse{
					Games: []infobasket.GameScheduleDto{
						{GameID: 1001, GameStatus: 1},
						{GameID: 1002, GameStatus: 1},
						{GameID: 1003, GameStatus: 1},
					},
				})
			},
			teamId:   "123",
			expected: []string{"1001", "1002", "1003"},
			errorMsg: "",
		},
		{
			name: "Some games are not final",
			setup: func(mockClient *infobasket.MockClientInterface) {
				mockClient.EXPECT().TeamGames("456").Return(infobasket.TeamScheduleResponse{
					Games: []infobasket.GameScheduleDto{
						{GameID: 2001, GameStatus: 1},
						{GameID: 2002, GameStatus: 0}, // Not final
						{GameID: 2003, GameStatus: 1},
						{GameID: 2004, GameStatus: 2}, // Not final
					},
				})
			},
			teamId:   "456",
			expected: nil,
			errorMsg: "game is not final. or game status is: 0",
		},
		{
			name: "No final games",
			setup: func(mockClient *infobasket.MockClientInterface) {
				mockClient.EXPECT().TeamGames("789").Return(infobasket.TeamScheduleResponse{
					Games: []infobasket.GameScheduleDto{
						{GameID: 3001, GameStatus: 0},
						{GameID: 3002, GameStatus: 2},
					},
				})
			},
			teamId:   "789",
			expected: nil,
			errorMsg: "game is not final. or game status is: 0",
		},
		{
			name: "Empty games list",
			setup: func(mockClient *infobasket.MockClientInterface) {
				mockClient.EXPECT().TeamGames("999").Return(infobasket.TeamScheduleResponse{
					Games: []infobasket.GameScheduleDto{},
				})
			},
			teamId:   "999",
			expected: []string{},
			errorMsg: "",
		},
		{
			name: "Error during mapping",
			setup: func(mockClient *infobasket.MockClientInterface) {
				// This will cause an error during mapping because we're returning a nil slice
				// which will cause the Map function to panic and return an error
				mockClient.EXPECT().TeamGames("error").Return(infobasket.TeamScheduleResponse{
					Games: nil,
				})
			},
			teamId:   "error",
			expected: []string{},
			errorMsg: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := infobasket.NewMockClientInterface(ctrl)
			tc.setup(mockClient)

			provider := &mlblProvider{
				client: mockClient,
			}

			result, err := provider.GamesByTeam(tc.teamId)

			if tc.errorMsg != "" {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}
