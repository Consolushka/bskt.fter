package sportoteka

import "time"

type GameBoxScoreEntity struct {
	Teams         []TeamBoxScoreEntity `json:"teams"`
	Team1         TeamInfoEntity       `json:"team1"`
	Team2         TeamInfoEntity       `json:"team2"`
	Ot            string               `json:"ot"`
	DistanceIndex int                  `json:"distanceIndex"`
	Game          GameInfoEntity       `json:"game"`
}

type TeamBoxScoreEntity struct {
	TeamNumber    int                       `json:"teamNumber"`
	TeamId        int                       `json:"teamId"`
	UniformNumber int                       `json:"uniformNumber"`
	UniformType   int                       `json:"uniformType"`
	DominantColor string                    `json:"dominantColor"`
	SecondColor   string                    `json:"secondColor"`
	BorderColor   string                    `json:"borderColor"`
	NumberColor   string                    `json:"numberColor"`
	Starts        []TeamBoxScoreStartEntity `json:"starts"`
	Total         TeamBoxScoreTotalsEntity  `json:"total"`
}

type TeamBoxScoreTotalsEntity struct {
	Points    int `json:"points"`
	Goal2     int `json:"goal2"`
	Shot2     int `json:"shot2"`
	Goal3     int `json:"goal3"`
	Shot3     int `json:"shot3"`
	Goal1     int `json:"goal1"`
	Shot1     int `json:"shot1"`
	Assist    int `json:"assist"`
	Pass      int `json:"pass"`
	Steal     int `json:"steal"`
	Block     int `json:"block"`
	Blocked   int `json:"blocked"`
	DefReb    int `json:"defReb"`
	OffReb    int `json:"offReb"`
	FoulsOn   int `json:"foulsOn"`
	Turnover  int `json:"turnover"`
	Foul      int `json:"foul"`
	FoulT     int `json:"foulT"`
	FoulD     int `json:"foulD"`
	FoulC     int `json:"foulC"`
	FoulB     int `json:"foulB"`
	Second    int `json:"second"`
	Dunk      int `json:"dunk"`
	FastBreak int `json:"fastBreak"`
	PlusMinus int `json:"plusMinus"`
}

type PlayerBoxScoreStatsEntity struct {
	StartNum     int    `json:"startNum"`
	PlayerNumber int    `json:"playerNumber"`
	Role         string `json:"role"`
	IsStart      bool   `json:"isStart"`
	IsOnCourt    bool   `json:"isOnCourt"`
	Points       int    `json:"points"`
	Goal2        int    `json:"goal2"`
	Shot2        int    `json:"shot2"`
	Goal3        int    `json:"goal3"`
	Shot3        int    `json:"shot3"`
	Goal1        int    `json:"goal1"`
	Shot1        int    `json:"shot1"`
	Assist       int    `json:"assist"`
	Pass         int    `json:"pass"`
	Steal        int    `json:"steal"`
	Block        int    `json:"block"`
	Blocked      int    `json:"blocked"`
	DefReb       int    `json:"defReb"`
	OffReb       int    `json:"offReb"`
	FoulsOn      int    `json:"foulsOn"`
	Turnover     int    `json:"turnover"`
	Foul         int    `json:"foul"`
	FoulT        int    `json:"foulT"`
	FoulD        int    `json:"foulD"`
	FoulC        int    `json:"foulC"`
	FoulB        int    `json:"foulB"`
	Second       int    `json:"second"`
	Dunk         int    `json:"dunk"`
	FastBreak    int    `json:"fastBreak"`
	PlusMinus    int    `json:"plusMinus"`
}

type TeamBoxScoreStartEntity struct {
	StartNum      int                       `json:"startNum"`
	TeamNumber    int                       `json:"teamNumber"`
	PlayerNumber  int                       `json:"playerNumber"`
	DisplayNumber string                    `json:"displayNumber"`
	StartRole     string                    `json:"startRole"`
	Position      string                    `json:"position"`
	PositionName  string                    `json:"positionName"`
	Height        *int                      `json:"height"`
	Weight        *int                      `json:"weight"`
	CountryId     *string                   `json:"countryId"`
	CountryName   *string                   `json:"countryName"`
	IsCapitan     bool                      `json:"isCapitan"`
	PersonId      *int                      `json:"personId"`
	LastName      string                    `json:"lastName"`
	FirstName     string                    `json:"firstName"`
	SecondName    string                    `json:"secondName"`
	Birthday      string                    `json:"birthday"`
	Age           *int                      `json:"age"`
	Photo         *string                   `json:"photo"`
	Stats         PlayerBoxScoreStatsEntity `json:"stats"`
}

type TeamInfoEntity struct {
	TeamId     int         `json:"teamId"`
	Start      int         `json:"start"`
	Label      string      `json:"label"`
	AbcName    string      `json:"abcName"`
	Name       string      `json:"name"`
	RegionName string      `json:"regionName"`
	ShortName  string      `json:"shortName"`
	Logo       string      `json:"logo"`
	ArenaId    interface{} `json:"arenaId"`
	Id         int         `json:"id"`
}

type GameInfoEntity struct {
	GameStatus          string      `json:"gameStatus"`
	ShowScore           bool        `json:"showScore"`
	Score1              int         `json:"score1"`
	Score2              int         `json:"score2"`
	Score               string      `json:"score"`
	FullScore           string      `json:"fullScore"`
	Periods             int         `json:"periods"`
	Attendance          int         `json:"attendance"`
	Number              string      `json:"number"`
	FinalTime           interface{} `json:"finalTime"`
	StartTime           interface{} `json:"startTime"`
	ScheduledTime       time.Time   `json:"scheduledTime"`
	DefaultZoneDateTime time.Time   `json:"defaultZoneDateTime"`
	LocalDate           string      `json:"localDate"`
	LocalTime           string      `json:"localTime"`
	DefaultZoneTime     string      `json:"defaultZoneTime"`
	HasTime             bool        `json:"hasTime"`
	RegionId            int         `json:"regionId"`
	ArenaId             int         `json:"arenaId"`
	CompTeam1Id         int         `json:"compTeam1Id"`
	CompTeam2Id         int         `json:"compTeam2Id"`
	Tv                  interface{} `json:"tv"`
	Video               interface{} `json:"video"`
	Id                  int         `json:"id"`
}

type CalendarGameEntity struct {
	Game          GameInfoEntity `json:"game"`
	Gender        int            `json:"gender"`
	Team1         TeamInfoEntity `json:"team1"`
	Team2         TeamInfoEntity `json:"team2"`
	Ot            string         `json:"ot"`
	DistanceIndex int            `json:"distanceIndex"`
}
