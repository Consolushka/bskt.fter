package responses

import (
	gamesModels "IMP/app/internal/modules/games/domain/models"
	gamesResources "IMP/app/internal/modules/games/domain/resources"
	"IMP/app/internal/utils/array_utils"
)

type PlayerGamesBoxScoreResponse struct {
	Games []gamesResources.Game
}

func NewPlayerGamesBoxScoreResponse(games []gamesModels.Game) PlayerGamesBoxScoreResponse {
	return PlayerGamesBoxScoreResponse{
		Games: array_utils.Map(games, func(game gamesModels.Game) gamesResources.Game {
			return gamesResources.NewGameResource(game)
		}),
	}

}
