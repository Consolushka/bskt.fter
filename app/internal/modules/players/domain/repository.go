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

func (r *Repository) FirstOrCreatePlayerGameStats(stats models.PlayerGameStats) error {
	tx := r.dbConnection.Attrs(
		models.PlayerGameStats{
			PlayedSeconds: stats.PlayedSeconds,
			PlsMin:        stats.PlsMin,
			IsBench:       stats.IsBench,
			IMPClean:      stats.IMPClean,
		}).
		FirstOrCreate(
			&models.PlayerGameStats{},
			models.PlayerGameStats{
				PlayerID:   stats.PlayerID,
				TeamGameId: stats.TeamGameId,
			})

	return tx.Error
}

func (r *Repository) ListByFullName(fullName string) ([]models.Player, error) {
	var result []models.Player

	tx := r.dbConnection.
		Where("full_name_local LIKE ?", "%"+fullName+"%").
		Or("full_name_en LIKE ?", "%"+fullName+"%").
		Find(&result)

	return result, tx.Error
}

func (r *Repository) ListOfGamesByPlayerId(playerId int) ([]int, error) {
	var gameIds []int

	tx := r.dbConnection.
		Table("team_game_stats").
		Select("team_game_stats.game_id").
		Where("team_game_stats.id IN (?)",
			r.dbConnection.Table("player_game_stats").
				Select("player_game_stats.team_game_id").
				Where("player_game_stats.player_id = ?", playerId)).
		Find(&gameIds)

	return gameIds, tx.Error
}

func (r *Repository) ListOfPlayersGamesStats() ([]models.PlayerGameStats, error) {
	var playerGameStats []models.PlayerGameStats

	tx := r.dbConnection.
		Model(models.PlayerGameStats{}).
		Find(&playerGameStats)

	return playerGameStats, tx.Error
}

func (r *Repository) ListOfPlayersIds() ([]int, error) {
	var playerIds []int

	tx := r.dbConnection.Model(models.Player{}).
		Select("id").
		Find(&playerIds)

	return playerIds, tx.Error
}
