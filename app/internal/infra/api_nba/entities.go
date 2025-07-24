package api_nba

import "time"

type GameEntity struct {
	Id          int                   `json:"id"`
	League      string                `json:"league"`
	Season      int                   `json:"season"`
	Date        GameDateEntity        `json:"date"`
	Stage       int                   `json:"stage"`
	Status      GameStatusEntity      `json:"status"`
	Periods     GamePeriodsEntity     `json:"periods"`
	Arena       ArenaEntity           `json:"arena"`
	Teams       GameTeamsEntity       `json:"teams"`
	Scores      GameTeamsScoresEntity `json:"scores"`
	Officials   []string              `json:"officials"`
	TimesTied   int                   `json:"timesTied"`
	LeadChanges int                   `json:"leadChanges"`
	Nugget      interface{}           `json:"nugget"`
}

type GameDateEntity struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration string    `json:"duration"`
}

type GameStatusEntity struct {
	Clock    interface{} `json:"clock"`
	Halftime bool        `json:"halftime"`
	Short    int         `json:"short"`
	Long     string      `json:"long"`
}

type GamePeriodsEntity struct {
	Current     int  `json:"current"`
	Total       int  `json:"total"`
	EndOfPeriod bool `json:"endOfPeriod"`
}

type ArenaEntity struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type GameTeamsEntity struct {
	Visitors TeamEntity `json:"visitors"`
	Home     TeamEntity `json:"home"`
}

type GameTeamsScoresEntity struct {
	Visitors TeamScoresEntity `json:"visitors"`
	Home     TeamScoresEntity `json:"home"`
}

type TeamEntity struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Code     string `json:"code"`
	Logo     string `json:"logo"`
}

type TeamScoresEntity struct {
	Win       int                    `json:"win"`
	Loss      int                    `json:"loss"`
	Series    TeamScoresSeriesEntity `json:"series"`
	Linescore []string               `json:"linescore"`
	Points    int                    `json:"points"`
}

type TeamScoresSeriesEntity struct {
	Win  int `json:"win"`
	Loss int `json:"loss"`
}

type PlayerStatisticsPlayerGeneralDataEntity struct {
	Id        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type PlayerStatisticEntity struct {
	Player PlayerStatisticsPlayerGeneralDataEntity `json:"player"`
	Team   TeamEntity                              `json:"team"`
	Game   struct {
		Id int `json:"id"`
	} `json:"game"`
	Points    int         `json:"points"`
	Pos       *string     `json:"pos"`
	Min       string      `json:"min"`
	Fgm       int         `json:"fgm"`
	Fga       int         `json:"fga"`
	Fgp       string      `json:"fgp"`
	Ftm       int         `json:"ftm"`
	Fta       int         `json:"fta"`
	Ftp       string      `json:"ftp"`
	Tpm       int         `json:"tpm"`
	Tpa       int         `json:"tpa"`
	Tpp       string      `json:"tpp"`
	OffReb    int         `json:"offReb"`
	DefReb    int         `json:"defReb"`
	TotReb    int         `json:"totReb"`
	Assists   int         `json:"assists"`
	PFouls    int         `json:"pFouls"`
	Steals    int         `json:"steals"`
	Turnovers int         `json:"turnovers"`
	Blocks    int         `json:"blocks"`
	PlusMinus string      `json:"plusMinus"`
	Comment   interface{} `json:"comment"`
}

type PlayerEntity struct {
	Id          int                 `json:"id"`
	Firstname   string              `json:"firstname"`
	Lastname    string              `json:"lastname"`
	Birth       PlayerBirthEntity   `json:"birth"`
	Nba         PlayerNbaEntity     `json:"nba"`
	Height      PlayerHeightEntity  `json:"height"`
	Weight      PlayerWeightEntity  `json:"weight"`
	College     string              `json:"college"`
	Affiliation string              `json:"affiliation"`
	Leagues     PlayerLeaguesEntity `json:"leagues"`
}

type PlayerBirthEntity struct {
	Date    string `json:"date"`
	Country string `json:"country"`
}

type PlayerNbaEntity struct {
	Start int `json:"start"`
	Pro   int `json:"pro"`
}

type PlayerHeightEntity struct {
	Feets  string `json:"feets"`
	Inches string `json:"inches"`
	Meters string `json:"meters"`
}

type PlayerWeightEntity struct {
	Pounds    string `json:"pounds"`
	Kilograms string `json:"kilograms"`
}

type PlayerLeaguesEntity struct {
	Standard PlayerLeagueStandardEntity `json:"standard"`
}

type PlayerLeagueStandardEntity struct {
	Jersey int    `json:"jersey"`
	Active bool   `json:"active"`
	Pos    string `json:"pos"`
}
