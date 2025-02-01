package games

import (
	"IMP/app/database"
	"gorm.io/gorm"
	"strconv"
)

type Repository struct {
	dbConnection *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		dbConnection: database.GetDB(),
	}
}

func (r *Repository) First(id int) (*GameModel, error) {
	var result GameModel

	tx := r.dbConnection.
		First(&result, GameModel{ID: id}).
		Preload("League")

	return &result, tx.Error
}

func (r *Repository) FirstOrCreate(game GameModel) (GameModel, error) {
	var result GameModel

	tx := r.dbConnection.
		Attrs(GameModel{
			PlayedMinutes: game.PlayedMinutes,
			OfficialId:    game.OfficialId,
		}).
		FirstOrCreate(&result, GameModel{
			HomeTeamID:  game.HomeTeamID,
			AwayTeamID:  game.AwayTeamID,
			LeagueID:    game.LeagueID,
			ScheduledAt: game.ScheduledAt,
		})

	return result, tx.Error
}

// Exists checks if game exists in db. Can check by id or official_id
func (r *Repository) Exists(game GameModel) (bool, error) {
	var exists bool
	var condition string

	if game.ID != 0 {
		condition = "id = " + strconv.Itoa(game.ID)
	} else {
		if game.OfficialId != "" {
			condition = "official_id = '" + game.OfficialId + "'"
		}
	}

	err := r.dbConnection.
		Model(&GameModel{}).
		Select("count(*) > 0").
		Where(condition).
		Find(&exists).
		Error

	if err != nil {
		return false, err
	}

	return exists, nil
}
