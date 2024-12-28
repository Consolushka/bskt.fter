package boxscore

import (
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/modules/statistics/enums"
	"fmt"
	"strings"
)

type GameInfo struct {
	IsOnline          bool           `json:"IsOnline"`
	GameStatus        int            `json:"GameStatus"`
	MaxPeriod         int            `json:"MaxPeriod"`
	FromDate          interface{}    `json:"FromDate"`
	GameDate          string         `json:"GameDate"`
	HasTime           bool           `json:"HasTime"`
	GameTime          string         `json:"GameTime"`
	GameTimeMsk       string         `json:"GameTimeMsk"`
	HasVideo          bool           `json:"HasVideo"`
	GameTeams         []TeamBoxscore `json:"GameTeams"`
	CompNameRu        string         `json:"CompNameRu"`
	CompNameEn        string         `json:"CompNameEn"`
	LeagueNameRu      string         `json:"LeagueNameRu"`
	LeagueNameEn      string         `json:"LeagueNameEn"`
	LeagueShortNameRu string         `json:"LeagueShortNameRu"`
	LeagueShortNameEn string         `json:"LeagueShortNameEn"`
	Gender            int            `json:"Gender"`
	CompID            int            `json:"CompID"`
	LeagueID          int            `json:"LeagueID"`
	Is3x3             bool           `json:"Is3x3"`
}

// ToImpModel converts a GameBoxScoreDTO to a models.GameModel which is neccessary for IMP package
func (dto GameInfo) ToImpModel() *models.GameModel {
	league := enums.MLBL
	duration := 0

	duration = 4 * league.QuarterDuration()
	for i := 0; i < dto.MaxPeriod-4; i++ {
		duration += league.OvertimeDuration()
	}

	// Convert date from dd.mm.yyyy to yyyy-mm-dd
	dateParts := strings.Split(dto.GameDate, ".")
	formattedDate := dto.GameDate
	if len(dateParts) == 3 {
		formattedDate = fmt.Sprintf("%s-%s-%s", dateParts[2], dateParts[1], dateParts[0])
	}

	return &models.GameModel{
		Scheduled:    formattedDate + " " + dto.GameTime,
		FullGameTime: duration,
		Home:         dto.GameTeams[0].ToImpModel(),
		Away:         dto.GameTeams[1].ToImpModel(),
		League:       enums.MLBL,
	}
}
