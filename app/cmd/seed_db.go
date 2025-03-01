package cmd

import (
	"IMP/app/database/seeders"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var seedDbCmd = &cobra.Command{
	Use:   "seed:db",
	Short: "seed database with existing seeders",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		seeders := strings.Split(args[0], ",")
		seed(seeders)
	},
}

func init() {
	rootCmd.AddCommand(seedDbCmd)
}

// seed math seeders for given models and then seed model
func seed(seederModels []string) {
	for _, seederModel := range seederModels {
		seeder := seeders.FindSeeder(seederModel)
		if seeder == nil {
			fmt.Print("Seeder not found for model: " + seederModel)
			continue
		}
		(*seeder).Seed()
	}
}
