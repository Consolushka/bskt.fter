package ports

import "IMP/app/internal/core/players"

type PlayersRepo interface {
	FirstOrCreatePlayer(player players.PlayerModel) (players.PlayerModel, error)
	PlayersByFullName(fullName string) ([]players.PlayerModel, error)
	FirstOrCreatePlayerStat(playerStat players.GameTeamPlayerStatModel) (players.GameTeamPlayerStatModel, error)
}
