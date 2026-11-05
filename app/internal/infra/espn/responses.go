package espn

type ScoreboardResponse struct {
	Events []EventEntity `json:"events"`
}

type SummaryResponse struct {
	Boxscore BoxscoreEntity `json:"boxscore"`
}

// AthleteResponse — карточка игрока из core API (sports.core.api.espn.com):
// только там дата рождения есть у всех игроков.
type AthleteResponse struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	DateOfBirth string `json:"dateOfBirth"` // "2002-09-30T07:00Z"; бывает пустым
}
