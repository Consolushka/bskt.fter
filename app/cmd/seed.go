package cmd

import (
	"IMP/app/db/seeders"
	"github.com/spf13/cobra"
	"strings"
)

var seedDbCmd = &cobra.Command{
	Use:   "seed",
	Short: "",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		seeders := strings.Split(args[0], ",")
		Seed(seeders)
	},
}

func init() {
	rootCmd.AddCommand(seedDbCmd)
}

// GamePdf takes gameId and generates file with IMP indicators for each player played in the given game
func Seed(seederModels []string) {
	for _, seederModel := range seederModels {
		seeders.SeedModel(seederModel)
	}
}
