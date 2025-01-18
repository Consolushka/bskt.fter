package games

import (
	"IMP/app/database"
	"IMP/app/internal/modules/imp"
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/modules/players"
	"IMP/app/internal/modules/teams"
	"IMP/app/internal/utils/array_utils"
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
	gameModel, err := s.retrieveGameModelById(id)

	if err != nil {
		return nil, err
	}

	return gameModel, nil
}

// GetGameMetrics returns game metrics by specific id
//
// Calculates IMP metrics for every player
func (s *Service) GetGameMetrics(id int) (*models.GameImpMetrics, error) {
	gameModel, err := s.retrieveGameModelById(id)

	if err != nil {
		return nil, err
	}

	gameImpMetrics := s.mapGameModelToImpMetricsModel(gameModel)

	return gameImpMetrics, nil
}

func (s *Service) retrieveGameModelById(id int) (*GameModel, error) {
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

	return &gameModel, tx.Error
}

func (s *Service) mapGameModelToImpMetricsModel(gameModel *GameModel) *models.GameImpMetrics {
	return &models.GameImpMetrics{
		Id:        gameModel.ID,
		Scheduled: &gameModel.ScheduledAt,
		Home: models.TeamImpMetrics{
			Alias:       gameModel.HomeTeamStats.Team.Alias,
			TotalPoints: gameModel.HomeTeamStats.Points,
			Players:     s.mapTeamPlayersMetrics(gameModel.HomeTeamStats, gameModel.AwayTeamStats, gameModel.PlayedMinutes),
		},
		Away: models.TeamImpMetrics{
			Alias:       gameModel.AwayTeamStats.Team.Alias,
			TotalPoints: gameModel.AwayTeamStats.Points,
			Players:     s.mapTeamPlayersMetrics(gameModel.AwayTeamStats, gameModel.HomeTeamStats, gameModel.PlayedMinutes),
		},
		FullGameTime: gameModel.PlayedMinutes,
	}
}

func (s *Service) mapTeamPlayersMetrics(currentTeam teams.TeamGameStats, oposingTeam teams.TeamGameStats, playedMinutes int) []models.PlayerImpMetrics {
	return array_utils.Map(currentTeam.PlayerGameStats, func(playerGameStats players.PlayerGameStats) models.PlayerImpMetrics {
		return models.PlayerImpMetrics{
			FullName:      playerGameStats.Player.FullName,
			SecondsPlayed: playerGameStats.PlayedSeconds,
			PlsMin:        playerGameStats.PlsMin,
			IMP:           imp.CalculatePlayerImpPerMinute(float64(playerGameStats.PlayedSeconds)/60, playerGameStats.PlsMin, currentTeam.Points-oposingTeam.Points, playedMinutes),
		}
	})
}
