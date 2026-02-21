package tournament_poll_logs_repo

import (
	"IMP/app/internal/core/tournament_poll_logs"

	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func (g Gorm) Create(log tournament_poll_logs.TournamentPollLogModel) (tournament_poll_logs.TournamentPollLogModel, error) {
	err := g.db.Create(&log).Error
	return log, err
}

func (g Gorm) GetLatestSuccess(tournamentId uint) (tournament_poll_logs.TournamentPollLogModel, error) {
	var log tournament_poll_logs.TournamentPollLogModel
	err := g.db.Where("tournament_id = ? AND status = ?", tournamentId, tournament_poll_logs.StatusSuccess).
		Order("interval_end DESC").
		First(&log).Error
	return log, err
}

func NewGormRepo(db *gorm.DB) Gorm {
	return Gorm{db: db}
}
