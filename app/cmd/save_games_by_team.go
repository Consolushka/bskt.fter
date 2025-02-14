package cmd

import (
	leaguesDomain "IMP/app/internal/modules/leagues/domain"
	"IMP/app/internal/modules/statistics"
	"IMP/app/log"
	"github.com/spf13/cobra"
)

var saveGameByTeamCmd = &cobra.Command{
	Use:   "save-game-by-team",
	Short: "Saves game into application",
	Long:  "Saves game results into database",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SaveGameByTeam(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(saveGameByTeamCmd)
}

func SaveGameByTeam(leagueName string, teamId string) {
	leagueRepository := leaguesDomain.NewRepository()
	league, err := leagueRepository.GetLeagueByAliasEn(leagueName)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	leagueProvider := statistics.NewLeagueProvider(league.AliasEn)

	gameIds, err := leagueProvider.GamesByTeam(teamId)
	if err != nil {
		panic(err)
	}

	for _, gameId := range gameIds {
		SaveGame(leagueName, gameId)
	}
}
