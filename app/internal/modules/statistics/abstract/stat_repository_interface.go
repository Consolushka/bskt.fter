package abstract

import (
	"FTER/app/internal/modules/imp/models"
)

type StatsRepository interface {
	// GameBoxScore returns boxscore data from stats provider
	GameBoxScore(gameId string) (*models.GameModel, error)
	// TodayGames returns date in string format and id's of games
	TodayGames() (string, []string, error)
}
