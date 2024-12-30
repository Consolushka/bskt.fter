package leagues

import (
	"IMP/app/database"
	"IMP/app/internal/modules/leagues/models"
)

type Seeder struct {
}

func (l Seeder) Model() interface{} {
	return models.League{}
}

func (l Seeder) Seed() {
	dbConnection := database.Connect()

	nbaLeague := &models.League{
		NameLocal:        "National Basketball Association",
		AliasLocal:       "NBA",
		NameEn:           "National Basketball Association",
		AliasEn:          "NBA",
		PeriodsNumber:    4,
		PeriodDuration:   12,
		OvertimeDuration: 6,
	}

	mlblLeague := &models.League{
		NameLocal:        "Межрегиональная любительская баскетбольная лига",
		AliasLocal:       "МЛБЛ",
		NameEn:           "Interregional Amateur Basketball League",
		AliasEn:          "MLBL",
		PeriodsNumber:    4,
		PeriodDuration:   10,
		OvertimeDuration: 5,
	}

	dbConnection.FirstOrCreate(nbaLeague, models.League{AliasEn: "NBA"})
	dbConnection.FirstOrCreate(mlblLeague, models.League{AliasEn: "MLBL"})
}
