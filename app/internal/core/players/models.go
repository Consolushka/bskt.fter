package players

import (
	"time"
)

type PlayerModel struct {
	Id        uint      `gorm:"column:id"`
	FullName  string    `gorm:"column:full_name"`
	BirthDate time.Time `gorm:"column:birth_date_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (PlayerModel) TableName() string {
	return "players"
}

type GameTeamPlayerStatModel struct {
	Id            uint      `gorm:"column:id"`
	GameId        uint      `gorm:"column:game_id"`
	TeamId        uint      `gorm:"column:team_id"`
	PlayerId      uint      `gorm:"column:player_id"`
	PlayedSeconds int       `gorm:"column:played_seconds"`
	PlsMin        int8      `gorm:"column:plus_minus"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (GameTeamPlayerStatModel) TableName() string {
	return "game_team_player_stats"
}
