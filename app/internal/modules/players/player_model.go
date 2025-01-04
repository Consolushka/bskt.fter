package players

import "time"

type Player struct {
	ID        int        `json:"id" db:"id"`
	FullName  string     `json:"full_name" db:"full_name"`
	DraftYear *int       `json:"draft_year" db:"draft_year"`
	BirthDate *time.Time `json:"birth_date" db:"birth_date"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}
