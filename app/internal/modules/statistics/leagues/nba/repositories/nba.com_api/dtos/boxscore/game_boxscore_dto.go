package boxscore

import (
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/modules/statistics/enums"
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

// ToImpModel converts a GameBoxScoreDTO to a models.GameModel which is neccessary for IMP package
func (dto GameDTO) ToImpModel() *models.GameModel {
	league := enums.NBA
	duration := 0

	duration = 4 * league.QuarterDuration()
	for i := 0; i < dto.Period-4; i++ {
		duration += league.OvertimeDuration()
	}

	return &models.GameModel{
		Scheduled:    dto.GameTimeUTC.Format("2006-01-02 15:04:05"),
		FullGameTime: duration,
		Home:         dto.HomeTeam.ToImpModel(),
		Away:         dto.AwayTeam.ToImpModel(),
		League:       enums.NBA,
	}
}
