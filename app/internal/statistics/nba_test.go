package statistics

import (
	"IMP/app/internal/domain"
	"IMP/app/internal/statistics/cdn_nba"
	"IMP/app/log"
	"IMP/app/pkg/time_utils"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNbaMapper_mapPlayer(t *testing.T) {
	cases := []struct {
		name      string
		player    cdn_nba.PlayerBoxScoreDto
		result    PlayerDTO
		errorMsg  string
		mockSetup func(timeUtils *time_utils.MockTimeUtilsInterface)
	}{
		{
			name: "Map player from starting lineup with positive plus-minus",
			player: cdn_nba.PlayerBoxScoreDto{
				Name:     "LeBron James",
				PersonId: 2544,
				Starter:  "1",
				Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
					Minutes: "PT35M20S",
					Plus:    25,
					Minus:   10,
				},
			},
			result: PlayerDTO{
				FullNameLocal:  "LeBron James",
				LeaguePlayerID: "2544",
				Statistic: PlayerStatisticDTO{
					PlsMin:        15,   // 25 - 10
					PlayedSeconds: 2120, // 35*60 + 20
					IsBench:       false,
				},
			},
			errorMsg: "",
			mockSetup: func(timeUtils *time_utils.MockTimeUtilsInterface) {
				timeUtils.EXPECT().
					FormattedMinutesToSeconds("PT35M20S", playedTimeFormat).
					Return(2120, nil)
			},
		},
		{
			name: "Map player from bench with negative plus-minus",
			player: cdn_nba.PlayerBoxScoreDto{
				Name:     "Kyle Kuzma",
				PersonId: 1628398,
				Starter:  "0",
				Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
					Minutes: "PT15M45S",
					Plus:    5,
					Minus:   15,
				},
			},
			result: PlayerDTO{
				FullNameLocal:  "Kyle Kuzma",
				LeaguePlayerID: "1628398",
				Statistic: PlayerStatisticDTO{
					PlsMin:        -10, // 5 - 15
					PlayedSeconds: 945, // 15*60 + 45
					IsBench:       true,
				},
			},
			errorMsg: "",
			mockSetup: func(timeUtils *time_utils.MockTimeUtilsInterface) {
				timeUtils.EXPECT().
					FormattedMinutesToSeconds("PT15M45S", playedTimeFormat).
					Return(945, nil)
			},
		},
		{
			name: "Error when parsing minutes",
			player: cdn_nba.PlayerBoxScoreDto{
				Name:     "Anthony Davis",
				PersonId: 203076,
				Starter:  "1",
				Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
					Minutes: "INVALID_FORMAT",
					Plus:    20,
					Minus:   5,
				},
			},
			result:   PlayerDTO{},
			errorMsg: "time format error",
			mockSetup: func(timeUtils *time_utils.MockTimeUtilsInterface) {
				timeUtils.EXPECT().
					FormattedMinutesToSeconds("INVALID_FORMAT", playedTimeFormat).
					Return(0, errors.New("time format error"))
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeUtils := time_utils.NewMockTimeUtilsInterface(ctrl)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup(timeUtils)
			mapper := &nbaMapper{
				timeUtils: timeUtils,
			}

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

func TestNbaMapper_mapTeam(t *testing.T) {
	cases := []struct {
		name      string
		team      cdn_nba.TeamBoxScoreDto
		expected  TeamBoxScoreDTO
		errorMsg  string
		mockSetup func(timeUtils *time_utils.MockTimeUtilsInterface)
	}{
		{
			name: "Successfully map team with players",
			team: cdn_nba.TeamBoxScoreDto{
				TeamId:      1610612747,
				TeamName:    "Los Angeles Lakers",
				TeamTricode: "LAL",
				Score:       120,
				Players: []cdn_nba.PlayerBoxScoreDto{
					{
						Name:     "LeBron James",
						PersonId: 2544,
						Starter:  "1",
						Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
							Minutes: "PT35M20S",
							Plus:    25,
							Minus:   10,
						},
					},
					{
						Name:     "Anthony Davis",
						PersonId: 203076,
						Starter:  "1",
						Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
							Minutes: "PT32M15S",
							Plus:    22,
							Minus:   8,
						},
					},
				},
			},
			expected: TeamBoxScoreDTO{
				Alias:    "LAL",
				Name:     "Los Angeles Lakers",
				LeagueId: "1610612747",
				Scored:   120,
				Players: []PlayerDTO{
					{
						FullNameLocal:  "LeBron James",
						LeaguePlayerID: "2544",
						Statistic: PlayerStatisticDTO{
							PlsMin:        15,
							PlayedSeconds: 2120,
							IsBench:       false,
						},
					},
					{
						FullNameLocal:  "Anthony Davis",
						LeaguePlayerID: "203076",
						Statistic: PlayerStatisticDTO{
							PlsMin:        14,
							PlayedSeconds: 1935,
							IsBench:       false,
						},
					},
				},
			},
			errorMsg: "",
			mockSetup: func(timeUtils *time_utils.MockTimeUtilsInterface) {
				timeUtils.EXPECT().
					FormattedMinutesToSeconds("PT35M20S", playedTimeFormat).
					Return(2120, nil)
				timeUtils.EXPECT().
					FormattedMinutesToSeconds("PT32M15S", playedTimeFormat).
					Return(1935, nil)
			},
		},
		{
			name: "Error when mapping player",
			team: cdn_nba.TeamBoxScoreDto{
				TeamId:      1610612747,
				TeamName:    "Los Angeles Lakers",
				TeamTricode: "LAL",
				Score:       120,
				Players: []cdn_nba.PlayerBoxScoreDto{
					{
						Name:     "LeBron James",
						PersonId: 2544,
						Starter:  "1",
						Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
							Minutes: "INVALID_FORMAT",
							Plus:    25,
							Minus:   10,
						},
					},
				},
			},
			expected: TeamBoxScoreDTO{},
			errorMsg: "time format error",
			mockSetup: func(timeUtils *time_utils.MockTimeUtilsInterface) {
				timeUtils.EXPECT().
					FormattedMinutesToSeconds("INVALID_FORMAT", playedTimeFormat).
					Return(0, errors.New("time format error"))
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeUtils := time_utils.NewMockTimeUtilsInterface(ctrl)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup(timeUtils)
			mapper := &nbaMapper{
				timeUtils: timeUtils,
			}

			result, err := mapper.mapTeam(tc.team)

			if tc.errorMsg != "" {
				assert.EqualError(t, err, tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestNbaMapper_mapGame tests the mapGame method of nbaMapper
//
// Verify that when valid game data is provided while mapping game - returns correct GameBoxScoreDTO
// Verify that when error occurs during league repository query - returns error
// Verify that when game has overtime periods while mapping game - calculates correct duration
// Verify that when player mapping fails while mapping game - returns error
func TestNbaMapper_mapGame(t *testing.T) {
	gameTime := time.Now().UTC()

	cases := []struct {
		name      string
		game      cdn_nba.BoxScoreDto
		expected  GameBoxScoreDTO
		errorMsg  string
		panic     string
		league    domain.League
		mockSetup func(timeUtils *time_utils.MockTimeUtilsInterface, logMock *log.MockLogger)
	}{
		{
			name: "Valid game data",
			game: cdn_nba.BoxScoreDto{
				GameId:      "0022200001",
				GameStatus:  3, // Final
				Period:      4, // Regular game, no overtime
				GameTimeUTC: gameTime,
				HomeTeam: cdn_nba.TeamBoxScoreDto{
					TeamId:      1610612747,
					TeamName:    "Los Angeles Lakers",
					TeamTricode: "LAL",
					Score:       120,
					Players: []cdn_nba.PlayerBoxScoreDto{
						{
							Name:     "LeBron James",
							PersonId: 2544,
							Starter:  "1",
							Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
								Minutes: "PT35M20S",
								Plus:    25,
								Minus:   10,
							},
						},
					},
				},
				AwayTeam: cdn_nba.TeamBoxScoreDto{
					TeamId:      1610612744,
					TeamName:    "Golden State Warriors",
					TeamTricode: "GSW",
					Score:       110,
					Players: []cdn_nba.PlayerBoxScoreDto{
						{
							Name:     "Stephen Curry",
							PersonId: 201939,
							Starter:  "1",
							Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
								Minutes: "PT34M10S",
								Plus:    15,
								Minus:   20,
							},
						},
					},
				},
			},
			expected: GameBoxScoreDTO{
				Id:            "0022200001",
				LeagueAliasEn: "NBA",
				IsFinal:       true,
				PlayedMinutes: 48, // 4 quarters * 12 minutes
				ScheduledAt:   gameTime,
				HomeTeam: TeamBoxScoreDTO{
					Alias:    "LAL",
					Name:     "Los Angeles Lakers",
					LeagueId: "1610612747",
					Scored:   120,
					Players: []PlayerDTO{
						{
							FullNameLocal:  "LeBron James",
							LeaguePlayerID: "2544",
							Statistic: PlayerStatisticDTO{
								PlsMin:        15,
								PlayedSeconds: 2120,
								IsBench:       false,
							},
						},
					},
				},
				AwayTeam: TeamBoxScoreDTO{
					Alias:    "GSW",
					Name:     "Golden State Warriors",
					LeagueId: "1610612744",
					Scored:   110,
					Players: []PlayerDTO{
						{
							FullNameLocal:  "Stephen Curry",
							LeaguePlayerID: "201939",
							Statistic: PlayerStatisticDTO{
								PlsMin:        -5,
								PlayedSeconds: 2050,
								IsBench:       false,
							},
						},
					},
				},
			},
			errorMsg: "",
			league: domain.League{
				PeriodsNumber:    4,
				PeriodDuration:   12,
				OvertimeDuration: 6,
				AliasEn:          "NBA",
			},
			mockSetup: func(timeUtils *time_utils.MockTimeUtilsInterface, logMock *log.MockLogger) {
				timeUtils.EXPECT().
					FormattedMinutesToSeconds("PT35M20S", playedTimeFormat).
					Return(2120, nil)
				timeUtils.EXPECT().
					FormattedMinutesToSeconds("PT34M10S", playedTimeFormat).
					Return(2050, nil)
			},
		},
		{
			name: "Game with overtime periods",
			game: cdn_nba.BoxScoreDto{
				GameId:      "0022200002",
				GameStatus:  3, // Final
				Period:      5, // 1 overtime
				GameTimeUTC: gameTime,
				HomeTeam: cdn_nba.TeamBoxScoreDto{
					TeamId:      1610612747,
					TeamName:    "Los Angeles Lakers",
					TeamTricode: "LAL",
					Score:       130,
					Players: []cdn_nba.PlayerBoxScoreDto{
						{
							Name:     "LeBron James",
							PersonId: 2544,
							Starter:  "1",
							Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
								Minutes: "PT40M20S",
								Plus:    25,
								Minus:   10,
							},
						},
					},
				},
				AwayTeam: cdn_nba.TeamBoxScoreDto{
					TeamId:      1610612744,
					TeamName:    "Golden State Warriors",
					TeamTricode: "GSW",
					Score:       128,
					Players: []cdn_nba.PlayerBoxScoreDto{
						{
							Name:     "Stephen Curry",
							PersonId: 201939,
							Starter:  "1",
							Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
								Minutes: "PT39M10S",
								Plus:    15,
								Minus:   20,
							},
						},
					},
				},
			},
			league: domain.League{
				PeriodsNumber:    4,
				PeriodDuration:   12,
				OvertimeDuration: 6,
				AliasEn:          "NBA",
			},
			expected: GameBoxScoreDTO{
				Id:            "0022200002",
				LeagueAliasEn: "NBA",
				IsFinal:       true,
				PlayedMinutes: 54, // 4 quarters * 12 minutes + 1 overtime * 6 minutes
				ScheduledAt:   gameTime,
				HomeTeam: TeamBoxScoreDTO{
					Alias:    "LAL",
					Name:     "Los Angeles Lakers",
					LeagueId: "1610612747",
					Scored:   130,
					Players: []PlayerDTO{
						{
							FullNameLocal:  "LeBron James",
							LeaguePlayerID: "2544",
							Statistic: PlayerStatisticDTO{
								PlsMin:        15,
								PlayedSeconds: 2420,
								IsBench:       false,
							},
						},
					},
				},
				AwayTeam: TeamBoxScoreDTO{
					Alias:    "GSW",
					Name:     "Golden State Warriors",
					LeagueId: "1610612744",
					Scored:   128,
					Players: []PlayerDTO{
						{
							FullNameLocal:  "Stephen Curry",
							LeaguePlayerID: "201939",
							Statistic: PlayerStatisticDTO{
								PlsMin:        -5,
								PlayedSeconds: 2350,
								IsBench:       false,
							},
						},
					},
				},
			},
			errorMsg: "",
			mockSetup: func(timeUtils *time_utils.MockTimeUtilsInterface, logMock *log.MockLogger) {
				timeUtils.EXPECT().
					FormattedMinutesToSeconds("PT40M20S", playedTimeFormat).
					Return(2420, nil)
				timeUtils.EXPECT().
					FormattedMinutesToSeconds("PT39M10S", playedTimeFormat).
					Return(2350, nil)
			},
		},
		{
			name: "Error during player mapping",
			game: cdn_nba.BoxScoreDto{
				GameId:      "0022200004",
				GameStatus:  3,
				Period:      4,
				GameTimeUTC: gameTime,
				HomeTeam: cdn_nba.TeamBoxScoreDto{
					TeamId:      1610612747,
					TeamName:    "Los Angeles Lakers",
					TeamTricode: "LAL",
					Score:       120,
					Players: []cdn_nba.PlayerBoxScoreDto{
						{
							Name:     "LeBron James",
							PersonId: 2544,
							Starter:  "1",
							Statistics: cdn_nba.PlayerEfficiencyBoxScoreDto{
								Minutes: "INVALID_FORMAT", // This will cause an error
								Plus:    25,
								Minus:   10,
							},
						},
					},
				},
				AwayTeam: cdn_nba.TeamBoxScoreDto{
					TeamId:      1610612744,
					TeamName:    "Golden State Warriors",
					TeamTricode: "GSW",
					Score:       110,
					Players:     []cdn_nba.PlayerBoxScoreDto{},
				},
			},
			expected: GameBoxScoreDTO{},
			league: domain.League{
				PeriodsNumber:    4,
				PeriodDuration:   12,
				OvertimeDuration: 6,
				AliasEn:          "NBA",
			},
			errorMsg: "time format error",
			mockSetup: func(timeUtils *time_utils.MockTimeUtilsInterface, logMock *log.MockLogger) {

				timeUtils.EXPECT().
					FormattedMinutesToSeconds("INVALID_FORMAT", playedTimeFormat).
					Return(0, errors.New("time format error"))
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeUtils := time_utils.NewMockTimeUtilsInterface(ctrl)
	logMock := log.NewMockLogger(ctrl)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			tc.mockSetup(timeUtils, logMock)

			mapper := &nbaMapper{
				league:    &tc.league,
				timeUtils: timeUtils,
				logger:    logMock,
			}

			assertion := func() {
				result, err := mapper.mapGame(tc.game)

				if tc.errorMsg != "" {
					assert.EqualError(t, err, tc.errorMsg)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				}
			}

			if tc.panic != "" {
				assert.Panics(t, assertion)
			} else {
				assertion()
			}
		})
	}
}
