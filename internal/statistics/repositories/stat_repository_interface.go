package repositories

import (
	"NBATrueEfficency/internal/statistics/repositories/sport_radar/dtos"
)

type StatsRepository interface {
	// GetGame returns game data from stats provider
	GetGame(gameId string) (*dtos.GameDTO, error)
}
