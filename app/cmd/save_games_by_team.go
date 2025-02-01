package cmd

import (
	"IMP/app/internal/modules/statistics"
	"IMP/app/internal/modules/statistics/enums"
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
	leagueProvider := statistics.NewLeagueProvider(enums.FromString(leagueName))

	gameIds, err := leagueProvider.GamesByTeam(teamId)
	if err != nil {
		panic(err)
	}

	for _, gameId := range gameIds {
		SaveGame(leagueName, gameId)
	}
}
