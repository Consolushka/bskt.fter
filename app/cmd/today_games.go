package cmd

import (
	ftercommands "FTER/app/internal/fter/commands"
	pdfcommands "FTER/app/internal/pdf/commands"
	"FTER/app/internal/statistics/factories"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var generateTodayGamesCmd = &cobra.Command{
	Use:   "generate:today",
	Short: "Generate today games results",
	Long:  "Generate files with IMP indicator for every todays games for each player who played in the given game",
	Run: func(cmd *cobra.Command, args []string) {
		todayGamesPdf()
	},
}

func init() {
	rootCmd.AddCommand(generateTodayGamesCmd)
}

// todayGamesPdf fetches id's of today games. And then generates pdf files for each game
func todayGamesPdf() {
	repo, err := factories.NewStatsRepository()
	if err != nil {
		log.Fatal(err)
		return
	}
	date, gamesIds, err := repo.TodayGames()

	for _, gameId := range gamesIds {
		game, err := repo.GameBoxScore(gameId)
		if err != nil {
			log.Fatal(err)
			return
		}

		gameRes := ftercommands.CalculateFullGame(game)
		pdfcommands.PrintGame(gameRes, &date)

		fmt.Println(gameRes.Title + " Game results file successfully generated")
	}
}
