package leagues

import "time"

type LeagueModel struct {
	Id        uint      `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	Alias     string    `gorm:"column:alias"`
	Tier      *int16    `gorm:"column:tier"`
	Country   *string   `gorm:"column:country"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (LeagueModel) TableName() string {
	return "leagues"
}
