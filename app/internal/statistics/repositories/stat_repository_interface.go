package repositories

import (
	"FTER/app/internal/models"
)

type StatsRepository interface {
	// GameBoxScore returns game data from stats provider
	GameBoxScore(gameId string) (*models.GameModel, error)
}
