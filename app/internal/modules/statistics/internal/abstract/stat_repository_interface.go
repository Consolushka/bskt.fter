package abstract

import (
	"IMP/app/internal/modules/statistics/models"
)

type StatsProvider interface {
	// GameBoxScore returns boxscore data from stats provider
	GameBoxScore(gameId string) (*models.GameBoxScoreDTO, error)
	// TodayGames returns date in string format and id's of games
	TodayGames() (string, []string, error)
}
