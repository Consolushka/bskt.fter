package schedule_league

type GameDTO struct {
	GameId           string        `json:"gameId"`
	GameCode         string        `json:"gameCode"`
	GameStatus       int           `json:"gameStatus"`
	GameStatusText   string        `json:"gameStatusText"`
	GameSequence     int           `json:"gameSequence"`
	GameDateEst      string        `json:"gameDateEst"`
	GameDateTimeEst  string        `json:"gameDateTimeEst"`
	GameDateUTC      string        `json:"gameDateUTC"`
	GameDateTimeUTC  string        `json:"gameDateTimeUTC"`
	AwayTeamTime     string        `json:"awayTeamTime"`
	HomeTeamTime     string        `json:"homeTeamTime"`
	Day              string        `json:"day"`
	MonthNum         int           `json:"monthNum"`
	WeekNumber       int           `json:"weekNumber"`
	WeekName         string        `json:"weekName"`
	IfNecessary      bool          `json:"ifNecessary"`
	SeriesGameNumber string        `json:"seriesGameNumber"`
	GameLabel        string        `json:"gameLabel"`
	GameSubLabel     string        `json:"gameSubLabel"`
	SeriesText       string        `json:"seriesText"`
	ArenaName        string        `json:"arenaName"`
	ArenaState       string        `json:"arenaState"`
	ArenaCity        string        `json:"arenaCity"`
	PostponedStatus  string        `json:"postponedStatus"`
	BranchLink       string        `json:"branchLink"`
	GameSubtype      string        `json:"gameSubtype"`
	IsNeutral        bool          `json:"isNeutral"`
	HomeTeam         TeamInfoDTO   `json:"homeTeam"`
	AwayTeam         TeamInfoDTO   `json:"awayTeam"`
	PointsLeaders    []interface{} `json:"pointsLeaders"`
}
