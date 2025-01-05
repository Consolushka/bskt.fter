package players

import (
	"IMP/app/database"
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
func (r *Repository) FirstOrCreate(player Player) (Player, error) {
	var result Player

	tx := r.dbConnection.
		FirstOrCreate(
			&result,
			Player{
				FullName:  player.FullName,
				BirthDate: player.BirthDate,
				DraftYear: player.DraftYear,
			})

	return result, tx.Error
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
				PlayerID: stats.PlayerID,
				GameID:   stats.GameID,
				TeamID:   stats.TeamID,
			})

	return tx.Error
}
