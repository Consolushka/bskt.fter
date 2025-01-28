package schedule_league

type TeamInfoDTO struct {
	TeamId      int    `json:"teamId"`
	TeamName    string `json:"teamName"`
	TeamCity    string `json:"teamCity"`
	TeamTricode string `json:"teamTricode"`
	TeamSlug    string `json:"teamSlug"`
	Wins        int    `json:"wins"`
	Losses      int    `json:"losses"`
	Score       int    `json:"score"`
	Seed        int    `json:"seed"`
}
