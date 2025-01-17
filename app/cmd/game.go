package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var generateGameCmd = &cobra.Command{
	Use:   "generate:game",
	Short: "Generate game results",
	Long:  "Generate file with IMP indicators for every player played in the given game",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		GamePdf(args[0])
		fmt.Println("Game results file successfully generated")
	},
}

func init() {
	rootCmd.AddCommand(generateGameCmd)
}

// GamePdf takes gameId and generates file with IMP indicators for each player played in the given game
func GamePdf(gameId string) {
	//gameService := games.NewService()
	//gameModel, err := gameService.GetGame(gameId)

	//repo := statistics.NewLeagueRepository(enums.FromString(leagueName))
	//model, _ := repo.GameBoxScore(gameId)
	//
	//gameRes := queries.CalculateFullGame(model)
	//pdfcommands.PrintGame(gameRes, nil)
}
