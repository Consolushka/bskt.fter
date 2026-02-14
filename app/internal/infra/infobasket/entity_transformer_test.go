package infobasket

import (
	"testing"
	"time"
)

func TestEntityTransformer_Transform(t *testing.T) {
	transformer := EntityTransformer{}

	game := GameBoxScoreResponse{
		GameDate:    "14.02.2026",
		GameTimeMsk: "21.30",
		MaxPeriod:   6,
		GameTeams: []TeamBoxScoreDto{
			{
				TeamName: TeamNameBoxScoreDto{
					CompTeamAbcNameEn:    "HOME",
					CompTeamNameEn:       "Home Team",
					CompTeamRegionNameEn: "Home City",
				},
				Score: 100,
				Players: []PlayerBoxScoreDto{
					{
						PersonID:     10,
						PersonNameEn: "New Player",
						PersonNameRu: "Новый Игрок",
						PersonBirth:  "01.01.2000",
						Seconds:      900,
						PlusMinus:    7,
					},
					{
						PersonID:     11,
						PersonNameEn: "Bad Date",
						PersonBirth:  "31-12-2000",
						Seconds:      100,
						PlusMinus:    -2,
					},
				},
			},
			{
				TeamName: TeamNameBoxScoreDto{
					CompTeamAbcNameEn:    "AWAY",
					CompTeamNameEn:       "Away Team",
					CompTeamRegionNameEn: "Away City",
				},
				Score: 94,
				Players: []PlayerBoxScoreDto{
					{
						PersonID:     20,
						PersonNameEn: "Away Player",
						PersonBirth:  "02.02.2001",
						Seconds:      850,
						PlusMinus:    -3,
					},
				},
			},
		},
	}

	entity, err := transformer.Transform(game)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if entity.GameModel.Title != "HOME - AWAY" {
		t.Fatalf("unexpected title: %s", entity.GameModel.Title)
	}
	// 5 periods * 10 + 1 overtime * 5
	if entity.GameModel.Duration != 55 {
		t.Fatalf("unexpected duration: %d", entity.GameModel.Duration)
	}

	expectedScheduled, _ := time.Parse("02.01.2006 15.04", "14.02.2026 21.30")
	if !entity.GameModel.ScheduledAt.Equal(expectedScheduled) {
		t.Fatalf("unexpected scheduled time: %v", entity.GameModel.ScheduledAt)
	}

	if entity.HomeTeamStat.TeamModel.Name != "Home Team" {
		t.Fatalf("unexpected home team name: %s", entity.HomeTeamStat.TeamModel.Name)
	}
	if entity.HomeTeamStat.GameTeamStatModel.Score != 100 {
		t.Fatalf("unexpected home score: %d", entity.HomeTeamStat.GameTeamStatModel.Score)
	}
	if entity.HomeTeamStat.GameTeamStatModel.FinalDiff != 6 {
		t.Fatalf("unexpected home final diff: %d", entity.HomeTeamStat.GameTeamStatModel.FinalDiff)
	}
	// transformTeam preallocates len(team.Players), so invalid player leaves zero-value slot.
	if len(entity.HomeTeamStat.PlayerStats) != 2 {
		t.Fatalf("unexpected home players count: %d", len(entity.HomeTeamStat.PlayerStats))
	}
	if entity.HomeTeamStat.PlayerStats[0].PlayerModel.FullName != "Новый Игрок" {
		t.Fatalf("expected fallback russian name, got %s", entity.HomeTeamStat.PlayerStats[0].PlayerModel.FullName)
	}
	if entity.HomeTeamStat.PlayerStats[1].PlayerExternalId != "" {
		t.Fatalf("expected zero-value player for invalid date, got external id %s", entity.HomeTeamStat.PlayerStats[1].PlayerExternalId)
	}

	if entity.AwayTeamStat.TeamModel.Name != "Away Team" {
		t.Fatalf("unexpected away team name: %s", entity.AwayTeamStat.TeamModel.Name)
	}
	if entity.AwayTeamStat.GameTeamStatModel.Score != 94 {
		t.Fatalf("unexpected away score: %d", entity.AwayTeamStat.GameTeamStatModel.Score)
	}
	if entity.AwayTeamStat.GameTeamStatModel.FinalDiff != -6 {
		t.Fatalf("unexpected away final diff: %d", entity.AwayTeamStat.GameTeamStatModel.FinalDiff)
	}
}

func TestEntityTransformer_Transform_InvalidGameDateTime(t *testing.T) {
	transformer := EntityTransformer{}

	game := GameBoxScoreResponse{
		GameDate:    "bad-date",
		GameTimeMsk: "21.30",
		GameTeams: []TeamBoxScoreDto{
			{TeamName: TeamNameBoxScoreDto{CompTeamAbcNameEn: "A"}},
			{TeamName: TeamNameBoxScoreDto{CompTeamAbcNameEn: "B"}},
		},
	}

	_, err := transformer.Transform(game)
	if err == nil {
		t.Fatal("expected error for invalid game date/time")
	}
}

func TestPlayersTrans(t *testing.T) {
	t.Run("success_regular_name", func(t *testing.T) {
		player := PlayerBoxScoreDto{
			PersonID:     1,
			PersonNameEn: "John Doe",
			PersonNameRu: "Джон Доу",
			PersonBirth:  "03.03.2003",
			Seconds:      777,
			PlusMinus:    9,
		}

		result, err := playersTrans(player)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.PlayerExternalId != "1" {
			t.Fatalf("unexpected external id: %s", result.PlayerExternalId)
		}
		if result.PlayerModel.FullName != "John Doe" {
			t.Fatalf("unexpected full name: %s", result.PlayerModel.FullName)
		}
		if result.GameTeamPlayerStatModel.PlayedSeconds != 777 {
			t.Fatalf("unexpected seconds: %d", result.GameTeamPlayerStatModel.PlayedSeconds)
		}
		if result.GameTeamPlayerStatModel.PlsMin != 9 {
			t.Fatalf("unexpected plus-minus: %d", result.GameTeamPlayerStatModel.PlsMin)
		}
	})

	t.Run("success_new_player_name_fallback", func(t *testing.T) {
		player := PlayerBoxScoreDto{
			PersonID:     2,
			PersonNameEn: "New Player",
			PersonNameRu: "Новый Игрок",
			PersonBirth:  "04.04.2004",
		}

		result, err := playersTrans(player)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.PlayerModel.FullName != "Новый Игрок" {
			t.Fatalf("expected russian fallback name, got %s", result.PlayerModel.FullName)
		}
	})

	t.Run("invalid_birth_date", func(t *testing.T) {
		player := PlayerBoxScoreDto{
			PersonID:     3,
			PersonNameEn: "Bad Birth",
			PersonBirth:  "2004-04-04",
		}

		_, err := playersTrans(player)
		if err == nil {
			t.Fatal("expected error for invalid birth date")
		}
	})
}
