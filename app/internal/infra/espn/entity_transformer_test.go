package espn

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fixtureEventId = "401909888"

func loadFixture[T any](t *testing.T, name string) T {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	require.NoError(t, err)
	var v T
	require.NoError(t, json.Unmarshal(data, &v))
	return v
}

func manualEvent(seasonType, period int) EventEntity {
	return EventEntity{
		Id:     "1",
		Date:   "2026-10-29T00:00Z",
		Season: SeasonEntity{Type: seasonType},
		Competitions: []CompetitionEntity{{
			Type:   CompetitionTypeEntity{Abbreviation: "STD"},
			Status: StatusEntity{Type: StatusTypeEntity{Completed: true}, Period: period},
			Competitors: []CompetitorEntity{
				{HomeAway: "home", Score: "100", Team: TeamEntity{Id: "16", Abbreviation: "MIN"}},
				{HomeAway: "away", Score: "90", Team: TeamEntity{Id: "9", Abbreviation: "GS"}},
			},
		}},
	}
}

func fixtureGame(t *testing.T) games.GameStatEntity {
	t.Helper()
	sb := loadFixture[ScoreboardResponse](t, "scoreboard_20261028.json")
	for _, ev := range sb.Events {
		if ev.Id != fixtureEventId {
			continue
		}
		game, ok, err := NewEntityTransformer().TransformEvent(ev)
		require.NoError(t, err)
		require.True(t, ok)
		return game
	}
	t.Fatalf("event %s not found in fixture", fixtureEventId)
	return games.GameStatEntity{}
}

func fixtureGameWithBoxScore(t *testing.T) games.GameStatEntity {
	t.Helper()
	game := fixtureGame(t)
	summary := loadFixture[SummaryResponse](t, "summary_401909888.json")
	require.NoError(t, NewEntityTransformer().MapBoxScore(summary, &game))
	return game
}

func findPlayer(t *testing.T, stats []players.PlayerStatisticEntity, id string) (players.PlayerStatisticEntity, bool) {
	t.Helper()
	for _, p := range stats {
		if p.PlayerExternalId == id {
			return p, true
		}
	}
	return players.PlayerStatisticEntity{}, false
}

func TestEntityTransformer_TransformEvent_Fixture(t *testing.T) {
	sb := loadFixture[ScoreboardResponse](t, "scoreboard_20261028.json")
	require.Len(t, sb.Events, 12)

	tr := NewEntityTransformer()
	var target games.GameStatEntity
	for _, ev := range sb.Events {
		game, ok, err := tr.TransformEvent(ev)
		require.NoError(t, err)
		require.True(t, ok)
		if ev.Id == fixtureEventId {
			target = game
		}
	}

	assert.Equal(t, "MIN - GSW", target.GameModel.Title)
	assert.Equal(t, time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC), target.GameModel.ScheduledAt)
	assert.Equal(t, 48, target.GameModel.Duration)
	assert.Equal(t, fixtureEventId, target.ExternalGameId)
	assert.Equal(t, 16, target.HomeTeamExternalId)
	assert.Equal(t, 9, target.AwayTeamExternalId)
	assert.Equal(t, 114, target.HomeTeamStat.GameTeamStatModel.Score)
	assert.Equal(t, -5, target.HomeTeamStat.GameTeamStatModel.FinalDiff)
	assert.Equal(t, teams.TeamModel{Name: "Timberwolves", HomeTown: "Minnesota", Alias: "MIN"}, target.HomeTeamStat.TeamModel)
	assert.Equal(t, 119, target.AwayTeamStat.GameTeamStatModel.Score)
	assert.Equal(t, 5, target.AwayTeamStat.GameTeamStatModel.FinalDiff)
	assert.Equal(t, teams.TeamModel{Name: "Warriors", HomeTown: "Golden State", Alias: "GSW"}, target.AwayTeamStat.TeamModel)
	assert.Nil(t, target.HomeTeamStat.PlayerStats)
	assert.Nil(t, target.AwayTeamStat.PlayerStats)
}

func TestEntityTransformer_TransformEvent_SkipsPreseason(t *testing.T) {
	sb := loadFixture[ScoreboardResponse](t, "scoreboard_20261010_preseason.json")
	require.NotEmpty(t, sb.Events)
	for _, ev := range sb.Events {
		_, ok, err := NewEntityTransformer().TransformEvent(ev)
		require.NoError(t, err)
		assert.False(t, ok, "event %s", ev.Id)
	}
}

func TestEntityTransformer_TransformEvent_SkipsAllStar(t *testing.T) {
	sb := loadFixture[ScoreboardResponse](t, "scoreboard_20260215_allstar.json")
	require.NotEmpty(t, sb.Events)
	for _, ev := range sb.Events {
		_, ok, err := NewEntityTransformer().TransformEvent(ev)
		require.NoError(t, err)
		assert.False(t, ok, "event %s", ev.Id)
	}
}

func TestEntityTransformer_TransformEvent_SkipsUnfinished(t *testing.T) {
	ev := manualEvent(seasonTypeRegular, 4)
	ev.Competitions[0].Status.Type.Completed = false

	_, ok, err := NewEntityTransformer().TransformEvent(ev)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestEntityTransformer_TransformEvent_KeepsPlayoffsAndPlayIn(t *testing.T) {
	for _, seasonType := range []int{seasonTypePostseason, seasonTypePlayIn} {
		_, ok, err := NewEntityTransformer().TransformEvent(manualEvent(seasonType, 4))
		require.NoError(t, err)
		assert.True(t, ok, "season type %d", seasonType)
	}
}

func TestEntityTransformer_TransformEvent_Overtime(t *testing.T) {
	for period, want := range map[int]int{5: 53, 6: 58} {
		game, ok, err := NewEntityTransformer().TransformEvent(manualEvent(seasonTypeRegular, period))
		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, want, game.GameModel.Duration, "period %d", period)
	}
}

func TestEntityTransformer_TransformEvent_ParsesDateWithSeconds(t *testing.T) {
	ev := manualEvent(seasonTypeRegular, 4)
	ev.Date = "2026-10-29T00:00:00Z"

	game, ok, err := NewEntityTransformer().TransformEvent(ev)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC), game.GameModel.ScheduledAt)
}

func TestEntityTransformer_TransformEvent_ReturnsErrorOnBadScore(t *testing.T) {
	ev := manualEvent(seasonTypeRegular, 4)
	ev.Competitions[0].Competitors[0].Score = "x"

	_, _, err := NewEntityTransformer().TransformEvent(ev)
	require.ErrorContains(t, err, "event 1")
}

func TestEntityTransformer_TransformEvent_ReturnsErrorOnBadDate(t *testing.T) {
	ev := manualEvent(seasonTypeRegular, 4)
	ev.Date = "yesterday"

	_, _, err := NewEntityTransformer().TransformEvent(ev)
	require.ErrorContains(t, err, "event 1")
	require.ErrorContains(t, err, "yesterday")
}

func TestEntityTransformer_TransformEvent_ReturnsErrorWhenHomeMissing(t *testing.T) {
	ev := manualEvent(seasonTypeRegular, 4)
	ev.Competitions[0].Competitors[0].HomeAway = "away"

	_, _, err := NewEntityTransformer().TransformEvent(ev)
	require.ErrorContains(t, err, "home or away competitor missing")
}

func TestEntityTransformer_MapBoxScore_Fixture(t *testing.T) {
	game := fixtureGameWithBoxScore(t)

	assert.Len(t, game.AwayTeamStat.PlayerStats, 11)
	assert.Len(t, game.HomeTeamStat.PlayerStats, 10)

	p, ok := findPlayer(t, game.AwayTeamStat.PlayerStats, "5175737")
	require.True(t, ok)
	assert.Equal(t, "Yaxel Lendeborg", p.PlayerModel.FullName)
	s := p.GameTeamPlayerStatModel
	assert.Equal(t, 2100, s.PlayedSeconds)
	assert.Equal(t, uint8(26), s.Points)
	assert.Equal(t, uint8(9), s.Rebounds)
	assert.Equal(t, uint8(1), s.Assists)
	assert.Equal(t, uint8(7), s.Turnovers)
	assert.Equal(t, uint8(1), s.Steals)
	assert.Equal(t, uint8(1), s.Blocks)
	assert.InDelta(t, 0.55, s.FieldGoalsPercentage, 1e-4)
	assert.Equal(t, int8(2), s.PlsMin)

	for _, team := range []teams.TeamStatEntity{game.HomeTeamStat, game.AwayTeamStat} {
		for _, pl := range team.PlayerStats {
			assert.NotEqual(t, "Gui Santos", pl.PlayerModel.FullName)
		}
	}

	var hyland *players.PlayerStatisticEntity
	for i, pl := range game.HomeTeamStat.PlayerStats {
		if pl.PlayerModel.FullName == "Bones Hyland" {
			hyland = &game.HomeTeamStat.PlayerStats[i]
		}
	}
	for i, pl := range game.AwayTeamStat.PlayerStats {
		if pl.PlayerModel.FullName == "Bones Hyland" {
			hyland = &game.AwayTeamStat.PlayerStats[i]
		}
	}
	require.NotNil(t, hyland)
	assert.Equal(t, int8(0), hyland.GameTeamPlayerStatModel.PlsMin)
}

func manualGame() games.GameStatEntity {
	return games.GameStatEntity{HomeTeamExternalId: 16, AwayTeamExternalId: 9}
}

func manualSummary(labels []string, stats []string) SummaryResponse {
	block := func(id string) BoxscorePlayersEntity {
		return BoxscorePlayersEntity{
			Team: TeamEntity{Id: id},
			Statistics: []PlayersStatisticsEntity{{
				Labels: labels,
				Athletes: []AthleteStatsEntity{{
					Athlete: AthleteEntity{Id: "7", DisplayName: "Test Player"},
					Stats:   stats,
				}},
			}},
		}
	}
	return SummaryResponse{Boxscore: BoxscoreEntity{Players: []BoxscorePlayersEntity{block("16"), block("9")}}}
}

func TestEntityTransformer_MapBoxScore_LabelsOrderIndependent(t *testing.T) {
	labels := []string{"+/-", "BLK", "STL", "TO", "AST", "REB", "FG", "PTS", "MIN"}
	stats := []string{"-5", "1", "2", "3", "4", "5", "4-8", "10", "20"}

	game := manualGame()
	require.NoError(t, NewEntityTransformer().MapBoxScore(manualSummary(labels, stats), &game))

	require.Len(t, game.HomeTeamStat.PlayerStats, 1)
	s := game.HomeTeamStat.PlayerStats[0].GameTeamPlayerStatModel
	assert.Equal(t, 1200, s.PlayedSeconds)
	assert.Equal(t, uint8(10), s.Points)
	assert.Equal(t, uint8(5), s.Rebounds)
	assert.Equal(t, uint8(4), s.Assists)
	assert.Equal(t, uint8(3), s.Turnovers)
	assert.Equal(t, uint8(2), s.Steals)
	assert.Equal(t, uint8(1), s.Blocks)
	assert.InDelta(t, 0.5, s.FieldGoalsPercentage, 1e-4)
	assert.Equal(t, int8(-5), s.PlsMin)
}

func TestEntityTransformer_MapBoxScore_ReturnsErrorOnMissingLabel(t *testing.T) {
	labels := []string{"MIN", "PTS", "REB", "AST", "TO", "STL", "BLK", "FG"}
	stats := []string{"20", "10", "5", "4", "3", "2", "1", "4-8"}

	game := manualGame()
	err := NewEntityTransformer().MapBoxScore(manualSummary(labels, stats), &game)
	require.ErrorContains(t, err, "+/-")
}

func TestEntityTransformer_MapBoxScore_ReturnsErrorOnStatsLengthMismatch(t *testing.T) {
	labels := []string{"MIN", "PTS", "REB", "AST", "TO", "STL", "BLK", "FG", "+/-"}

	game := manualGame()
	err := NewEntityTransformer().MapBoxScore(manualSummary(labels, []string{"20", "10"}), &game)
	require.ErrorContains(t, err, "stats length 2 != labels 9")
}

func TestEntityTransformer_MapBoxScore_ReturnsErrorOnUnknownTeam(t *testing.T) {
	labels := []string{"MIN", "PTS", "REB", "AST", "TO", "STL", "BLK", "FG", "+/-"}
	stats := []string{"20", "10", "5", "4", "3", "2", "1", "4-8", "+2"}
	summary := manualSummary(labels, stats)
	summary.Boxscore.Players[0].Team.Id = "99"

	game := manualGame()
	err := NewEntityTransformer().MapBoxScore(summary, &game)
	require.ErrorContains(t, err, "99")
}

func TestEntityTransformer_MapBoxScore_ReturnsErrorWhenTeamMissing(t *testing.T) {
	labels := []string{"MIN", "PTS", "REB", "AST", "TO", "STL", "BLK", "FG", "+/-"}
	stats := []string{"20", "10", "5", "4", "3", "2", "1", "4-8", "+2"}
	summary := manualSummary(labels, stats)
	summary.Boxscore.Players = summary.Boxscore.Players[:1]

	game := manualGame()
	err := NewEntityTransformer().MapBoxScore(summary, &game)
	require.ErrorContains(t, err, "missing")
}

func TestEntityTransformer_VerifyTotals_Fixture(t *testing.T) {
	game := fixtureGameWithBoxScore(t)
	require.NoError(t, NewEntityTransformer().VerifyTotals(game))
}

func TestEntityTransformer_VerifyTotals_ReturnsErrorOnPointsMismatch(t *testing.T) {
	game := fixtureGameWithBoxScore(t)
	kept := []players.PlayerStatisticEntity{}
	for _, p := range game.AwayTeamStat.PlayerStats {
		if p.PlayerExternalId != "5175737" {
			kept = append(kept, p)
		}
	}
	game.AwayTeamStat.PlayerStats = kept

	err := NewEntityTransformer().VerifyTotals(game)
	require.ErrorContains(t, err, "GSW")
	require.ErrorContains(t, err, "93")
	require.ErrorContains(t, err, "119")
}

func TestEntityTransformer_VerifyTotals_ReturnsErrorOnMinutesOutOfTolerance(t *testing.T) {
	game := fixtureGameWithBoxScore(t)
	game.HomeTeamStat.PlayerStats[0].GameTeamPlayerStatModel.PlayedSeconds += 20 * 60

	err := NewEntityTransformer().VerifyTotals(game)
	require.ErrorContains(t, err, "expected 240")
}

func TestEntityTransformer_VerifyTotals_ReturnsErrorOnMinutesShortfall(t *testing.T) {
	// Выпавший из протокола игрок без очков заметен только по нехватке минут.
	game := fixtureGameWithBoxScore(t)
	game.HomeTeamStat.PlayerStats[0].GameTeamPlayerStatModel.PlayedSeconds -= 20 * 60

	err := NewEntityTransformer().VerifyTotals(game)
	require.ErrorContains(t, err, "expected 240")
}

func TestEntityTransformer_VerifyTotals_MinutesToleranceBoundary(t *testing.T) {
	// Нечётное число игроков различает допуск ⌈n/2⌉ от n/2: при 5 игравших допуск 3.
	build := func(shortfall int) games.GameStatEntity {
		game := manualGame()
		game.GameModel.Duration = 48
		game.HomeTeamStat.TeamModel.Alias = "MIN"
		game.AwayTeamStat.TeamModel.Alias = "GSW"
		for i := 0; i < 5; i++ {
			minutes := 48
			if i == 0 {
				minutes -= shortfall
			}
			stat := players.PlayerStatisticEntity{GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{PlayedSeconds: minutes * 60, Points: 20}}
			game.HomeTeamStat.PlayerStats = append(game.HomeTeamStat.PlayerStats, stat)
			game.AwayTeamStat.PlayerStats = append(game.AwayTeamStat.PlayerStats, players.PlayerStatisticEntity{GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{PlayedSeconds: 48 * 60, Points: 20}})
		}
		game.HomeTeamStat.GameTeamStatModel.Score = 100
		game.AwayTeamStat.GameTeamStatModel.Score = 100
		return game
	}

	require.NoError(t, NewEntityTransformer().VerifyTotals(build(3)))
	require.ErrorContains(t, NewEntityTransformer().VerifyTotals(build(4)), "expected 240 ± 3")
}

func TestEntityTransformer_VerifyTotals_ReturnsErrorOnNoPlayers(t *testing.T) {
	game := fixtureGame(t)

	err := NewEntityTransformer().VerifyTotals(game)
	require.ErrorContains(t, err, "no players with stats")
}

func TestEntityTransformer_MapAthleteBio_Fixture(t *testing.T) {
	tr := NewEntityTransformer()

	bio, err := tr.MapAthleteBio(loadFixture[AthleteResponse](t, "core_athlete_5175737.json"))
	require.NoError(t, err)
	assert.Equal(t, "Yaxel Lendeborg", bio.FullName)
	assert.Equal(t, time.Date(2002, 9, 30, 0, 0, 0, 0, time.UTC), bio.BirthDate)

	bio, err = tr.MapAthleteBio(loadFixture[AthleteResponse](t, "core_athlete_5279133.json"))
	require.NoError(t, err)
	assert.Equal(t, "Joan Beringer", bio.FullName)
	assert.Equal(t, time.Date(2006, 11, 11, 0, 0, 0, 0, time.UTC), bio.BirthDate)
}

func TestEntityTransformer_MapBoxScore_ErrorBranches(t *testing.T) {
	labels := []string{"MIN", "PTS", "REB", "AST", "TO", "STL", "BLK", "FG", "+/-"}
	good := []string{"20", "10", "5", "4", "3", "2", "1", "4-8", "+2"}

	tests := []struct {
		name    string
		mutate  func(s *SummaryResponse)
		wantErr string
	}{
		{"non-numeric stat", func(s *SummaryResponse) { s.Boxscore.Players[0].Statistics[0].Athletes[0].Stats[1] = "ten" }, `PTS "ten"`},
		{"malformed FG", func(s *SummaryResponse) { s.Boxscore.Players[0].Statistics[0].Athletes[0].Stats[7] = "4/8" }, `FG "4/8"`},
		{"non-numeric team id", func(s *SummaryResponse) { s.Boxscore.Players[0].Team.Id = "x" }, `team id "x"`},
		{"team without statistics", func(s *SummaryResponse) { s.Boxscore.Players[0].Statistics = nil }, "has no statistics"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			summary := manualSummary(labels, append([]string(nil), good...))
			tc.mutate(&summary)
			game := manualGame()
			require.ErrorContains(t, NewEntityTransformer().MapBoxScore(summary, &game), tc.wantErr)
		})
	}
}

func TestEntityTransformer_MapBoxScore_SkipsDidNotPlayEvenWithStats(t *testing.T) {
	labels := []string{"MIN", "PTS", "REB", "AST", "TO", "STL", "BLK", "FG", "+/-"}
	summary := manualSummary(labels, []string{"20", "10", "5", "4", "3", "2", "1", "4-8", "+2"})
	summary.Boxscore.Players[0].Statistics[0].Athletes[0].DidNotPlay = true

	game := manualGame()
	require.NoError(t, NewEntityTransformer().MapBoxScore(summary, &game))
	assert.Empty(t, game.HomeTeamStat.PlayerStats)
	assert.Len(t, game.AwayTeamStat.PlayerStats, 1)
}

func TestEntityTransformer_MapBoxScore_SkipsBenchedWithDashMinutes(t *testing.T) {
	labels := []string{"MIN", "PTS", "REB", "AST", "TO", "STL", "BLK", "FG", "+/-"}
	summary := manualSummary(labels, []string{"--", "0", "0", "0", "0", "0", "0", "0-0", "0"})

	game := manualGame()
	require.NoError(t, NewEntityTransformer().MapBoxScore(summary, &game))
	assert.Empty(t, game.HomeTeamStat.PlayerStats)
	assert.Empty(t, game.AwayTeamStat.PlayerStats)
}

// BOS - NYK 23.10.2026: Pacome Dadiet (5211983) пришёл с didNotPlay=false и минутами "--".
func TestEntityTransformer_MapBoxScore_BenchedWithDashMinutesFixture(t *testing.T) {
	game := games.GameStatEntity{HomeTeamExternalId: 2, AwayTeamExternalId: 18}
	game.GameModel.Duration = 48
	game.HomeTeamStat.GameTeamStatModel.Score = 107
	game.AwayTeamStat.GameTeamStatModel.Score = 103
	summary := loadFixture[SummaryResponse](t, "summary_401909095.json")

	transformer := NewEntityTransformer()
	require.NoError(t, transformer.MapBoxScore(summary, &game))
	require.NoError(t, transformer.VerifyTotals(game))

	_, found := findPlayer(t, game.AwayTeamStat.PlayerStats, "5211983")
	assert.False(t, found)
}

func TestEntityTransformer_TransformEvent_ReturnsErrorOnEmptyId(t *testing.T) {
	event := manualEvent(seasonTypeRegular, 4)
	event.Id = ""

	_, _, err := NewEntityTransformer().TransformEvent(event)
	require.ErrorContains(t, err, "without id")
}

func TestEntityTransformer_MapAthleteBio_ReturnsErrorOnEmptyName(t *testing.T) {
	_, err := NewEntityTransformer().MapAthleteBio(AthleteResponse{Id: "1", DateOfBirth: "2002-09-30T07:00Z"})
	require.ErrorContains(t, err, "displayName is empty")
}

func TestEntityTransformer_MapAthleteBio_ReturnsErrorOnEmptyDate(t *testing.T) {
	_, err := NewEntityTransformer().MapAthleteBio(AthleteResponse{Id: "1", DisplayName: "X"})
	require.ErrorContains(t, err, "dateOfBirth is empty")
}

func TestEntityTransformer_MapAthleteBio_ReturnsErrorOnInvalidDate(t *testing.T) {
	_, err := NewEntityTransformer().MapAthleteBio(AthleteResponse{Id: "1", DisplayName: "X", DateOfBirth: "not-a-date!"})
	require.ErrorContains(t, err, "not-a-date!")
}

func TestTeamCode(t *testing.T) {
	for in, want := range map[string]string{"GS": "GSW", "UTAH": "UTA", "MIN": "MIN", "": ""} {
		assert.Equal(t, want, teamCode(in), "input %q", in)
	}
}
