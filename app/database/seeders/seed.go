package seeders

import (
	"IMP/app/database/seeders/leagues"
	"IMP/app/internal/abstract/database"
	"reflect"
	"strings"
)

var AvailableModels = []database.Seeder{
	leagues.Seeder{},
}

func FindSeeder(model string) *database.Seeder {
	for _, seeder := range AvailableModels {
		if strings.ToLower(reflect.TypeOf(seeder.Model()).Name()) == strings.ToLower(model) {
			return &seeder
		}
	}

	return nil
}
