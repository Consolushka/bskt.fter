package cmd

import (
	"IMP/app/internal/modules/statistics"
	"IMP/app/internal/modules/statistics/enums"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var saveGameByDateCmd = &cobra.Command{
	Use:   "save-game-by-date",
	Short: "Saves game into application",
	Long:  "Saves game results into database",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var dateString string
		if len(args) == 1 {
			dateString = time.Now().Format("02-01-2006")
		} else {
			dateString = args[1]
		}

		date, err := time.Parse("02-01-2006", dateString)
		if err != nil {
			fmt.Println("Incorrect date format. Please use format: dd-mm-yyyy")
		}
		SaveGameByDate(args[0], date)
		fmt.Println("Game results file successfully generated")
	},
}

func init() {
	rootCmd.AddCommand(saveGameByDateCmd)
}

func SaveGameByDate(leagueName string, date time.Time) {
	nbaProvider := statistics.NewLeagueProvider(enums.FromString(leagueName))

	gameIds, err := nbaProvider.GamesByDate(date)
	if err != nil {
		panic(err)
	}

	for _, gameId := range gameIds {
		SaveGame(leagueName, gameId)
	}
}
