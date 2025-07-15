package leagues

import "time"

type LeagueModel struct {
	Id        uint      `db:"id"`
	Name      string    `db:"name"`
	Alias     string    `db:"alias"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (LeagueModel) TableName() string {
	return "leagues"
}
