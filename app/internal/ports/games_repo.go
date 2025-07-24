package ports

import "IMP/app/internal/core/games"

type GamesRepo interface {
	FindOrCreateGame(model games.GameModel) (games.GameModel, error)
}
