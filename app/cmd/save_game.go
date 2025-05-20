package cmd

import (
	"IMP/app/internal/domain"
	"IMP/app/internal/persistence"
	"IMP/app/internal/statistics"
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
		fmt.Println(err)
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

	leagueProvider, err := statistics.NewLeagueProvider(league)
	if err != nil {
		fmt.Println(err)
		logger.Fatalln(err)
	}
	model, err := leagueProvider.GameBoxScore(gameId)
	if err != nil {
		logger.Fatalln(err)
	}
	if !model.IsFinal {
		message := "Game with Id" + gameId + " is not final"
		fmt.Println(message)
		logger.Fatalln(message)
	}

	service := persistence.NewService()
	err = service.SaveGameBoxScore(model)
	if err != nil {
		fmt.Println(err)
		logger.Fatalln(err)
	}

	fmt.Println("Game with id " + gameId + " was saved into db")
}
