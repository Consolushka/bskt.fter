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
