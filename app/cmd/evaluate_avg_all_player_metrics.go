package cmd

import (
	"IMP/app/internal/modules/players/domain"
	"IMP/app/log"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var evaluateAllPlayersAvgMetrics = &cobra.Command{
	Use:   "evaluate-avg-all-player-metrics",
	Short: "seed database with existing seeders",
	Run: func(cmd *cobra.Command, args []string) {
		evaluateAllPlayersAvgMetricsMeth()
	},
}

func init() {
	rootCmd.AddCommand(evaluateAllPlayersAvgMetrics)
}

// seed math seeders for given models and then seed model
func evaluateAllPlayersAvgMetricsMeth() {
	playersIds, err := domain.NewRepository().ListOfPlayersIds()
	if err != nil {
		log.Error(err)
		return
	}
	for _, id := range playersIds {
		EvaluatePlayerAvgMetrics(id)
		fmt.Println("Saved AVG Metrics for player with id: " + strconv.Itoa(id))
	}
}
