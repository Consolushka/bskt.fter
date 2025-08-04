package tournaments

import (
	"IMP/app/internal/core/leagues"
	"gorm.io/gorm"
	"time"
)

type TournamentModel struct {
	Id        uint           `gorm:"column:id"`
	LeagueId  uint           `gorm:"column:league_id"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	League leagues.LeagueModel `gorm:"foreignKey:LeagueId"`
}

func (TournamentModel) TableName() string {
	return "tournaments"
}
