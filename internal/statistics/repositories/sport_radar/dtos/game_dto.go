package dtos

import (
	"FTER/internal/enums"
	"FTER/internal/fter/models"
)

type GameDTO struct {
	ID           string       `json:"id"`
	Status       string       `json:"status"`
	Coverage     string       `json:"coverage"`
	Scheduled    string       `json:"scheduled"`
	Duration     string       `json:"duration"`
	Attendance   int          `json:"attendance"`
	LeadChanges  int          `json:"lead_changes"`
	TimesTied    int          `json:"times_tied"`
	Clock        string       `json:"clock"`
	Quarter      int          `json:"quarter"`
	TrackOnCourt bool         `json:"track_on_court"`
	Reference    string       `json:"reference"`
	EntryMode    string       `json:"entry_mode"`
	SrID         string       `json:"sr_id"`
	ClockDecimal string       `json:"clock_decimal"`
	Home         TeamStatsDTO `json:"home"`
	Away         TeamStatsDTO `json:"away"`
}

// ToFterModel converts a GameDTO to a models.GameModel which is neccessary for FTER package
func (dto GameDTO) ToFterModel() *models.GameModel {
	league := enums.NBA
	duration := 0
	if dto.Quarter > 4 {
		duration = 4 * league.QuarterDuration()
		for i := 0; i < dto.Quarter-4; i++ {
			duration += league.OvertimeDuration()
		}
	}

	return &models.GameModel{
		Scheduled:    dto.Scheduled,
		FullGameTime: duration,
		Home:         dto.Home.ToFterModel(),
		Away:         dto.Away.ToFterModel(),
		League:       enums.NBA,
	}
}
