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

func (g Gorm) PlayersByFullName(fullName string) ([]players.PlayerModel, error) {
	var playerModels []players.PlayerModel

	tx := g.db.Model(&players.PlayerModel{}).Where(&playerModels, players.PlayerModel{
		FullName: fullName,
	})

	return playerModels, tx.Error
}

func (g Gorm) FirstOrCreatePlayerStat(playerStat players.GameTeamPlayerStatModel) (players.GameTeamPlayerStatModel, error) {
	tx := g.db.FirstOrCreate(&playerStat, players.GameTeamPlayerStatModel{
		GameId:   playerStat.GameId,
		TeamId:   playerStat.TeamId,
		PlayerId: playerStat.PlayerId,
	})

	return playerStat, tx.Error
}

func NewGormRepo(db *gorm.DB) Gorm {
	return Gorm{db: db}
}
