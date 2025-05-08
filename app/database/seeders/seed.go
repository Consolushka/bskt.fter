package seeders

import (
	"IMP/app/database/seeders/leagues"
	"reflect"
	"strings"
)

var AvailableModels = []Seeder{
	leagues.Seeder{},
}

func FindSeeder(model string) *Seeder {
	for _, seeder := range AvailableModels {
		if strings.ToLower(reflect.TypeOf(seeder.Model()).Name()) == strings.ToLower(model) {
			return &seeder
		}
	}

	return nil
}
