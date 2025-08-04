package games

import "time"

type GameModel struct {
	Id           uint      `gorm:"column:id"`
	TournamentId uint      `gorm:"column:tournament_id"`
	ScheduledAt  time.Time `gorm:"column:scheduled_at"`
	Title        string    `gorm:"column:title"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (GameModel) TableName() string {
	return "games"
}
