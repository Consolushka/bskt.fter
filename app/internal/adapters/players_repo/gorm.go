package players_repo

import (
	"IMP/app/internal/core/players"
	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func (g Gorm) FirstOrCreatePlayer(player players.PlayerModel) (players.PlayerModel, error) {
	tx := g.db.FirstOrCreate(&player, players.PlayerModel{
		FullName:  player.FullName,
		BirthDate: player.BirthDate,
	})

	return player, tx.Error
}

func (g Gorm) FirstOrCreatePlayerStat(playerStat players.GameTeamPlayerStatModel) (players.GameTeamPlayerStatModel, error) {
	tx := g.db.FirstOrCreate(&playerStat, players.GameTeamPlayerStatModel{
		GameTeamId: playerStat.GameTeamId,
		PlayerId:   playerStat.PlayerId,
	})

	return playerStat, tx.Error
}

func NewGormRepo(db *gorm.DB) Gorm {
	return Gorm{db: db}
}
