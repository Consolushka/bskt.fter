package leagues

import (
	gamesDomain "IMP/app/internal/modules/games/domain"
	gamesModels "IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/modules/leagues/domain"
	leaguesModels "IMP/app/internal/modules/leagues/domain/models"
	teamsDomain "IMP/app/internal/modules/teams/domain"
	teamsModels "IMP/app/internal/modules/teams/domain/models"
	"time"
)

type Service struct {
	leaguesRepository *domain.Repository
	gamesRepository   *gamesDomain.Repository
	teamsRepository   *teamsDomain.Repository
}

func NewService() *Service {
	return &Service{
		leaguesRepository: domain.NewRepository(),
		gamesRepository:   gamesDomain.NewRepository(),
		teamsRepository:   teamsDomain.NewRepository(),
	}
}

func (s *Service) GetAllLeagues() ([]leaguesModels.League, error) {
	return s.leaguesRepository.List()
}

func (s *Service) GetGamesByLeagueAndDate(leagueId int, date time.Time) ([]gamesModels.Game, error) {
	return s.gamesRepository.GamesStatsByDateList(date, &leagueId)
}

func (s *Service) GetLeague(id int) (*leaguesModels.League, error) {
	return s.leaguesRepository.FirstById(id)
}

func (s *Service) GetTeamsByLeague(leagueId int) ([]teamsModels.Team, error) {
	return s.teamsRepository.ListByLeague(leagueId)
}

func (s *Service) GetPlayersRanking(leagueId int, limit *int, minMinuterPerGame int, avgMinutes int, gamesPlayed int) (*[]leaguesModels.PlayerImpRanking, error) {
	return s.leaguesRepository.ListPlayerRankingByImpClean(leagueId, limit, minMinuterPerGame, avgMinutes, gamesPlayed)
}
