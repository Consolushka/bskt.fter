package teams

import (
	"strings"
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
	var alias string
	if len(runes) > 3 {
		alias = string(runes[:3])
	} else {
		alias = t.Name
	}
	return strings.ToUpper(alias)
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
