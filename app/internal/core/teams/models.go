package teams

import "time"

type TeamModel struct {
	Id        uint
	Name      string
	HomeTown  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Id        uint      `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	HomeTown  string    `gorm:"column:home_town"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
}

type GameTeamStatModel struct {
	Id        uint
	GameId    uint
	TeamId    uint
	Score     int
	CreatedAt time.Time
	UpdatedAt time.Time
	Id        uint      `gorm:"column:id"`
	GameId    uint      `gorm:"column:game_id"`
	TeamId    uint      `gorm:"column:team_id"`
	Score     int       `gorm:"column:score"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
}
