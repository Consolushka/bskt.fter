package teams

import (
	"IMP/app/internal/modules/games"
	gamesModels "IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/modules/imp/domain/enums"
	impModels "IMP/app/internal/modules/imp/domain/models"
	teamsDomain "IMP/app/internal/modules/teams/domain"
	teamsModels "IMP/app/internal/modules/teams/domain/models"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repository *teamsDomain.Repository

	logger *logrus.Logger
}

func NewService() *Service {
	return &Service{
		repository: teamsDomain.NewRepository(),
	}
}

func (s *Service) GetTeams() ([]teamsModels.Team, error) {
	var teams []teamsModels.Team

	tx := s.repository.DbConnection.Model(&teamsModels.Team{}).
		Preload("League").
		Find(&teams)

	return teams, tx.Error
}

func (s *Service) GetTeamById(id int) (teamsModels.Team, error) {
	var team teamsModels.Team

	tx := s.repository.DbConnection.Model(&teamsModels.Team{}).
		Preload("League").
		Where("id = ?", id).
		First(&team)

	return team, tx.Error
}

func (s *Service) GetTeamWithGames(teamId int) ([]*gamesModels.Game, error) {
	var gamesIds []int
	var gamesCollection []*gamesModels.Game

	gamesService := games.NewService()

	tx := s.repository.DbConnection.Model(&teamsModels.TeamGameStats{}).
		Where("team_id = ?", teamId).
		Select("game_id").
		Find(&gamesIds)

	for _, gameId := range gamesIds {
		game, err := gamesService.GetGame(gameId)
		if err != nil {
			return nil, err
		}

		gamesCollection = append(gamesCollection, game)
	}

	return gamesCollection, tx.Error
}

func (s *Service) GetTeamWithGamesMetrics(teamId int, impPers []enums.ImpPERs) ([]*impModels.GameImpMetrics, error) {
	var gamesIds []int
	var gamesCollection []*impModels.GameImpMetrics

	gamesService := games.NewService()

	tx := s.repository.DbConnection.Model(&teamsModels.TeamGameStats{}).
		Where("team_id = ?", teamId).
		Select("game_id").
		Find(&gamesIds)

	for _, gameId := range gamesIds {
		game, err := gamesService.GetGameMetrics(gameId, impPers)
		if err != nil {
			return nil, err
		}

		gamesCollection = append(gamesCollection, game)
	}

	return gamesCollection, tx.Error
}
