package seeders

import (
	"IMP/app/database/seeders/abstract"
	"IMP/app/database/seeders/leagues"
	"reflect"
	"strings"
)

var AvailableModels = []abstract.Seeder{
	leagues.Seeder{},
}

func FindSeeder(model string) *abstract.Seeder {
	for _, seeder := range AvailableModels {
		if strings.ToLower(reflect.TypeOf(seeder.Model()).Name()) == strings.ToLower(model) {
			return &seeder
		}
	}

	return nil
}
