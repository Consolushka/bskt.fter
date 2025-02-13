package domain

import (
	"IMP/app/database"
	"IMP/app/internal/modules/players/domain/models"
	"errors"
	"gorm.io/gorm"
)

type Repository struct {
	dbConnection *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		dbConnection: database.GetDB(),
	}
}

func (r *Repository) FirstByOfficialId(id string) (*models.Player, error) {
	var result models.Player

	tx := r.dbConnection.
		First(
			&result,
			models.Player{
				OfficialId: id,
			})

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &result, tx.Error
}

func (r *Repository) FirstOrCreate(player models.Player) (*models.Player, error) {
	var result models.Player

	tx := r.dbConnection.
		Attrs(models.Player{
			FullNameLocal: player.FullNameLocal,
			FullNameEn:    player.FullNameEn,
			BirthDate:     player.BirthDate,
		}).
		FirstOrCreate(
			&result,
			models.Player{
				OfficialId: player.OfficialId,
			})

	return &result, tx.Error
}

func (r *Repository) FirstOrCreateGameStat(stats models.PlayerGameStats) error {
	tx := r.dbConnection.Attrs(
		models.PlayerGameStats{
			PlayedSeconds: stats.PlayedSeconds,
			PlsMin:        stats.PlsMin,
			IsBench:       stats.IsBench,
		}).
		FirstOrCreate(
			&models.PlayerGameStats{},
			models.PlayerGameStats{
				PlayerID:   stats.PlayerID,
				TeamGameId: stats.TeamGameId,
			})

	return tx.Error
}
