package cmd

import (
	"IMP/app/internal/domain"
	"IMP/app/internal/persistence"
	statistics2 "IMP/app/internal/statistics"
	"IMP/app/log"
	"fmt"
	"github.com/spf13/cobra"
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
	logger := log.NewLogger()

	leagueRepository := domain.NewLeaguesRepository()
	league, err := leagueRepository.FirstByAliasEn(leagueName)
	if err != nil {
		logger.Fatalln(err)
	}

	exists, err := domain.NewGamesRepository().Exists(domain.Game{OfficialId: gameId})
	if err != nil {
		logger.Fatalln(err)
	}
	if exists {
		fmt.Println("Game with official_id " + gameId + " already exists")
		return
	}

	leagueProvider := statistics2.NewLeagueProvider(league.AliasEn)
	model, err := leagueProvider.GameBoxScore(gameId)
	if err != nil {
		logger.Fatalln(err)
	}
	if !model.IsFinal {
		logger.Fatalln("Game with Id" + gameId + " is not final")
	}

	service := persistence.NewService()
	err = service.SaveGameBoxScore(model)
	if err != nil {
		logger.Fatalln(err)
	}

	fmt.Println("Game with id " + gameId + " was saved into db")
}
