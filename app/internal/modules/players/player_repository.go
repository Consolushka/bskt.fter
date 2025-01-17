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

func (r *Repository) FirstByLeaguePlayerId(id int) (*Player, error) {
	var result Player

	tx := r.dbConnection.
		First(
			&result,
			Player{
				LeaguePlayerID: id,
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
			FullName:  player.FullName,
			BirthDate: player.BirthDate,
		}).
		FirstOrCreate(
			&result,
			Player{
				LeaguePlayerID: player.LeaguePlayerID,
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
