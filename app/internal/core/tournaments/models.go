package tournaments

import (
	"IMP/app/internal/core/leagues"
	"gorm.io/gorm"
	"time"
)

type TournamentModel struct {
	Id        uint           `db:"id"`
	LeagueId  uint           `db:"league_id"`
	Name      string         `db:"name"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt gorm.DeletedAt `db:"deleted_at"`

	League leagues.LeagueModel `gorm:"foreignKey:LeagueId"`
}

func (TournamentModel) TableName() string {
	return "tournaments"
}
