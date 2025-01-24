package boxscore

type TeamDTO struct {
	TeamId            int         `json:"teamId"`
	TeamName          string      `json:"teamName"`
	TeamCity          string      `json:"teamCity"`
	TeamTricode       string      `json:"teamTricode"`
	Score             int         `json:"score"`
	InBonus           string      `json:"inBonus"`
	TimeoutsRemaining int         `json:"timeoutsRemaining"`
	Players           []PlayerDTO `json:"players"`
}
