package teams

import (
	"time"
)

type TeamModel struct {
	Id        uint      `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	HomeTown  string    `gorm:"column:home_town"`
	Alias     string    `gorm:"column:alias"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (TeamModel) TableName() string {
	return "teams"
}

func (t *TeamModel) AutoGenerateAlias() string {
	if t.Name == "" {
		return ""
	}
	runes := []rune(t.Name)
	if len(runes) > 3 {
		return string(runes[:3])
	}
	return t.Name
}

type GameTeamStatModel struct {
	Id        uint      `gorm:"column:id"`
	GameId    uint      `gorm:"column:game_id"`
	TeamId    uint      `gorm:"column:team_id"`
	Score     int       `gorm:"column:score"`
	FinalDiff int       `gorm:"column:final_differential"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (GameTeamStatModel) TableName() string {
	return "game_team_stats"
}
