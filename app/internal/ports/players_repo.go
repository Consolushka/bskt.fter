package ports

import "IMP/app/internal/core/players"

type PlayersRepo interface {
	FirstOrCreate(player players.PlayerModel) (players.PlayerModel, error)
	ListByFullName(fullName string) ([]players.PlayerModel, error)
	FirstOrCreateStat(playerStat players.GameTeamPlayerStatModel) (players.GameTeamPlayerStatModel, error)
}
