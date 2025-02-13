package cmd

import (
	gameDomain "IMP/app/internal/modules/games/domain"
	gameModel "IMP/app/internal/modules/games/domain/models"
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
	exists, err := gameDomain.NewRepository().Exists(gameModel.Game{OfficialId: gameId})
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		fmt.Println("Game with official_id " + gameId + " already exists")
		return
	}

	leagueProvider := statistics.NewLeagueProvider(enums.FromString(leagueName))
	model, err := leagueProvider.GameBoxScore(gameId)
	if err != nil {
		log.Fatal(err)
	}

	persistence := statistics.NewPersistence()
	err = persistence.SaveGameBoxScore(model)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Game with id " + gameId + " was saved into db")
}
