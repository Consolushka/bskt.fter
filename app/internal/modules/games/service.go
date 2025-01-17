package games

import (
	"IMP/app/database"
	"IMP/app/internal/modules/imp/models"
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
		Preload("HomeTeamStats").
		Preload("HomeTeamStats.Team").
		Preload("HomeTeamStats.PlayerGameStats").
		Preload("HomeTeamStats.PlayerGameStats.Player").
		Preload("AwayTeamStats").
		Preload("AwayTeamStats.Team").
		Preload("AwayTeamStats.PlayerGameStats").
		Preload("AwayTeamStats.PlayerGameStats.Player").
		First(&gameModel, GameModel{ID: id})

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &gameModel, nil
}

func (s *Service) GetGameMetrics(id int) (*models.GameImpMetrics, error) {
	var gameModel GameModel

	tx := s.dbConnection.Debug().
		Preload("League").
		Preload("HomeTeamStats").
		Preload("HomeTeamStats.PlayerGameStats", "player_game_stats.game_id = ?", id).
		Preload("HomeTeamStats.PlayerGameStats.Player").
		Preload("AwayTeamStats").
		Preload("AwayTeamStats.PlayerGameStats", "player_game_stats.game_id = ?", id).
		Preload("AwayTeamStats.PlayerGameStats.Player").
		First(&gameModel, GameModel{ID: id})

	if tx.Error != nil {
		return nil, tx.Error
	}

	gameImpMetrics := models.GameImpMetrics{
		Id:        gameModel.ID,
		Scheduled: &gameModel.ScheduledAt,
		Home: models.TeamImpMetrics{
			Alias:       "",
			TotalPoints: 0,
			Players:     nil,
		},
		Away:         models.TeamImpMetrics{},
		FullGameTime: gameModel.PlayedMinutes,
	}

	return &gameImpMetrics, nil
}
