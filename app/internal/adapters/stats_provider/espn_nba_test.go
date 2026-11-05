package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/infra/espn"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAdapter(t *testing.T) (EspnNbaStatsProviderAdapter, *espn.MockClientInterface) {
	t.Helper()

	ctrl := gomock.NewController(t)
	client := espn.NewMockClientInterface(ctrl)
	adapter, err := NewEspnNbaStatsProviderAdapter(client)
	require.NoError(t, err)

	return adapter, client
}

func finishedEvent(id, dateUTC, homeAbbr, awayAbbr string, homeScore, awayScore int) espn.EventEntity {
	return espn.EventEntity{
		Id:     id,
		Date:   dateUTC,
		Season: espn.SeasonEntity{Type: 2},
		Competitions: []espn.CompetitionEntity{{
			Type:   espn.CompetitionTypeEntity{Abbreviation: "STD"},
			Status: espn.StatusEntity{Type: espn.StatusTypeEntity{Completed: true}, Period: 4},
			Competitors: []espn.CompetitorEntity{
				{HomeAway: "home", Score: strconv.Itoa(homeScore), Team: espn.TeamEntity{Id: "1", Abbreviation: homeAbbr}},
				{HomeAway: "away", Score: strconv.Itoa(awayScore), Team: espn.TeamEntity{Id: "2", Abbreviation: awayAbbr}},
			},
		}},
	}
}

func TestEspnNbaStatsProviderAdapter_GetGamesStatsByPeriod(t *testing.T) {
	t.Run("queries one ESPN day for a window inside one US Eastern date", func(t *testing.T) {
		adapter, client := newAdapter(t)
		client.EXPECT().Scoreboard("20261028").Times(1).Return(espn.ScoreboardResponse{
			Events: []espn.EventEntity{finishedEvent("1001", "2026-10-29T00:00Z", "MIN", "GS", 100, 90)},
		}, nil)

		result, err := adapter.GetGamesStatsByPeriod(
			time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 10, 29, 3, 0, 0, 0, time.UTC),
		)

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "1001", result[0].ExternalGameId)
	})

	t.Run("queries each US Eastern day in the window", func(t *testing.T) {
		adapter, client := newAdapter(t)
		client.EXPECT().Scoreboard("20261028").Return(espn.ScoreboardResponse{
			Events: []espn.EventEntity{finishedEvent("1001", "2026-10-29T00:00Z", "MIN", "GS", 100, 90)},
		}, nil)
		client.EXPECT().Scoreboard("20261029").Return(espn.ScoreboardResponse{
			Events: []espn.EventEntity{finishedEvent("1002", "2026-10-30T00:00Z", "MIN", "GS", 100, 90)},
		}, nil)

		result, err := adapter.GetGamesStatsByPeriod(
			time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC),
			time.Date(2026, 10, 30, 2, 0, 0, 0, time.UTC),
		)

		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("skips malformed event and keeps the rest", func(t *testing.T) {
		adapter, client := newAdapter(t)
		broken := finishedEvent("1001", "2026-10-29T00:00Z", "MIN", "GS", 100, 90)
		broken.Competitions[0].Competitors[0].Score = "x"
		client.EXPECT().Scoreboard("20261028").Return(espn.ScoreboardResponse{
			Events: []espn.EventEntity{broken, finishedEvent("1002", "2026-10-29T00:00Z", "MIN", "GS", 100, 90)},
		}, nil)

		result, err := adapter.GetGamesStatsByPeriod(
			time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 10, 29, 3, 0, 0, 0, time.UTC),
		)

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "1002", result[0].ExternalGameId)
	})

	t.Run("returns error when scoreboard fails", func(t *testing.T) {
		adapter, client := newAdapter(t)
		client.EXPECT().Scoreboard("20261028").Return(espn.ScoreboardResponse{}, nil)
		client.EXPECT().Scoreboard("20261029").Return(espn.ScoreboardResponse{}, errors.New("api error"))

		result, err := adapter.GetGamesStatsByPeriod(
			time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC),
			time.Date(2026, 10, 30, 2, 0, 0, 0, time.UTC),
		)

		require.ErrorContains(t, err, "20261029")
		require.ErrorContains(t, err, "api error")
		assert.Nil(t, result)
	})

	t.Run("returns empty slice when there are no events", func(t *testing.T) {
		adapter, client := newAdapter(t)
		client.EXPECT().Scoreboard("20261028").Return(espn.ScoreboardResponse{}, nil)

		result, err := adapter.GetGamesStatsByPeriod(
			time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 10, 29, 3, 0, 0, 0, time.UTC),
		)

		require.NoError(t, err)
		assert.Empty(t, result)
		assert.NotNil(t, result)
	})
}

func boxScore(homePoints1, homePoints2 int) espn.SummaryResponse {
	labels := []string{"MIN", "PTS", "FG", "REB", "AST", "TO", "STL", "BLK", "+/-"}
	athlete := func(id string, points int) espn.AthleteStatsEntity {
		return espn.AthleteStatsEntity{
			Athlete: espn.AthleteEntity{Id: id, DisplayName: "Player " + id},
			Stats:   []string{"120", strconv.Itoa(points), "5-10", "1", "1", "1", "1", "1", "0"},
		}
	}
	block := func(teamId string, athletes ...espn.AthleteStatsEntity) espn.BoxscorePlayersEntity {
		return espn.BoxscorePlayersEntity{
			Team:       espn.TeamEntity{Id: teamId},
			Statistics: []espn.PlayersStatisticsEntity{{Labels: labels, Athletes: athletes}},
		}
	}

	return espn.SummaryResponse{Boxscore: espn.BoxscoreEntity{Players: []espn.BoxscorePlayersEntity{
		block("1", athlete("11", homePoints1), athlete("12", homePoints2)),
		block("2", athlete("21", 45), athlete("22", 45)),
	}}}
}

func TestEspnNbaStatsProviderAdapter_EnrichGameStats(t *testing.T) {
	newGame := func(t *testing.T) games.GameStatEntity {
		t.Helper()
		game, ok, err := espn.NewEntityTransformer().TransformEvent(finishedEvent("1001", "2026-10-29T00:00Z", "MIN", "GS", 100, 90))
		require.NoError(t, err)
		require.True(t, ok)
		return game
	}

	t.Run("returns game as is when ExternalGameId is empty", func(t *testing.T) {
		adapter, _ := newAdapter(t)
		game := games.GameStatEntity{}

		result, err := adapter.EnrichGameStats(game)

		require.NoError(t, err)
		assert.Equal(t, game, result)
	})

	t.Run("maps players and verifies totals", func(t *testing.T) {
		adapter, client := newAdapter(t)
		client.EXPECT().Summary("1001").Return(boxScore(50, 50), nil)

		result, err := adapter.EnrichGameStats(newGame(t))

		require.NoError(t, err)
		require.Len(t, result.HomeTeamStat.PlayerStats, 2)
		require.Len(t, result.AwayTeamStat.PlayerStats, 2)
		assert.Equal(t, uint8(50), result.HomeTeamStat.PlayerStats[0].GameTeamPlayerStatModel.Points)
	})

	t.Run("returns error when summary fails", func(t *testing.T) {
		adapter, client := newAdapter(t)
		game := newGame(t)
		client.EXPECT().Summary("1001").Return(espn.SummaryResponse{}, errors.New("api error"))

		result, err := adapter.EnrichGameStats(game)

		require.ErrorContains(t, err, "ESPN event 1001")
		assert.Equal(t, game.GameModel, result.GameModel)
	})

	t.Run("returns error when box score does not map", func(t *testing.T) {
		adapter, client := newAdapter(t)
		game := newGame(t)
		summary := boxScore(50, 50)
		summary.Boxscore.Players[0].Team.Id = "99"
		client.EXPECT().Summary("1001").Return(summary, nil)

		result, err := adapter.EnrichGameStats(game)

		require.ErrorContains(t, err, "ESPN event 1001")
		require.ErrorContains(t, err, "transform")
		assert.Equal(t, game.GameModel, result.GameModel)
	})

	t.Run("returns error with event id when totals mismatch", func(t *testing.T) {
		adapter, client := newAdapter(t)
		game := newGame(t)
		client.EXPECT().Summary("1001").Return(boxScore(50, 40), nil)

		result, err := adapter.EnrichGameStats(game)

		require.ErrorContains(t, err, "ESPN event 1001")
		require.ErrorContains(t, err, "verify")
		assert.Equal(t, game.GameModel, result.GameModel)
	})
}

func TestEspnNbaStatsProviderAdapter_GetPlayerBio(t *testing.T) {
	t.Run("returns bio from core athlete", func(t *testing.T) {
		adapter, client := newAdapter(t)
		client.EXPECT().Athlete("5175737").Return(espn.AthleteResponse{
			Id: "5175737", DisplayName: "Yaxel Lendeborg", DateOfBirth: "2002-09-30T07:00Z",
		}, nil)

		bio, err := adapter.GetPlayerBio("5175737")

		require.NoError(t, err)
		assert.Equal(t, "Yaxel Lendeborg", bio.FullName)
		assert.Equal(t, time.Date(2002, 9, 30, 0, 0, 0, 0, time.UTC), bio.BirthDate)
	})

	t.Run("returns error when athlete request fails", func(t *testing.T) {
		adapter, client := newAdapter(t)
		client.EXPECT().Athlete("5175737").Return(espn.AthleteResponse{}, errors.New("api error"))

		_, err := adapter.GetPlayerBio("5175737")

		require.ErrorContains(t, err, "ESPN athlete 5175737")
	})

	t.Run("returns error when date of birth is empty", func(t *testing.T) {
		adapter, client := newAdapter(t)
		client.EXPECT().Athlete("5175737").Return(espn.AthleteResponse{Id: "5175737", DisplayName: "X"}, nil)

		_, err := adapter.GetPlayerBio("5175737")

		require.ErrorContains(t, err, "map returned error")
		require.ErrorContains(t, err, "dateOfBirth is empty")
	})
}
