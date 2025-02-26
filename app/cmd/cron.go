package cmd

import (
	"IMP/app/internal/modules/leagues/domain/models"
	"IMP/app/log"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"time"
)

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start cron scheduler for background tasks",
	Long:  `This command starts the cron scheduler that runs periodic tasks like downloading yesterday's games.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Started cron scheduler")

		timeZone := "Europe/Moscow"
		loc, _ := time.LoadLocation(timeZone)
		cronJob := cron.New(cron.WithLocation(loc))
		cronJob.AddJob("0 12 * * *", saveYesterdayGamesJob{})

		cronJob.Start()
		select {}
	},
}

func init() {
	rootCmd.AddCommand(cronCmd)
}

type saveYesterdayGamesJob struct {
}

func (s saveYesterdayGamesJob) Run() {
	yesterday := time.Now().AddDate(0, 0, -1)
	SaveGameByDate(models.NBAAlias, yesterday)
	SaveGameByDate(models.MLBLAlias, yesterday)
}
