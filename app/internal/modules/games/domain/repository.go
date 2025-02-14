package domain

import (
	"IMP/app/database"
	"IMP/app/internal/modules/games/domain/models"
	"gorm.io/gorm"
	"strconv"
)

// todo: refactor all repositories
type Repository struct {
	dbConnection *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		dbConnection: database.GetDB(),
	}
}

func (r *Repository) First(id int) (*models.Game, error) {
	var result models.Game

	tx := r.dbConnection.
		First(&result, models.Game{ID: id}).
		Preload("League")

	return &result, tx.Error
}

func (r *Repository) FirstOrCreate(game models.Game) (models.Game, error) {
	var result models.Game

	tx := r.dbConnection.
		Attrs(models.Game{
			PlayedMinutes: game.PlayedMinutes,
			OfficialId:    game.OfficialId,
		}).
		FirstOrCreate(&result, models.Game{
			HomeTeamID:  game.HomeTeamID,
			AwayTeamID:  game.AwayTeamID,
			LeagueID:    game.LeagueID,
			ScheduledAt: game.ScheduledAt,
		})

	return result, tx.Error
}

// Exists checks if game exists in db. Can check by id or official_id
func (r *Repository) Exists(game models.Game) (bool, error) {
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
		Model(&models.Game{}).
		Select("count(*) > 0").
		Where(condition).
		Find(&exists).
		Error

	if err != nil {
		return false, err
	}

	return exists, nil
}
