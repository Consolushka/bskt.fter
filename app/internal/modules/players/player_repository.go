package players

import (
	"IMP/app/database"
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

func (r *Repository) FirstByOfficialId(id string) (*Player, error) {
	var result Player

	tx := r.dbConnection.
		First(
			&result,
			Player{
				OfficialId: id,
			})

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &result, tx.Error
}

func (r *Repository) FirstOrCreate(player Player) (*Player, error) {
	var result Player

	tx := r.dbConnection.
		Attrs(Player{
			FullNameLocal: player.FullNameLocal,
			FullNameEn:    player.FullNameEn,
			BirthDate:     player.BirthDate,
		}).
		FirstOrCreate(
			&result,
			Player{
				OfficialId: player.OfficialId,
			})

	return &result, tx.Error
}

func (r *Repository) FirstOrCreateGameStat(stats PlayerGameStats) error {
	tx := r.dbConnection.Attrs(
		PlayerGameStats{
			PlayedSeconds: stats.PlayedSeconds,
			PlsMin:        stats.PlsMin,
			IsBench:       stats.IsBench,
		}).
		FirstOrCreate(
			&PlayerGameStats{},
			PlayerGameStats{
				PlayerID:   stats.PlayerID,
				TeamGameId: stats.TeamGameId,
			})

	return tx.Error
}
