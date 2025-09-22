package tournaments

import (
	"IMP/app/internal/core/leagues"
	"time"

	"gorm.io/gorm"
)

type TournamentModel struct {
	Id        uint           `gorm:"column:id"`
	LeagueId  uint           `gorm:"column:league_id"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	League      leagues.LeagueModel         `gorm:"foreignKey:LeagueId"`
	ExternalIds []TournamentExternalIdModel `gorm:"foreignKey:TournamentId"`
}

func (TournamentModel) TableName() string {
	return "tournaments"
}

type TournamentExternalIdModel struct {
	TournamentId uint      `gorm:"column:tournament_id"`
	ProviderName string    `gorm:"column:provider_name"`
	ExternalId   string    `gorm:"column:external_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`

	Tournament TournamentModel `gorm:"foreignKey:TournamentId"`
}

func (TournamentExternalIdModel) TableName() string {
	return "tournament_external_ids"
}
