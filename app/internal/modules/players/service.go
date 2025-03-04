package players

import (
	"IMP/app/internal/modules/games"
	gamesDomain "IMP/app/internal/modules/games/domain"
	gamesModels "IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/modules/imp/domain/enums"
	impModels "IMP/app/internal/modules/imp/domain/models"
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

func (s *Service) GetPlayerGamesMetrics(playerId int, impPers []enums.ImpPERs) ([]*impModels.GameImpMetrics, error) {
	playerGames, err := s.GetPlayerGamesBoxScore(playerId)
	if err != nil {
		return nil, err
	}

	gamesService := games.NewService()

	return gamesService.GetGamesMetrics(playerGames, impPers)
}
