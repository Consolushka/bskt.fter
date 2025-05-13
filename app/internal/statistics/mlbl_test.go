package statistics

import (
	"IMP/app/internal/statistics/infobasket"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMlblMapper_mapPlayer(t *testing.T) {
	firstPlayerDate := time.Date(1990, 11, 25, 0, 0, 0, 0, time.UTC)
	secondPlayerDate := time.Date(1970, 11, 11, 0, 0, 0, 0, time.UTC)
	thirdPlayerDate := time.Date(2000, 12, 13, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		player   infobasket.PlayerBoxScoreDto
		result   PlayerDTO
		errorMsg string
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
		},
	}

	mapper := newMlblMapper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

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

	mapper := newMlblMapper()
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

	mapper := newMlblMapper()
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
