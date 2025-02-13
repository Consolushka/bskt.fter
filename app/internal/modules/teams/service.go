package teams

import (
	"IMP/app/internal/modules/games"
	gamesModels "IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/modules/imp/domain/enums"
	impModels "IMP/app/internal/modules/imp/domain/models"
	teamModels "IMP/app/internal/modules/teams/models"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repository *Repository

	logger *logrus.Logger
}

func NewService() *Service {
	return &Service{
		repository: NewRepository(),
	}
}

func (s *Service) GetTeams() ([]teamModels.Team, error) {
	var teams []teamModels.Team

	tx := s.repository.dbConnection.Model(&teamModels.Team{}).
		Preload("League").
		Find(&teams)

	return teams, tx.Error
}

func (s *Service) GetTeamById(id int) (teamModels.Team, error) {
	var team teamModels.Team

	tx := s.repository.dbConnection.Model(&teamModels.Team{}).
		Preload("League").
		Where("id = ?", id).
		First(&team)

	return team, tx.Error
}

func (s *Service) GetTeamWithGames(teamId int) ([]*gamesModels.Game, error) {
	var gamesIds []int
	var gamesModels []*gamesModels.Game

	gamesService := games.NewService()

	tx := s.repository.dbConnection.Model(&teamModels.TeamGameStats{}).
		Where("team_id = ?", teamId).
		Select("game_id").
		Find(&gamesIds)

	for _, gameId := range gamesIds {
		game, err := gamesService.GetGame(gameId)
		if err != nil {
			return nil, err
		}

		gamesModels = append(gamesModels, game)
	}

	return gamesModels, tx.Error
}

func (s *Service) GetTeamWithGamesMetrics(teamId int, impPers []enums.ImpPERs) ([]*impModels.GameImpMetrics, error) {
	var gamesIds []int
	var gamesModels []*impModels.GameImpMetrics

	gamesService := games.NewService()

	tx := s.repository.dbConnection.Model(&teamModels.TeamGameStats{}).
		Where("team_id = ?", teamId).
		Select("game_id").
		Find(&gamesIds)

	for _, gameId := range gamesIds {
		game, err := gamesService.GetGameMetrics(gameId, impPers)
		if err != nil {
			return nil, err
		}

		gamesModels = append(gamesModels, game)
	}

	return gamesModels, tx.Error
}
