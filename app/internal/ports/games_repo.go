package ports

import "IMP/app/internal/core/games"

type GamesRepo interface {
	Exists(model games.GameModel) (bool, error)
	FirstOrCreate(model games.GameModel) (games.GameModel, error)
}
