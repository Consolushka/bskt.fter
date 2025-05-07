package mock_data

import (
	"IMP/app/internal/statistics/infobasket"
	"fmt"
	"github.com/go-faker/faker/v4"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

// CreateMockTeamBoxScoreDto creates a mock TeamBoxScoreDto with specified parameters
func CreateMockTeamBoxScoreDto(alias string, teamName string, score int, teamID int, playersCount int) infobasket.TeamBoxScoreDto {
	rand.Seed(time.Now().UnixNano())

	// Create players
	players := make([]infobasket.PlayerBoxScoreDto, 0, playersCount)
	for i := 0; i < playersCount; i++ {
		isStart := i < 5                // First 5 players are starters
		plusMinus := rand.Intn(41) - 20 // Random plus-minus between -20 and +20
		seconds := rand.Intn(2400)      // Random playing time up to 40 minutes (2400 seconds)

		// player name will be generated using faker
		player := CreateMockPlayer(
			1000+i,
			plusMinus,
			seconds,
			"",
			"",
			"",
			"",
			time.Now().AddDate(-rand.Intn(15)-18, -rand.Intn(12), -rand.Intn(30)), // Random birthdate
			isStart,
		)
		players = append(players, player)
	}

	// Generate random stats
	shot1 := rand.Intn(40) + 10
	goal1 := rand.Intn(shot1)
	shot2 := rand.Intn(70) + 30
	goal2 := rand.Intn(shot2)
	shot3 := rand.Intn(30) + 10
	goal3 := rand.Intn(shot3)
	paintShot := rand.Intn(40) + 20
	paintGoal := rand.Intn(paintShot)

	// Calculate percentages
	shot1Percent := "0.0"
	if shot1 > 0 {
		shot1Percent = fmt.Sprintf("%.1f", float64(goal1)/float64(shot1)*100)
	}

	shot2Percent := "0.0"
	if shot2 > 0 {
		shot2Percent = fmt.Sprintf("%.1f", float64(goal2)/float64(shot2)*100)
	}

	shot3Percent := "0.0"
	if shot3 > 0 {
		shot3Percent = fmt.Sprintf("%.1f", float64(goal3)/float64(shot3)*100)
	}

	paintShotPercent := "0.0"
	if paintShot > 0 {
		paintShotPercent = fmt.Sprintf("%.1f", float64(paintGoal)/float64(paintShot)*100)
	}

	// Random values for other stats
	assist := rand.Intn(30) + 5
	blocks := rand.Intn(10) + 1
	defRebound := rand.Intn(30) + 10
	offRebound := rand.Intn(20) + 5
	rebound := defRebound + offRebound
	steal := rand.Intn(15) + 3
	turnover := rand.Intn(20) + 5
	teamDefRebound := rand.Intn(10)
	teamOffRebound := rand.Intn(10)
	teamRebound := teamDefRebound + teamOffRebound
	teamSteal := rand.Intn(5)
	teamTurnover := rand.Intn(8)
	foul := rand.Intn(25) + 10
	opponentFoul := rand.Intn(25) + 10
	seconds := rand.Intn(2400) + 1200 // Team played between 20 and 60 minutes
	minutes := seconds / 60
	remainingSeconds := seconds % 60

	return infobasket.TeamBoxScoreDto{
		TeamNumber: 1,
		TeamID:     teamID,
		TeamName: infobasket.TeamNameBoxScoreDto{
			CompTeamNameID:       teamID,
			TeamID:               teamID,
			TeamType:             1,
			CompTeamShortNameRu:  alias,
			CompTeamShortNameEn:  alias,
			CompTeamNameRu:       teamName,
			CompTeamNameEn:       teamName,
			CompTeamRegionNameRu: "Регион",
			CompTeamRegionNameEn: "Region",
			CompTeamAbcNameRu:    alias,
			CompTeamAbcNameEn:    alias,
			CompTeamNameDefault:  true,
			SysStatus:            1,
			SysLastChanged:       time.Now().Format("2006-01-02 15:04:05"),
			IsRealTeam:           true,
		},
		Score:            score,
		Points:           score,
		Shot1:            shot1,
		Goal1:            goal1,
		Shot2:            shot2,
		Goal2:            goal2,
		Shot3:            shot3,
		Goal3:            goal3,
		PaintShot:        paintShot,
		PaintGoal:        paintGoal,
		Shots1:           fmt.Sprintf("%d/%d", goal1, shot1),
		Shot1Percent:     shot1Percent,
		Shots2:           fmt.Sprintf("%d/%d", goal2, shot2),
		Shot2Percent:     shot2Percent,
		Shots3:           fmt.Sprintf("%d/%d", goal3, shot3),
		Shot3Percent:     shot3Percent,
		PaintShots:       fmt.Sprintf("%d/%d", paintGoal, paintShot),
		PaintShotPercent: paintShotPercent,
		Assist:           assist,
		Blocks:           blocks,
		DefRebound:       defRebound,
		OffRebound:       offRebound,
		Rebound:          rebound,
		Steal:            steal,
		Turnover:         turnover,
		TeamDefRebound:   teamDefRebound,
		TeamOffRebound:   &teamOffRebound,
		TeamRebound:      teamRebound,
		TeamSteal:        &teamSteal,
		TeamTurnover:     teamTurnover,
		Foul:             foul,
		OpponentFoul:     opponentFoul,
		Seconds:          seconds,
		PlayedTime:       fmt.Sprintf("%d:%02d", minutes, remainingSeconds),
		PlusMinus:        nil, // Team plus-minus is often nil in the model
		Players:          players,
	}
}

func CreateMockGameScheduleDto(gameId int, gameDate string, status int) infobasket.GameScheduleDto {
	return infobasket.GameScheduleDto{
		GameID:                 gameId,
		IsToday:                false,
		DaysFromToday:          0,
		GameNumber:             "G-001",
		GameDateInt:            20230101,
		GameDate:               gameDate,
		HasTime:                true,
		GameTime:               "19:00",
		GameTimeMsk:            "21:00",
		GameLocalDate:          gameDate,
		GameDateTime:           gameDate + " 19:00",
		GameDateTimeMoscow:     gameDate + " 21:00",
		DisplayDateTimeLocal:   gameDate + " 19:00",
		DisplayDateTimeMsk:     gameDate + " 21:00",
		DisplayDateTimeLocalEn: gameDate + " 7:00 PM",
		DisplayDateTimeMskEn:   gameDate + " 9:00 PM",
		GameStatus:             status,
		TeamAid:                101,
		TeamBid:                102,
		ShortTeamNameAru:       "Команда A",
		ShortTeamNameBru:       "Команда B",
		TeamNameAru:            "Полное имя команды A",
		TeamNameBru:            "Полное имя команды B",
		ShortTeamNameAen:       "Team A",
		ShortTeamNameBen:       "Team B",
		TeamNameAen:            "Full Team A Name",
		TeamNameBen:            "Full Team B Name",
		CompTeamNameAen:        "Team A",
		CompTeamNameBen:        "Team B",
		CompTeamNameAru:        "Команда A",
		CompTeamNameBru:        "Команда B",
		RegionTeamNameAru:      "Регион A",
		RegionTeamNameBru:      "Регион B",
		RegionTeamNameAen:      "Region A",
		RegionTeamNameBen:      "Region B",
		ScoreA:                 85,
		ScoreB:                 77,
		Periods:                4,
		TeamLogoA:              "https://example.com/logo_a.png",
		TeamLogoB:              "https://example.com/logo_b.png",
		HasVideo:               true,
		LiveID:                 "live123",
		TvRu:                   "Матч ТВ",
		TvEn:                   "Match TV",
		VideoID:                nil,
		RegionRu:               "Москва",
		RegionEn:               "Moscow",
		RegionId:               1,
		ArenaEn:                "Basketball Arena",
		ArenaRu:                "Баскетбольная Арена",
		ArenaId:                42,
		GameAttendance:         5000,
		CompNameRu:             "Чемпионат",
		CompNameEn:             "Championship",
		LeagueNameRu:           "Лига",
		LeagueNameEn:           "League",
		LeagueShortNameRu:      "Л",
		LeagueShortNameEn:      "L",
		Gender:                 1,
	}
}

// CreateMockPlayer creates a mock player with the provided parameters. If Name is empty, it will be generated using faker.
func CreateMockPlayer(personID, plusMinus, seconds int, firstNameRu, lastNameRu, firstNameEn, lastNameEn string, birthDate time.Time, isStart bool) infobasket.PlayerBoxScoreDto {
	rand.Seed(time.Now().UnixNano())

	if firstNameRu == "" {
		firstNameVal, _ := faker.GetPerson().RussianFirstNameMale(reflect.Value{})

		firstNameRu = firstNameVal.(string)
	}

	if lastNameRu == "" {
		lastNameVal, _ := faker.GetPerson().RussianLastNameMale(reflect.Value{})

		lastNameRu = lastNameVal.(string)
	}

	if firstNameEn == "" {
		firstNameVal, _ := faker.GetPerson().FirstNameMale(reflect.Value{})

		firstNameEn = firstNameVal.(string)
	}

	if lastNameEn == "" {
		lastNameVal, _ := faker.GetPerson().LastName(reflect.Value{})

		lastNameEn = lastNameVal.(string)
	}

	shot1 := rand.Intn(10) + 1
	goal1 := rand.Intn(shot1 + 1)
	shot1Percent := fmt.Sprintf("%.1f", float64(goal1)/float64(shot1)*100)

	shot2 := rand.Intn(15) + 1
	goal2 := rand.Intn(shot2 + 1)
	shot2Percent := fmt.Sprintf("%.1f", float64(goal2)/float64(shot2)*100)

	paintShot := rand.Intn(10) + 1
	paintGoal := rand.Intn(paintShot + 1)
	paintShotPercent := fmt.Sprintf("%.1f", float64(paintGoal)/float64(paintShot)*100)

	shot3 := rand.Intn(8) + 1
	goal3 := rand.Intn(shot3 + 1)
	shot3Percent := fmt.Sprintf("%.1f", float64(goal3)/float64(shot3)*100)

	assist := rand.Intn(10)
	blocks := rand.Intn(5)
	defRebound := rand.Intn(10)
	offRebound := rand.Intn(5)
	rebound := defRebound + offRebound
	steal := rand.Intn(5)
	turnover := rand.Intn(5)
	foul := rand.Intn(5)
	opponentFoul := rand.Intn(6)

	minutes := seconds / 60
	remainingSeconds := seconds % 60
	playedTime := fmt.Sprintf("%d:%02d", minutes, remainingSeconds)

	return infobasket.PlayerBoxScoreDto{
		PersonID:         personID,
		TeamNumber:       1,
		PlayerNumber:     rand.Intn(99) + 1,
		DisplayNumber:    strconv.Itoa(rand.Intn(99) + 1),
		LastNameRu:       lastNameRu,
		LastNameEn:       lastNameEn,
		FirstNameRu:      firstNameRu,
		FirstNameEn:      firstNameEn,
		PersonNameRu:     lastNameRu + " " + firstNameRu,
		PersonNameEn:     lastNameEn + " " + firstNameEn,
		Capitan:          rand.Intn(2),
		PersonBirth:      birthDate.Format("02.01.2006"),
		PosID:            rand.Intn(5) + 1,
		CountryCodeIOC:   "RUS",
		CountryNameRu:    "Россия",
		CountryNameEn:    "Russia",
		RankRu:           []string{"КМС", "МС", "МСМК"}[rand.Intn(3)],
		Height:           170 + rand.Intn(40),
		Weight:           70 + rand.Intn(40),
		Points:           goal1 + goal2*2 + goal3*3,
		Shot1:            shot1,
		Goal1:            goal1,
		Shots1:           fmt.Sprintf("%d/%d", goal1, shot1),
		Shot1Percent:     shot1Percent,
		Shot2:            shot2,
		Goal2:            goal2,
		Shots2:           fmt.Sprintf("%d/%d", goal2, shot2),
		Shot2Percent:     shot2Percent,
		PaintShot:        paintShot,
		PaintGoal:        paintGoal,
		PaintShots:       fmt.Sprintf("%d/%d", paintGoal, paintShot),
		PaintShotPercent: paintShotPercent,
		Shot3:            shot3,
		Goal3:            goal3,
		Shots3:           fmt.Sprintf("%d/%d", goal3, shot3),
		Shot3Percent:     shot3Percent,
		Assist:           assist,
		Blocks:           blocks,
		DefRebound:       defRebound,
		OffRebound:       offRebound,
		Rebound:          rebound,
		Steal:            steal,
		Turnover:         turnover,
		Foul:             foul,
		OpponentFoul:     opponentFoul,
		PlusMinus:        plusMinus,
		Seconds:          seconds,
		PlayedTime:       playedTime,
		IsStart:          isStart,
		StartMark: func() string {
			if isStart {
				return "+"
			} else {
				return ""
			}
		}(),
	}
}
