package domain

import (
	"IMP/app/database"
	"IMP/app/internal/modules/leagues/domain/models"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		db: database.GetDB(),
	}
}

func (r *Repository) List() ([]models.League, error) {
	var result []models.League

	tx := r.db.Find(&result)

	return result, tx.Error
}

func (r *Repository) FirstByAliasEn(aliasEn string) (*models.League, error) {
	var result models.League
	tx := r.db.First(&result, models.League{AliasEn: strings.ToUpper(aliasEn)})
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &result, nil
}

func (r *Repository) FirstById(id int) (*models.League, error) {
	var leagueModel models.League

	tx := r.db.First(&leagueModel, models.League{ID: id})
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &leagueModel, nil
}

func (r *Repository) ListPlayerRankingByImpClean(leagueId int, limit *int, minMinuterPerGame int, avgMinutes int, gamesPlayed int) (*[]models.PlayerImpRanking, error) {
	var rankings []models.PlayerImpRanking

	query := `
		SELECT 
			ROW_NUMBER() OVER (ORDER BY avg_imp_clean DESC) as rank, * 
		FROM 
			(SELECT  
				 p.id, p.full_name_local, AVG(pgs.imp_clean) as avg_imp_clean, AVG(pgs.played_seconds) as avg_played_seconds, COUNT(*) as games_played 
			 FROM player_game_stats pgs 
				 JOIN players p ON p.id = pgs.player_id 
				 JOIN team_game_stats tgs ON tgs.id = pgs.team_game_id 
				 JOIN games g ON g.id = tgs.game_id 
			 WHERE 
				 g.league_id = ? 
			   AND 
				 pgs.played_seconds > ? 
			 GROUP BY p.id, p.full_name_local, p.full_name_en 
			 ORDER BY avg_imp_clean DESC
			 ) AS players_ranking 
		WHERE 
			games_played > ? 
		  AND 
			avg_played_seconds > ?
		`

	if limit != nil {
		query += " LIMIT " + strconv.Itoa(*limit)
	}

	err := r.db.Raw(query, leagueId, minMinuterPerGame*60, gamesPlayed, avgMinutes*60).Scan(&rankings).Error
	if err != nil {
		return nil, err
	}

	return &rankings, nil
}
