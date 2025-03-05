package cmd

import (
	gamesModels "IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/modules/imp"
	"IMP/app/internal/modules/players"
	playersDomain "IMP/app/internal/modules/players/domain"
	"IMP/app/internal/modules/players/domain/models"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/log"
	"github.com/spf13/cobra"
	"strconv"
)

var evaluatePlayerAvgMetricsCmd = &cobra.Command{
	Use:   "evaluate-avg-player-metrics",
	Short: "Saves games into application by team",
	Long:  "Fetch game ids for team and saves them into application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		intVal, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error(err)
			return
		}
		EvaluatePlayerAvgMetrics(intVal)
	},
}

func init() {
	rootCmd.AddCommand(evaluatePlayerAvgMetricsCmd)
}

// playerGameData contains player stats and game score differential
type playerGameData struct {
	Game              gamesModels.Game
	PlayerStats       models.PlayerGameStats
	ScoreDifferential int // Positive if player's team won, negative if lost
	IsHomeTeam        bool
}

func EvaluatePlayerAvgMetrics(playerId int) {
	playersRepository := playersDomain.NewRepository()

	sumClearImp := 0.0
	sumPlayedSeconds := 0
	countFromBench := 0

	service := players.NewService()

	games, _ := service.GetPlayerGamesBoxScore(playerId)
	gamesPlayed := len(games)

	gameData := array_utils.Map(games, func(game gamesModels.Game) playerGameData {
		return getPlayerGameStats(game, playerId)
	})

	for _, gameData := range gameData {
		sumPlayedSeconds += gameData.PlayerStats.PlayedSeconds
		sumClearImp += imp.EvaluateClean(gameData.PlayerStats.PlayedSeconds, gameData.PlayerStats.PlsMin, gameData.ScoreDifferential, gameData.Game.PlayedMinutes)

		if gameData.PlayerStats.IsBench {
			countFromBench++
		}
	}

	playerMetric := models.PlayerMetrics{
		PlayerID:         playerId,
		Player:           nil,
		AvgClearImp:      sumClearImp / float64(gamesPlayed),
		AvgPlayedSeconds: float64(sumPlayedSeconds) / float64(gamesPlayed),
		PlayedGamesCount: gamesPlayed,
		FromBenchCount:   countFromBench,
		FromStartCount:   gamesPlayed - countFromBench,
	}

	err := playersRepository.CreatePlayerMetric(playerMetric)
	if err != nil {
		log.Fatalln(err)
	}
}

func getPlayerGameStats(game gamesModels.Game, playerId int) playerGameData {
	playerAwayStats := array_utils.Filter(game.AwayTeamStats.PlayerGameStats, func(stats models.PlayerGameStats) bool {
		return stats.PlayerID == playerId
	})

	// Check if player is in home team
	playerHomeStats := array_utils.Filter(game.HomeTeamStats.PlayerGameStats, func(stats models.PlayerGameStats) bool {
		return stats.PlayerID == playerId
	})

	var playerStats models.PlayerGameStats
	var scoreDiff int
	var isHomeTeam bool

	if len(playerAwayStats) > 0 {
		// Player is in away team
		playerStats = playerAwayStats[0]
		scoreDiff = game.AwayTeamStats.Points - game.HomeTeamStats.Points
		isHomeTeam = false
	} else if len(playerHomeStats) > 0 {
		// Player is in home team
		playerStats = playerHomeStats[0]
		scoreDiff = game.HomeTeamStats.Points - game.AwayTeamStats.Points
		isHomeTeam = true
	}

	return playerGameData{
		Game:              game,
		PlayerStats:       playerStats,
		ScoreDifferential: scoreDiff,
		IsHomeTeam:        isHomeTeam,
	}
}
