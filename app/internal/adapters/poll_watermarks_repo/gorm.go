package poll_watermarks_repo

import (
	"IMP/app/internal/core/poll_watermarks"
	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func (g Gorm) Create(log poll_watermarks.TournamentPollLogModel) (poll_watermarks.TournamentPollLogModel, error) {
	err := g.db.Create(&log).Error
	return log, err
}

func (g Gorm) GetLatestSuccess(tournamentId uint) (poll_watermarks.TournamentPollLogModel, error) {
	var log poll_watermarks.TournamentPollLogModel
	err := g.db.Where("tournament_id = ? AND status = ?", tournamentId, poll_watermarks.StatusSuccess).
		Order("interval_end DESC").
		First(&log).Error
	return log, err
}

func NewGormRepo(db *gorm.DB) Gorm {
	return Gorm{db: db}
}
