package infobasket

import (
	"IMP/app/internal/infrastructure/infobasket"
	"IMP/app/internal/modules/statistics/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"strconv"
	"testing"
	"time"
)

func TestMapper_mapPlayer(t *testing.T) {
	mapper := newMapper()

	id := 23
	plusMnius := 10
	seconds := 2611
	isStart := false
	localFirstName := "Nikola"
	localLastName := "Jokić"
	enFirstName := "Nikola"
	enLastName := "Jokic"
	birthdate := time.Date(1995, 5, 12, 0, 0, 0, 0, time.UTC)

	infobasketDto := infobasket.CreateMockPlayer(id, plusMnius, seconds, localFirstName, localLastName, enFirstName, enLastName, birthdate, isStart)

	excpectedStatisticDto := models.PlayerDTO{
		FullNameLocal:  localLastName + " " + localFirstName,
		FullNameEn:     enLastName + " " + enFirstName,
		BirthDate:      &birthdate,
		LeaguePlayerID: strconv.Itoa(id),
		Statistic: models.PlayerStatisticDTO{
			PlsMin:        plusMnius,
			PlayedSeconds: seconds,
			IsBench:       !isStart,
		},
	}

	mapperDto := mapper.mapPlayer(infobasketDto)
	if diff := cmp.Diff(excpectedStatisticDto, mapperDto); diff != "" {
		t.Errorf("mapPlayer() mismatch (-want +got):\n%s", diff)
	}
}

func TestMapper_mapTeam(t *testing.T) {
	mapper := newMapper()

	teamId := 144
	alias := "PHX"
	name := "Phoenix Suns"
	score := 110
	playersCount := 12

	infobasketDto := infobasket.CreateMockTeamBoxScoreDto(alias, name, score, teamId, playersCount)

	excpectedStatisticDto := models.TeamBoxScoreDTO{
		Alias:    alias,
		Name:     name,
		LeagueId: strconv.Itoa(teamId),
		Scored:   score,
	}

	mapperDto := mapper.mapTeam(infobasketDto)
	if diff := cmp.Diff(excpectedStatisticDto, mapperDto, cmpopts.IgnoreFields(models.TeamBoxScoreDTO{}, "Players")); diff != "" {
		t.Errorf("mapTeam() mismatch (-want +got):\n%s", diff)
	}

	if equal := cmp.Equal(len(mapperDto.Players), playersCount); !equal {
		t.Errorf("Players count mismatch (-want +got):\n%s", "- "+strconv.Itoa(len(mapperDto.Players))+"\n+ "+strconv.Itoa(playersCount))
	}
}

func TestMapper_mapGame(t *testing.T) {
	mapper := newMapper()

	leagueAlias := "NBA"
	regulationPeriodsNumber := 4
	periodDuration := 12
	overtimeDuration := 6
	isFinal := true
	maxPeriods := 6
	var gameStatus int
	if isFinal {
		gameStatus = 1
	} else {
		gameStatus = 0
	}
	homeTeamId := 144
	awayTeamId := 145
	scheduled := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	excpectedStatisticDto := models.GameBoxScoreDTO{
		LeagueAliasEn: leagueAlias,
		IsFinal:       isFinal,
		HomeTeam: models.TeamBoxScoreDTO{
			LeagueId: strconv.Itoa(homeTeamId),
		},
		AwayTeam: models.TeamBoxScoreDTO{
			LeagueId: strconv.Itoa(awayTeamId),
		},
		PlayedMinutes: 60,
		ScheduledAt:   scheduled,
	}

	infobasketDto := infobasket.CreateMockGameBoxScoreResponse(homeTeamId, awayTeamId, gameStatus, maxPeriods, scheduled.Format("02.01.2006"), scheduled.Format("15.04"))

	mapperDto := mapper.mapGame(infobasketDto, regulationPeriodsNumber, periodDuration, overtimeDuration, leagueAlias)

	if diff := cmp.Diff(excpectedStatisticDto, *mapperDto, cmpopts.IgnoreFields(models.GameBoxScoreDTO{}, "HomeTeam", "AwayTeam")); diff != "" {
		t.Errorf("mapGame() mismatch (-want +got):\n%s", diff)
	}
}
