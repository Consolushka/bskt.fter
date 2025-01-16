package leagues

import (
	"IMP/app/database"
	"IMP/app/internal/modules/leagues"
)

type Seeder struct {
}

func (l Seeder) Model() interface{} {
	return leagues.League{}
}

func (l Seeder) Seed() {
	dbConnection := database.OpenDbConnection()

	nbaLeague := &leagues.League{
		NameLocal:        "National Basketball Association",
		AliasLocal:       "NBA",
		NameEn:           "National Basketball Association",
		AliasEn:          "NBA",
		PeriodsNumber:    4,
		PeriodDuration:   12,
		OvertimeDuration: 6,
	}

	mlblLeague := &leagues.League{
		NameLocal:        "Межрегиональная любительская баскетбольная лига",
		AliasLocal:       "МЛБЛ",
		NameEn:           "Interregional Amateur Basketball League",
		AliasEn:          "MLBL",
		PeriodsNumber:    4,
		PeriodDuration:   10,
		OvertimeDuration: 5,
	}

	dbConnection.FirstOrCreate(nbaLeague, leagues.League{AliasEn: "NBA"})
	dbConnection.FirstOrCreate(mlblLeague, leagues.League{AliasEn: "MLBL"})
}
