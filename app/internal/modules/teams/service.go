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
	return s.repository.List()
}

func (s *Service) GetTeamById(id int) (teamsModels.Team, error) {
	return s.repository.FirstById(id)
}

func (s *Service) GetTeamWithGames(teamId int) ([]*gamesModels.Game, error) {
	var gamesCollection []*gamesModels.Game

	gamesService := games.NewService()

	gamesIds, err := s.repository.TeamGameIdListByTeamId(teamId)
	if err != nil {
		return nil, err
	}

	for _, gameId := range gamesIds {
		game, err := gamesService.GetGame(gameId)
		if err != nil {
			return nil, err
		}

		gamesCollection = append(gamesCollection, game)
	}

	return gamesCollection, nil
}

func (s *Service) GetTeamWithGamesMetrics(teamId int, impPers []enums.ImpPERs) ([]*impModels.GameImpMetrics, error) {
	var gamesCollection []*impModels.GameImpMetrics

	gamesService := games.NewService()

	gamesIds, err := s.repository.TeamGameIdListByTeamId(teamId)
	if err != nil {
		return nil, err
	}

	for _, gameId := range gamesIds {
		game, err := gamesService.GetGameMetrics(gameId, impPers)
		if err != nil {
			return nil, err
		}

		gamesCollection = append(gamesCollection, game)
	}

	return gamesCollection, nil
}
