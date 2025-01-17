package games

import (
	"IMP/app/database"
	"gorm.io/gorm"
)

type Service struct {
	repository   *Repository
	dbConnection *gorm.DB
}

func NewService() *Service {
	return &Service{
		repository:   NewRepository(),
		dbConnection: database.GetDB().Debug(),
	}
}

// GetGame returns game by specific id
//
// Also preloads:
//   - League
//   - Teams
//   - Players stats
//   - Players models
func (s *Service) GetGame(id int) (*GameModel, error) {
	var gameModel GameModel

	tx := s.dbConnection.Debug().
		Preload("League").
		Preload("HomeTeam").
		Preload("HomeTeam.PlayerGameStats", "player_game_stats.game_id = ?", id).
		Preload("HomeTeam.PlayerGameStats.Player").
		Preload("AwayTeam").
		Preload("AwayTeam.PlayerGameStats", "player_game_stats.game_id = ?", id).
		Preload("AwayTeam.PlayerGameStats.Player").
		First(&gameModel, GameModel{ID: id})

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &gameModel, nil
}
