package todays_games

import "time"

type GameDTO struct {
	GameId            string    `json:"gameId"`
	GameCode          string    `json:"gameCode"`
	GameStatus        int       `json:"gameStatus"`
	GameStatusText    string    `json:"gameStatusText"`
	Period            int       `json:"period"`
	GameClock         string    `json:"gameClock"`
	GameTimeUTC       time.Time `json:"gameTimeUTC"`
	GameET            time.Time `json:"gameEt"`
	RegulationPeriods int       `json:"regulationPeriods"`
	IfNecessary       bool      `json:"ifNecessary"`
	SeriesGameNumber  string    `json:"seriesGameNumber"`
	GameLabel         string    `json:"gameLabel"`
	GameSubLabel      string    `json:"gameSubLabel"`
	SeriesText        string    `json:"seriesText"`
	SeriesConference  string    `json:"seriesConference"`
	PoRoundDesc       string    `json:"poRoundDesc"`
	GameSubtype       string    `json:"gameSubtype"`
	IsNeutral         bool      `json:"isNeutral"`
}
