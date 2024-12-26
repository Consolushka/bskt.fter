package cmd

import (
	"FTER/app/internal/enums"
	ftercommands "FTER/app/internal/fter/commands"
	pdfcommands "FTER/app/internal/pdf/commands"
	"FTER/app/internal/statistics/factory"
	"fmt"
	"github.com/spf13/cobra"
)

var generateGameCmd = &cobra.Command{
	Use:   "generate:game",
	Short: "Generate game results",
	Long:  "Generate file with FTER indicator for every player played in the given game",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		GamePdf(args[0], args[0])
		fmt.Println("Game results file successfully generated")
	},
}

func init() {
	rootCmd.AddCommand(generateGameCmd)
}

// GamePdf takes sportRadar gameId and generates pdf with players FTER
func GamePdf(leagueName string, gameId string) {
	//gameId := "2aa29340-f4ca-4e43-be10-02a7415eece4"
	//repo, err := factories.NewNbaStatsRepository()
	repo := factory.NewLeagueRepository(enums.FromString(leagueName))
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	model, _ := repo.GameBoxScore(gameId)

	gameRes := ftercommands.CalculateFullGame(model)
	pdfcommands.PrintGame(gameRes, nil)
}
