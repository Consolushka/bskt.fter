package leagues

import (
	"IMP/app/database"
	"IMP/app/internal/domain"
	"strings"
)

type Seeder struct {
}

func (l Seeder) Model() interface{} {
	return domain.League{}
}

func (l Seeder) Seed() {
	dbConnection := database.OpenDbConnection()

	nbaLeague := &domain.League{
		NameLocal:        "National Basketball Association",
		AliasLocal:       strings.ToUpper(domain.NBAAlias),
		NameEn:           "National Basketball Association",
		AliasEn:          strings.ToUpper(domain.NBAAlias),
		PeriodsNumber:    4,
		PeriodDuration:   12,
		OvertimeDuration: 6,
	}

	mlblLeague := &domain.League{
		NameLocal:        "Межрегиональная любительская баскетбольная лига",
		AliasLocal:       "МЛБЛ",
		NameEn:           "Interregional Amateur Basketball League",
		AliasEn:          strings.ToUpper(domain.MLBLAlias),
		PeriodsNumber:    4,
		PeriodDuration:   10,
		OvertimeDuration: 5,
	}

	dbConnection.FirstOrCreate(nbaLeague, domain.League{AliasEn: strings.ToUpper(domain.NBAAlias)})
	dbConnection.FirstOrCreate(mlblLeague, domain.League{AliasEn: strings.ToUpper(domain.MLBLAlias)})
}
