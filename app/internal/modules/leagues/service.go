package leagues

import (
	gamesDomain "IMP/app/internal/modules/games/domain"
	gamesModels "IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/modules/leagues/domain"
	leaguesModels "IMP/app/internal/modules/leagues/domain/models"
	"time"
)

type Service struct {
	leaguesRepository *domain.Repository
	gamesRepository   *gamesDomain.Repository
}

func NewService() *Service {
	return &Service{
		leaguesRepository: domain.NewRepository(),
		gamesRepository:   gamesDomain.NewRepository(),
	}
}

func (s *Service) GetAllLeagues() ([]leaguesModels.League, error) {
	return s.leaguesRepository.List()
}

func (s *Service) GetGamesByLeagueAndDate(leagueId int, date time.Time) ([]gamesModels.Game, error) {
	return s.gamesRepository.GamesStatsByDateList(date, &leagueId)
}
