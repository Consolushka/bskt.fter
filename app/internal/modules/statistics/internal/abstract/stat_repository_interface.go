package abstract

import (
	"IMP/app/internal/modules/statistics/models"
	"time"
)

type StatsProvider interface {
	// GameBoxScore returns boxscore data from stats provider
	GameBoxScore(gameId string) (*models.GameBoxScoreDTO, error)
	// GamesByDate returns list of games for given date
	GamesByDate(date time.Time) ([]string, error)
	// GamesByTeam returns list of already played games for given team
	GamesByTeam(teamId string) ([]string, error)
}
