package cdn_nba

import "time"

type MetaDto struct {
	Version int    `json:"version"`
	Code    int    `json:"code"`
	Request string `json:"request"`
	Time    string `json:"time"`
}

type BoxScoreDto struct {
	GameId            string          `json:"gameId"`
	GameTimeLocal     time.Time       `json:"gameTimeLocal"`
	GameTimeUTC       time.Time       `json:"gameTimeUTC"`
	GameTimeHome      time.Time       `json:"gameTimeHome"`
	GameTimeAway      time.Time       `json:"gameTimeAway"`
	GameET            time.Time       `json:"gameEt"`
	Duration          int             `json:"duration"`
	GameCode          string          `json:"gameCode"`
	GameStatusText    string          `json:"gameStatusText"`
	GameStatus        int             `json:"gameStatus"`
	RegulationPeriods int             `json:"regulationPeriods"`
	Period            int             `json:"period"`
	GameClock         string          `json:"gameClock"`
	Attendance        int             `json:"attendance"`
	Sellout           string          `json:"sellout"`
	HomeTeam          TeamBoxScoreDto `json:"homeTeam"`
	AwayTeam          TeamBoxScoreDto `json:"awayTeam"`
}

type TeamBoxScoreDto struct {
	TeamId            int                 `json:"teamId"`
	TeamName          string              `json:"teamName"`
	TeamCity          string              `json:"teamCity"`
	TeamTricode       string              `json:"teamTricode"`
	Score             int                 `json:"score"`
	InBonus           string              `json:"inBonus"`
	TimeoutsRemaining int                 `json:"timeoutsRemaining"`
	Players           []PlayerBoxScoreDto `json:"players"`
}

type PlayerBoxScoreDto struct {
	Status     string                      `json:"status"`
	Order      int                         `json:"order"`
	PersonId   int                         `json:"personId"`
	JerseyNum  string                      `json:"jerseyNum"`
	Position   string                      `json:"position"`
	Starter    string                      `json:"starter"`
	Oncourt    string                      `json:"oncourt"`
	Played     string                      `json:"played"`
	Statistics PlayerEfficiencyBoxScoreDto `json:"statistics"`
	Name       string                      `json:"name"`
	NameI      string                      `json:"nameI"`
	FirstName  string                      `json:"firstName"`
	FamilyName string                      `json:"familyName"`
}

type PlayerEfficiencyBoxScoreDto struct {
	Assists                 int     `json:"assists"`
	Blocks                  int     `json:"blocks"`
	BlocksReceived          int     `json:"blocksReceived"`
	FieldGoalsAttempted     int     `json:"fieldGoalsAttempted"`
	FieldGoalsMade          int     `json:"fieldGoalsMade"`
	FieldGoalsPercentage    float64 `json:"fieldGoalsPercentage"`
	FoulsOffensive          int     `json:"foulsOffensive"`
	FoulsDrawn              int     `json:"foulsDrawn"`
	FoulsPersonal           int     `json:"foulsPersonal"`
	FoulsTechnical          int     `json:"foulsTechnical"`
	FreeThrowsAttempted     int     `json:"freeThrowsAttempted"`
	FreeThrowsMade          int     `json:"freeThrowsMade"`
	FreeThrowsPercentage    float64 `json:"freeThrowsPercentage"`
	Minus                   float64 `json:"minus"`
	Minutes                 string  `json:"minutes"`
	MinutesCalculated       string  `json:"minutesCalculated"`
	Plus                    float64 `json:"plus"`
	PlusMinusPoints         float64 `json:"plusMinusPoints"`
	Points                  int     `json:"points"`
	PointsFastBreak         int     `json:"pointsFastBreak"`
	PointsInThePaint        int     `json:"pointsInThePaint"`
	PointsSecondChance      int     `json:"pointsSecondChance"`
	ReboundsDefensive       int     `json:"reboundsDefensive"`
	ReboundsOffensive       int     `json:"reboundsOffensive"`
	ReboundsTotal           int     `json:"reboundsTotal"`
	Steals                  int     `json:"steals"`
	ThreePointersAttempted  int     `json:"threePointersAttempted"`
	ThreePointersMade       int     `json:"threePointersMade"`
	ThreePointersPercentage float64 `json:"threePointersPercentage"`
	Turnovers               int     `json:"turnovers"`
	TwoPointersAttempted    int     `json:"twoPointersAttempted"`
	TwoPointersMade         int     `json:"twoPointersMade"`
	TwoPointersPercentage   float64 `json:"twoPointersPercentage"`
}

type SeasonScheduleDto struct {
	SeasonYear string                      `json:"seasonYear"`
	LeagueId   string                      `json:"leagueId"`
	Games      []GameDateSeasonScheduleDto `json:"gameDates"`
	Weeks      []WeekSeasonScheduleDto     `json:"weeks"`
}

type GameDateSeasonScheduleDto struct {
	GameDate string                  `json:"gameDate"`
	Games    []GameSeasonScheduleDto `json:"games"`
}

type WeekSeasonScheduleDto struct {
	WeekNumber int    `json:"weekNumber"`
	WeekName   string `json:"weekName"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}

type GameSeasonScheduleDto struct {
	GameId           string                    `json:"gameId"`
	GameCode         string                    `json:"gameCode"`
	GameStatus       int                       `json:"gameStatus"`
	GameStatusText   string                    `json:"gameStatusText"`
	GameSequence     int                       `json:"gameSequence"`
	GameDateEst      string                    `json:"gameDateEst"`
	GameDateTimeEst  string                    `json:"gameDateTimeEst"`
	GameDateUTC      string                    `json:"gameDateUTC"`
	GameDateTimeUTC  string                    `json:"gameDateTimeUTC"`
	AwayTeamTime     string                    `json:"awayTeamTime"`
	HomeTeamTime     string                    `json:"homeTeamTime"`
	Day              string                    `json:"day"`
	MonthNum         int                       `json:"monthNum"`
	WeekNumber       int                       `json:"weekNumber"`
	WeekName         string                    `json:"weekName"`
	IfNecessary      bool                      `json:"ifNecessary"`
	SeriesGameNumber string                    `json:"seriesGameNumber"`
	GameLabel        string                    `json:"gameLabel"`
	GameSubLabel     string                    `json:"gameSubLabel"`
	SeriesText       string                    `json:"seriesText"`
	ArenaName        string                    `json:"arenaName"`
	ArenaState       string                    `json:"arenaState"`
	ArenaCity        string                    `json:"arenaCity"`
	PostponedStatus  string                    `json:"postponedStatus"`
	BranchLink       string                    `json:"branchLink"`
	GameSubtype      string                    `json:"gameSubtype"`
	IsNeutral        bool                      `json:"isNeutral"`
	HomeTeam         TeamInfoSeasonScheduleDto `json:"homeTeam"`
	AwayTeam         TeamInfoSeasonScheduleDto `json:"awayTeam"`
	PointsLeaders    []interface{}             `json:"pointsLeaders"`
}

type TeamInfoSeasonScheduleDto struct {
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
