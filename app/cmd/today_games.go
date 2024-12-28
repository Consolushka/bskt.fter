package cmd

import (
	"IMP/app/internal/modules/imp/queries"
	pdfcommands "IMP/app/internal/modules/pdf/commands"
	"IMP/app/internal/modules/statistics/leagues/nba/repositories_factory"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var generateTodayGamesCmd = &cobra.Command{
	Use:   "generate:today",
	Short: "Generate pdf files for every today games",
	Long:  "Generate files with IMP indicator for every todays games for each player who played in the given game",
	Run: func(cmd *cobra.Command, args []string) {
		TodayGamesPdf()
	},
}

func init() {
	rootCmd.AddCommand(generateTodayGamesCmd)
}

// TodayGamesPdf fetches id's of today games. And then generates pdf files for each game
func TodayGamesPdf() {
	repo := repositories_factory.NewNbaStatsRepository()
	date, gamesIds, _ := repo.TodayGames()

	for _, gameId := range gamesIds {
		game, err := repo.GameBoxScore(gameId)
		if err != nil {
			log.Fatal(err)
			return
		}

		gameRes := queries.CalculateFullGame(game)
		pdfcommands.PrintGame(gameRes, &date)

		fmt.Println(gameRes.Title + " Game results file successfully generated")
	}
}
