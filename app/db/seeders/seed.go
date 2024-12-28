package seeders

import "IMP/app/db/seeders/leagues"

func SeedModel(model string) {
	switch model {
	case "league":
		leagues.Seed()
		break
	default:
		panic("No such model")
	}
}
