package leagues

import (
	"IMP/app/database"
	leaguesModels "IMP/app/internal/modules/leagues/domain/models"
	"strings"
)

type Seeder struct {
}

func (l Seeder) Model() interface{} {
	return leaguesModels.League{}
}

func (l Seeder) Seed() {
	dbConnection := database.OpenDbConnection()

	nbaLeague := &leaguesModels.League{
		NameLocal:        "National Basketball Association",
		AliasLocal:       strings.ToUpper(leaguesModels.NBAAlias),
		NameEn:           "National Basketball Association",
		AliasEn:          strings.ToUpper(leaguesModels.NBAAlias),
		PeriodsNumber:    4,
		PeriodDuration:   12,
		OvertimeDuration: 6,
	}

	mlblLeague := &leaguesModels.League{
		NameLocal:        "Межрегиональная любительская баскетбольная лига",
		AliasLocal:       "МЛБЛ",
		NameEn:           "Interregional Amateur Basketball League",
		AliasEn:          strings.ToUpper(leaguesModels.MLBLAlias),
		PeriodsNumber:    4,
		PeriodDuration:   10,
		OvertimeDuration: 5,
	}

	dbConnection.FirstOrCreate(nbaLeague, leaguesModels.League{AliasEn: strings.ToUpper(leaguesModels.NBAAlias)})
	dbConnection.FirstOrCreate(mlblLeague, leaguesModels.League{AliasEn: strings.ToUpper(leaguesModels.MLBLAlias)})
}
