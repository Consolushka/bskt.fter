package ports

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"time"
)

type StatsProvider interface {
	GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error)
	// EnrichGameStats обязан вернуть игру с теми же командами и названием, что получил:
	// к этому моменту они уже взяты из базы, и потеря команды не даст сохранить игру.
	EnrichGameStats(game games.GameStatEntity) (games.GameStatEntity, error)
	GetPlayerBio(id string) (players.PlayerBioEntity, error)
}
