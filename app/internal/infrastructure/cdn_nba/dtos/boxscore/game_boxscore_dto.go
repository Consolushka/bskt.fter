package boxscore

import (
	"time"
)

type GameDTO struct {
	GameId            string    `json:"gameId"`
	GameTimeLocal     time.Time `json:"gameTimeLocal"`
	GameTimeUTC       time.Time `json:"gameTimeUTC"`
	GameTimeHome      time.Time `json:"gameTimeHome"`
	GameTimeAway      time.Time `json:"gameTimeAway"`
	GameET            time.Time `json:"gameEt"`
	Duration          int       `json:"duration"`
	GameCode          string    `json:"gameCode"`
	GameStatusText    string    `json:"gameStatusText"`
	GameStatus        int       `json:"gameStatus"`
	RegulationPeriods int       `json:"regulationPeriods"`
	Period            int       `json:"period"`
	GameClock         string    `json:"gameClock"`
	Attendance        int       `json:"attendance"`
	Sellout           string    `json:"sellout"`
	HomeTeam          TeamDTO   `json:"homeTeam"`
	AwayTeam          TeamDTO   `json:"awayTeam"`
}
