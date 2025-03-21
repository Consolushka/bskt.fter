package domain

import (
	"IMP/app/database"
	"IMP/app/internal/modules/games/domain/models"
	"database/sql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Repository struct {
	dbConnection *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		dbConnection: database.GetDB(),
	}
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

func (r *Repository) GamesStatsByDateList(date time.Time, leagueId *int) ([]models.Game, error) {
	var gamesModel []models.Game

	builder := r.gameStatsBuilder(nil).
		Where("DATE(scheduled_at) = @date", sql.Named("date", date.Format("2006-01-02")))

	if leagueId != nil {
		builder = builder.Where("league_id = ?", *leagueId)
	}

	tx := builder.
		Debug().
		Find(&gamesModel)

	return gamesModel, tx.Error
}

func (r *Repository) ListOfGamesByGamesIdsAndPlayerId(gamesIds []int, playerId int) ([]models.Game, error) {
	var gamesModel []models.Game

	tx := r.gameStatsBuilder(&playerId).
		Where("id IN (?)", gamesIds).
		Find(&gamesModel)

	return gamesModel, tx.Error
}

func (r *Repository) FirstGameStatsById(id int) (*models.Game, error) {
	var gameModel models.Game

	tx := r.gameStatsBuilder(nil).
		First(&gameModel, models.Game{ID: id})

	return &gameModel, tx.Error
}

func (r *Repository) gameBuilder() *gorm.DB {
	return r.dbConnection.
		Preload("League")
}

func (r *Repository) gameStatsBuilder(playerId *int) *gorm.DB {
	tx := r.gameBuilder().
		Preload("HomeTeamStats").
		Preload("HomeTeamStats.Team").
		Preload("HomeTeamStats.PlayerGameStats.Player").
		Preload("AwayTeamStats").
		Preload("AwayTeamStats.Team").
		Preload("AwayTeamStats.PlayerGameStats.Player")

	if playerId != nil {
		tx = tx.
			Preload("HomeTeamStats.PlayerGameStats", "player_id", *playerId).
			Preload("AwayTeamStats.PlayerGameStats", "player_id", *playerId)
	}
	return tx
}
