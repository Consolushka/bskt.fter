package players

import (
	"time"
)

type PlayerModel struct {
	Id        uint
	FullName  string
	BirthDate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Id        uint      `gorm:"column:id"`
	FullName  string    `gorm:"column:full_name"`
	BirthDate time.Time `gorm:"column:birth_date_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
}

type GameTeamPlayerStatModel struct {
	Id            uint
	GameTeamId    uint
	PlayerId      uint
	PlayedSeconds int
	PlsMin        int8
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Id            uint      `gorm:"column:id"`
	PlayerId      uint      `gorm:"column:player_id"`
	PlayedSeconds int       `gorm:"column:played_seconds"`
	PlsMin        int8      `gorm:"column:plus_minus"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

}
