package dtos

type GameBoxScoreDTO struct {
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
