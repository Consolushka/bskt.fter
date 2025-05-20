package cmd

import (
	"IMP/app/internal/domain"
	"IMP/app/internal/statistics"
	"IMP/app/log"
	"fmt"
	"github.com/spf13/cobra"
)

var saveGameByTeamCmd = &cobra.Command{
	Use:   "save-game-by-team",
	Short: "Saves games into application by team",
	Long:  "Fetch game ids for team and saves them into application",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SaveGameByTeam(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(saveGameByTeamCmd)
}

func SaveGameByTeam(leagueName string, teamId string) {
	leagueRepository := domain.NewLeaguesRepository()
	league, err := leagueRepository.FirstByAliasEn(leagueName)
	if err != nil {
		log.NewLogger().Fatalln(err)
		panic(err)
	}
	leagueProvider, err := statistics.NewLeagueProvider(league)
	if err != nil {
		fmt.Println(err)
		log.NewLogger().Fatalln(err)
	}

	gameIds, err := leagueProvider.GamesByTeam(teamId)
	if err != nil {
		panic(err)
	}

	for _, gameId := range gameIds {
		SaveGame(leagueName, gameId)
	}
}
