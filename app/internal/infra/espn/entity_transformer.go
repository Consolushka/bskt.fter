package espn

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/pkg/statsutil"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Коды команд в БД заданы первым провайдером (api-sports), а заголовок игры
// входит в unique-индекс games — расходящиеся аббревиатуры ESPN приводим к ним.
var canonicalTeamCode = map[string]string{
	"GS":   "GSW",
	"NO":   "NOP",
	"NY":   "NYK",
	"SA":   "SAS",
	"UTAH": "UTA",
	"WSH":  "WAS",
}

func teamCode(abbreviation string) string {
	if code, ok := canonicalTeamCode[abbreviation]; ok {
		return code
	}
	return abbreviation
}

const (
	seasonTypeRegular    = 2
	seasonTypePostseason = 3
	seasonTypePlayIn     = 5

	competitionTypeAllStar = "ALLSTAR"

	regulationPeriods       = 4
	regulationPeriodMinutes = 12
	overtimePeriodMinutes   = 5

	// ESPN отдаёт минуты без секунд; запасной разбор — time.RFC3339.
	eventDateLayout = "2006-01-02T15:04Z"
)

var requiredLabels = []string{"MIN", "PTS", "REB", "AST", "TO", "STL", "BLK", "FG", "+/-"}

type EntityTransformer struct{}

func NewEntityTransformer() EntityTransformer {
	return EntityTransformer{}
}

// TransformEvent превращает событие scoreboard в игру. ok == false: событие отфильтровано
// (не завершено, предсезонка, звёздный уикенд). Ошибка — только для завершённого
// события, которое не разбирается; вызывающий решает, что с ним делать.
func (e EntityTransformer) TransformEvent(event EventEntity) (games.GameStatEntity, bool, error) {
	if !isImportable(event) {
		return games.GameStatEntity{}, false, nil
	}
	// Без id игра прошла бы мимо обогащения и сохранилась бы без игроков — навсегда.
	if event.Id == "" {
		return games.GameStatEntity{}, false, errors.New("event without id")
	}
	comp := event.Competitions[0]

	var home, away *CompetitorEntity
	for i := range comp.Competitors {
		switch comp.Competitors[i].HomeAway {
		case "home":
			home = &comp.Competitors[i]
		case "away":
			away = &comp.Competitors[i]
		}
	}
	if home == nil || away == nil {
		return games.GameStatEntity{}, false, fmt.Errorf("event %s: home or away competitor missing", event.Id)
	}

	scheduledAt, err := parseEventDate(event.Date)
	if err != nil {
		return games.GameStatEntity{}, false, fmt.Errorf("event %s: date %q: %w", event.Id, event.Date, err)
	}

	homeScore, err := strconv.Atoi(home.Score)
	if err != nil {
		return games.GameStatEntity{}, false, fmt.Errorf("event %s: home score: %w", event.Id, err)
	}
	awayScore, err := strconv.Atoi(away.Score)
	if err != nil {
		return games.GameStatEntity{}, false, fmt.Errorf("event %s: away score: %w", event.Id, err)
	}
	homeId, err := strconv.Atoi(home.Team.Id)
	if err != nil {
		return games.GameStatEntity{}, false, fmt.Errorf("event %s: home team id: %w", event.Id, err)
	}
	awayId, err := strconv.Atoi(away.Team.Id)
	if err != nil {
		return games.GameStatEntity{}, false, fmt.Errorf("event %s: away team id: %w", event.Id, err)
	}

	overtimes := max(0, comp.Status.Period-regulationPeriods)

	return games.GameStatEntity{
		GameModel: games.GameModel{
			ScheduledAt: scheduledAt,
			Title:       teamCode(home.Team.Abbreviation) + " - " + teamCode(away.Team.Abbreviation),
			Duration:    regulationPeriods*regulationPeriodMinutes + overtimes*overtimePeriodMinutes,
		},
		ExternalGameId:     event.Id,
		HomeTeamExternalId: homeId,
		AwayTeamExternalId: awayId,
		HomeTeamStat:       teamStat(home.Team, homeScore, awayScore),
		AwayTeamStat:       teamStat(away.Team, awayScore, homeScore),
	}, true, nil
}

func isImportable(event EventEntity) bool {
	if len(event.Competitions) == 0 {
		return false
	}
	comp := event.Competitions[0]
	if !comp.Status.Type.Completed || comp.Type.Abbreviation == competitionTypeAllStar {
		return false
	}
	switch event.Season.Type {
	case seasonTypeRegular, seasonTypePostseason, seasonTypePlayIn:
		return true
	}
	return false
}

func parseEventDate(date string) (time.Time, error) {
	if t, err := time.Parse(eventDateLayout, date); err == nil {
		return t, nil
	}
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse date: %w", err)
	}
	return t, nil
}

func teamStat(team TeamEntity, score, opponentScore int) teams.TeamStatEntity {
	return teams.TeamStatEntity{
		TeamModel: teams.TeamModel{
			Name:     team.Name,
			HomeTown: team.Location,
			Alias:    teamCode(team.Abbreviation),
		},
		GameTeamStatModel: teams.GameTeamStatModel{
			Score:     score,
			FinalDiff: score - opponentScore,
		},
	}
}

// MapBoxScore раскладывает игроков из summary по командам игры.
func (e EntityTransformer) MapBoxScore(resp SummaryResponse, game *games.GameStatEntity) error {
	seenHome, seenAway := false, false

	for _, block := range resp.Boxscore.Players {
		teamId, err := strconv.Atoi(block.Team.Id)
		if err != nil {
			return fmt.Errorf("box score team id %q: %w", block.Team.Id, err)
		}

		var target *teams.TeamStatEntity
		switch teamId {
		case game.HomeTeamExternalId:
			target, seenHome = &game.HomeTeamStat, true
		case game.AwayTeamExternalId:
			target, seenAway = &game.AwayTeamStat, true
		default:
			return fmt.Errorf("box score team %s is neither home nor away", block.Team.Id)
		}

		if len(block.Statistics) == 0 {
			return fmt.Errorf("box score for team %d has no statistics", teamId)
		}
		stats := block.Statistics[0]

		index := make(map[string]int, len(stats.Labels))
		for i, label := range stats.Labels {
			index[label] = i
		}
		for _, label := range requiredLabels {
			if _, ok := index[label]; !ok {
				return fmt.Errorf("box score label %q missing", label)
			}
		}

		result := []players.PlayerStatisticEntity{}
		for _, a := range stats.Athletes {
			if a.DidNotPlay || len(a.Stats) == 0 {
				continue
			}
			if len(a.Stats) != len(stats.Labels) {
				return fmt.Errorf("athlete %s: stats length %d != labels %d", a.Athlete.Id, len(a.Stats), len(stats.Labels))
			}
			// ESPN иногда отдаёт игрока со скамейки без пометки «не играл», но с минутами "--".
			// Если у него всё же есть очки, это поймает сверка итогов команды.
			if a.Stats[index["MIN"]] == "--" {
				continue
			}
			p, err := mapAthlete(index, a)
			if err != nil {
				return err
			}
			result = append(result, p)
		}
		target.PlayerStats = result
	}

	if !seenHome {
		return fmt.Errorf("box score for team %d missing", game.HomeTeamExternalId)
	}
	if !seenAway {
		return fmt.Errorf("box score for team %d missing", game.AwayTeamExternalId)
	}
	return nil
}

func mapAthlete(index map[string]int, a AthleteStatsEntity) (players.PlayerStatisticEntity, error) {
	num := func(label string) (int, error) {
		v := a.Stats[index[label]]
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("athlete %s: %s %q: %w", a.Athlete.Id, label, v, err)
		}
		return n, nil
	}

	vals := map[string]int{}
	for _, label := range requiredLabels {
		if label == "FG" {
			continue
		}
		n, err := num(label)
		if err != nil {
			return players.PlayerStatisticEntity{}, err
		}
		vals[label] = n
	}

	fg := a.Stats[index["FG"]]
	madeStr, attemptedStr, _ := strings.Cut(fg, "-")
	made, err := strconv.Atoi(madeStr)
	if err != nil {
		return players.PlayerStatisticEntity{}, fmt.Errorf("athlete %s: FG %q: %w", a.Athlete.Id, fg, err)
	}
	attempted, err := strconv.Atoi(attemptedStr)
	if err != nil {
		return players.PlayerStatisticEntity{}, fmt.Errorf("athlete %s: FG %q: %w", a.Athlete.Id, fg, err)
	}

	return players.PlayerStatisticEntity{
		PlayerModel:      players.PlayerModel{FullName: a.Athlete.DisplayName},
		PlayerExternalId: a.Athlete.Id,
		GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
			PlayedSeconds:        vals["MIN"] * 60,
			PlsMin:               int8(vals["+/-"]),
			Points:               uint8(vals["PTS"]),
			Assists:              uint8(vals["AST"]),
			Rebounds:             uint8(vals["REB"]),
			Steals:               uint8(vals["STL"]),
			Blocks:               uint8(vals["BLK"]),
			Turnovers:            uint8(vals["TO"]),
			FieldGoalsPercentage: statsutil.FromRatio(made, attempted),
		},
	}, nil
}

// VerifyTotals сверяет сумму по игрокам с итогами команды: если не сходится,
// разбор box score нельзя сохранять — данные неполные или перепутаны.
func (e EntityTransformer) VerifyTotals(game games.GameStatEntity) error {
	if err := verifyTeam(game, game.HomeTeamStat); err != nil {
		return err
	}
	return verifyTeam(game, game.AwayTeamStat)
}

func verifyTeam(game games.GameStatEntity, stat teams.TeamStatEntity) error {
	alias := stat.TeamModel.Alias

	played := len(stat.PlayerStats)
	if played == 0 {
		return fmt.Errorf("team %s: no players with stats", alias)
	}

	sumPoints, sumSeconds := 0, 0
	for _, p := range stat.PlayerStats {
		sumPoints += int(p.GameTeamPlayerStatModel.Points)
		sumSeconds += p.GameTeamPlayerStatModel.PlayedSeconds
	}
	if sumPoints != stat.GameTeamStatModel.Score {
		return fmt.Errorf("team %s: players points %d, team total %d", alias, sumPoints, stat.GameTeamStatModel.Score)
	}

	sumMinutes := sumSeconds / 60
	expectedMinutes := 5 * game.GameModel.Duration
	tolerance := (played + 1) / 2
	diff := sumMinutes - expectedMinutes
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		return fmt.Errorf("team %s: players minutes %d, expected %d ± %d", alias, sumMinutes, expectedMinutes, tolerance)
	}
	return nil
}

func (e EntityTransformer) MapAthleteBio(resp AthleteResponse) (players.PlayerBioEntity, error) {
	if len(resp.DateOfBirth) < 10 {
		return players.PlayerBioEntity{}, fmt.Errorf("athlete %s: dateOfBirth is empty", resp.Id)
	}
	birthDate, err := time.Parse("2006-01-02", resp.DateOfBirth[:10])
	if err != nil {
		return players.PlayerBioEntity{}, fmt.Errorf("athlete %s: dateOfBirth %q: %w", resp.Id, resp.DateOfBirth, err)
	}
	if resp.DisplayName == "" {
		return players.PlayerBioEntity{}, fmt.Errorf("athlete %s: displayName is empty", resp.Id)
	}
	return players.PlayerBioEntity{FullName: resp.DisplayName, BirthDate: birthDate}, nil
}
