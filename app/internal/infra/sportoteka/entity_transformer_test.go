package sportoteka

import (
	"testing"
	"time"
)

func intPtr(v int) *int {
	return &v
}

func TestEntityTransformer_Transform(t *testing.T) {
	transformer := EntityTransformer{}

	scheduledAt := time.Date(2026, 2, 14, 12, 30, 0, 0, time.UTC)
	game := GameBoxScoreEntity{
		Game: GameInfoEntity{
			Periods:       6,
			ScheduledTime: scheduledAt,
		},
		Team1: TeamInfoEntity{
			AbcName:    "T1",
			Name:       "Home Team",
			RegionName: "Home City",
		},
		Team2: TeamInfoEntity{
			AbcName:    "T2",
			Name:       "Away Team",
			RegionName: "Away City",
		},
		Teams: []TeamBoxScoreEntity{
			{
				TeamNumber: 2,
				Total: TeamBoxScoreTotalsEntity{
					Points: 88,
				},
				Starts: []TeamBoxScoreStartEntity{
					{
						StartRole: "Team",
					},
					{
						StartRole: "Player",
						PersonId:  intPtr(2),
						LastName:  "Away",
						FirstName: "Player",
						Birthday:  "2001-02-03T00:00:00",
						Stats: PlayerBoxScoreStatsEntity{
							Second:    600,
							PlusMinus: -5,
						},
					},
				},
			},
			{
				TeamNumber: 1,
				Total: TeamBoxScoreTotalsEntity{
					Points: 95,
				},
				Starts: []TeamBoxScoreStartEntity{
					{
						StartRole: "Player",
						PersonId:  intPtr(1),
						LastName:  "Home",
						FirstName: "Good",
						Birthday:  "2000-01-01T00:00:00",
						Stats: PlayerBoxScoreStatsEntity{
							Second:    900,
							PlusMinus: 10,
						},
					},
					{
						StartRole: "Player",
						PersonId:  intPtr(3),
						LastName:  "Home",
						FirstName: "BadBirth",
						Birthday:  "bad-date",
						Stats: PlayerBoxScoreStatsEntity{
							Second:    120,
							PlusMinus: 0,
						},
					},
				},
			},
		},
	}

	entity, err := transformer.Transform(game)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entity.GameModel.Title != "T1 - T2" {
		t.Fatalf("unexpected title: %s", entity.GameModel.Title)
	}
	if entity.GameModel.ScheduledAt != scheduledAt {
		t.Fatalf("unexpected scheduledAt: %v", entity.GameModel.ScheduledAt)
	}
	// 5 periods * 10 + 1 overtime period * 5 = 55
	if entity.GameModel.Duration != 55 {
		t.Fatalf("unexpected duration: %d", entity.GameModel.Duration)
	}

	if entity.HomeTeamStat.TeamModel.Name != "Home Team" {
		t.Fatalf("unexpected home team name: %s", entity.HomeTeamStat.TeamModel.Name)
	}
	if entity.HomeTeamStat.GameTeamStatModel.Score != 95 {
		t.Fatalf("unexpected home score: %d", entity.HomeTeamStat.GameTeamStatModel.Score)
	}
	if entity.HomeTeamStat.GameTeamStatModel.FinalDiff != 7 {
		t.Fatalf("unexpected home final diff: %d", entity.HomeTeamStat.GameTeamStatModel.FinalDiff)
	}
	// One valid player, one invalid birthday skipped.
	if len(entity.HomeTeamStat.PlayerStats) != 1 {
		t.Fatalf("unexpected home players count: %d", len(entity.HomeTeamStat.PlayerStats))
	}
	if entity.HomeTeamStat.PlayerStats[0].PlayerExternalId != "1" {
		t.Fatalf("unexpected home player external id: %s", entity.HomeTeamStat.PlayerStats[0].PlayerExternalId)
	}

	if entity.AwayTeamStat.TeamModel.Name != "Away Team" {
		t.Fatalf("unexpected away team name: %s", entity.AwayTeamStat.TeamModel.Name)
	}
	if entity.AwayTeamStat.GameTeamStatModel.Score != 88 {
		t.Fatalf("unexpected away score: %d", entity.AwayTeamStat.GameTeamStatModel.Score)
	}
	if entity.AwayTeamStat.GameTeamStatModel.FinalDiff != -7 {
		t.Fatalf("unexpected away final diff: %d", entity.AwayTeamStat.GameTeamStatModel.FinalDiff)
	}
	// One Team-role player skipped, one valid player remains.
	if len(entity.AwayTeamStat.PlayerStats) != 1 {
		t.Fatalf("unexpected away players count: %d", len(entity.AwayTeamStat.PlayerStats))
	}
	if entity.AwayTeamStat.PlayerStats[0].PlayerExternalId != "2" {
		t.Fatalf("unexpected away player external id: %s", entity.AwayTeamStat.PlayerStats[0].PlayerExternalId)
	}
}

func TestEntityTransformer_PlayerTransform(t *testing.T) {
	transformer := EntityTransformer{}

	t.Run("success_case", func(t *testing.T) {
		player := TeamBoxScoreStartEntity{
			PersonId:  intPtr(42),
			LastName:  "Doe",
			FirstName: "John",
			Birthday:  "1999-12-31T00:00:00",
			Stats: PlayerBoxScoreStatsEntity{
				Second:    321,
				PlusMinus: 7,
			},
		}

		stat, err := transformer.playerTransform(player)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stat.PlayerExternalId != "42" {
			t.Fatalf("unexpected external id: %s", stat.PlayerExternalId)
		}
		if stat.PlayerModel.FullName != "Doe John" {
			t.Fatalf("unexpected full name: %s", stat.PlayerModel.FullName)
		}
		if stat.GameTeamPlayerStatModel.PlayedSeconds != 321 {
			t.Fatalf("unexpected played seconds: %d", stat.GameTeamPlayerStatModel.PlayedSeconds)
		}
		if stat.GameTeamPlayerStatModel.PlsMin != 7 {
			t.Fatalf("unexpected plus-minus: %d", stat.GameTeamPlayerStatModel.PlsMin)
		}
	})

	t.Run("invalid_birthday", func(t *testing.T) {
		player := TeamBoxScoreStartEntity{
			PersonId:  intPtr(1),
			LastName:  "Invalid",
			FirstName: "Birth",
			Birthday:  "bad-date",
		}

		_, err := transformer.playerTransform(player)
		if err == nil {
			t.Fatal("expected error for invalid birthday")
		}
	})
}
