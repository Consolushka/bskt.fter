package ports

import "IMP/app/internal/core/games"

type GamesRepo interface {
	GameExists(model games.GameModel) (bool, error)
	FindOrCreateGame(model games.GameModel) (games.GameModel, error)
}
