package cmd

import (
	"IMP/app/internal/modules/statistics"
	"IMP/app/internal/modules/statistics/enums"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var saveGameCmd = &cobra.Command{
	Use:   "save-game",
	Short: "Saves game into application",
	Long:  "Saves game results into database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SaveGame(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(saveGameCmd)
}

func SaveGame(leagueName string, gameId string) {
	repo := statistics.NewLeagueProvider(enums.FromString(leagueName))
	model, err := repo.GameBoxScore(gameId)
	if err != nil {
		log.Fatal(err)
	}

	persistence := statistics.NewPersistence()
	err = persistence.SaveGameBoxScore(model)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Game results file successfully generated")
}
