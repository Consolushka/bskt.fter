package players

import (
	gamesDomain "IMP/app/internal/modules/games/domain"
	gamesModels "IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/modules/players/domain"
	"IMP/app/internal/modules/players/domain/models"
)

type Service struct {
	repository *domain.Repository

	gamesRepository *gamesDomain.Repository
}

func NewService() *Service {
	return &Service{
		repository:      domain.NewRepository(),
		gamesRepository: gamesDomain.NewRepository(),
	}
}

func (s *Service) GetPlayerByFullName(fullName string) ([]models.Player, error) {
	return s.repository.ListByFullName(fullName)
}

func (s *Service) GetPlayerGamesBoxScore(playerId int) ([]gamesModels.Game, error) {
	playerGames, err := s.repository.ListOfGamesByPlayerId(playerId)
	if err != nil {
		return nil, err
	}

	return s.gamesRepository.ListOfGamesByGamesIdsAndPlayerId(playerGames, playerId)
}
