package cmd

import (
	"IMP/app/internal/domain"
	"IMP/app/log"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"time"
)

// cronCmd starts cron scheduler for background tasks at 12:00 AM Moscow time.
//
// If the current time is after 12:00 AM, the job will run immediately.
//
// Saves yesterday's games for each league.
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start cron scheduler for background tasks",
	Long:  `This command starts the cron scheduler that runs periodic tasks like downloading yesterday's games.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Started cron scheduler")

		timeZone := "Europe/Moscow"
		loc, _ := time.LoadLocation(timeZone)
		cronJob := cron.New(cron.WithLocation(loc))

		job := newSaveYesterdayGamesJob()
		cronJob.AddJob("0 12 * * *", job)

		now := time.Now().In(loc)
		noon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, loc)

		// Если сейчас позже 12:00, запускаем задачу немедленно
		if now.After(noon) {
			log.Info("Current time is after 12:00, running job immediately...")
			go job.Run()
		}

		// Запускаем планировщик
		cronJob.Start()

		select {}
	},
}

func init() {
	rootCmd.AddCommand(cronCmd)
}

type saveYesterdayGamesJob struct {
	leaguesRepository *domain.LeaguesRepository
}

func newSaveYesterdayGamesJob() *saveYesterdayGamesJob {
	return &saveYesterdayGamesJob{
		leaguesRepository: domain.NewLeaguesRepository(),
	}
}

func (s saveYesterdayGamesJob) Run() {
	yesterday := time.Now().AddDate(0, 0, -1)

	leagues, err := s.leaguesRepository.List()
	if err != nil {
		log.Fatalln(err)
	}

	for _, league := range leagues {
		SaveGameByDate(league.AliasEn, yesterday)
	}
}
