package api_basketball

import "time"

type GameEntity struct {
	Id        int       `json:"id"`
	Date      time.Time `json:"date"`
	Time      string    `json:"time"`
	Timestamp int64     `json:"timestamp"`
	Timezone  string    `json:"timezone"`
	Stage     *string   `json:"stage"`
	Week      string    `json:"week"`
	Venue     string    `json:"venue"`
	Status    Status    `json:"status"`
	League    League    `json:"league"`
	Country   Country   `json:"country"`
	Teams     Teams     `json:"teams"`
	Scores    Scores    `json:"scores"`
}

type Status struct {
	Long  string  `json:"long"`
	Short string  `json:"short"`
	Timer *string `json:"timer"`
}

type League struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Season any    `json:"season"`
	Logo   string `json:"logo"`
}

type Country struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Code *string `json:"code"`
	Flag *string `json:"flag"`
}

type Teams struct {
	Home Team `json:"home"`
	Away Team `json:"away"`
}

type Team struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type Scores struct {
	Home ScoreDetails `json:"home"`
	Away ScoreDetails `json:"away"`
}

type ScoreDetails struct {
	Quarter1 int  `json:"quarter_1"`
	Quarter2 int  `json:"quarter_2"`
	Quarter3 int  `json:"quarter_3"`
	Quarter4 int  `json:"quarter_4"`
	OverTime *int `json:"over_time"`
	Total    int  `json:"total"`
}

type TeamStatsEntity struct {
	Game            GameRef         `json:"game"`
	Team            TeamRef         `json:"team"`
	FieldGoals      StatsDetails    `json:"field_goals"`
	ThreepointGoals StatsDetails    `json:"threepoint_goals"`
	FreethrowsGoals StatsDetails    `json:"freethrows_goals"`
	Rebounds        ReboundsDetails `json:"rebounds"`
	Assists         int             `json:"assists"`
	Steals          int             `json:"steals"`
	Blocks          int             `json:"blocks"`
	Turnovers       int             `json:"turnovers"`
	PersonalFouls   *int            `json:"personal_fouls"`
}

type GameRef struct {
	Id int `json:"id"`
}

type TeamRef struct {
	Id int `json:"id"`
}

type StatsDetails struct {
	Total      *int     `json:"total"`
	Attempts   *int     `json:"attempts"`
	Percentage *float64 `json:"percentage"`
}

type ReboundsDetails struct {
	Total   int  `json:"total"`
	Offence *int `json:"offence"`
	Defense *int `json:"defense"`
}

type PlayerStatsEntity struct {
	Game            GameRef        `json:"game"`
	Team            TeamRef        `json:"team"`
	Player          PlayerRef      `json:"player"`
	Type            string         `json:"type"`
	Minutes         string         `json:"minutes"`
	FieldGoals      StatsDetails   `json:"field_goals"`
	ThreepointGoals StatsDetails   `json:"threepoint_goals"`
	FreethrowsGoals StatsDetails   `json:"freethrows_goals"`
	Rebounds        PlayerRebounds `json:"rebounds"`
	Assists         int            `json:"assists"`
	Points          int            `json:"points"`
}

type PlayerRef struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PlayerRebounds struct {
	Total int `json:"total"`
}

type PlayerInfoEntity struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Number   string     `json:"number"`
	Country  string     `json:"country"`
	Position string     `json:"position"`
	Age      int        `json:"age"`
	Birth    *BirthInfo `json:"birth"` // Based on API NBA, but let's check if API Basketball has it
}

type BirthInfo struct {
	Date string `json:"date"`
}
