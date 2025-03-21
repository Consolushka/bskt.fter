package cmd

import (
	leaguesDomain "IMP/app/internal/modules/leagues/domain"
	"IMP/app/internal/modules/statistics"
	"IMP/app/log"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"time"
)

var saveGameByDateCmd = &cobra.Command{
	Use:   "save-games-by-date",
	Short: "Saves games into application by date",
	Long:  "Fetch games ids from league schedule and saves them into database",
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
	leagueRepository := leaguesDomain.NewRepository()
	league, err := leagueRepository.FirstByAliasEn(leagueName)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	leagueProvider := statistics.NewLeagueProvider(league.AliasEn)

	gameIds, err := leagueProvider.GamesByDate(date)
	log.Info("On date " + date.Format("02-01-2006") + " there are " + strconv.Itoa(len(gameIds)) + " games in " + strings.ToUpper(leagueName))
	if err != nil {
		panic(err)
	}

	for _, gameId := range gameIds {
		SaveGame(leagueName, gameId)
	}
}
