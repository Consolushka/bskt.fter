package sport_radar

type GameBoxScoreDTO struct {
	ID           string          `json:"id"`
	Status       string          `json:"status"`
	Coverage     string          `json:"coverage"`
	Scheduled    string          `json:"scheduled"`
	Duration     string          `json:"duration"`
	Attendance   int             `json:"attendance"`
	LeadChanges  int             `json:"lead_changes"`
	TimesTied    int             `json:"times_tied"`
	Clock        string          `json:"clock"`
	Quarter      int             `json:"quarter"`
	TrackOnCourt bool            `json:"track_on_court"`
	Reference    string          `json:"reference"`
	EntryMode    string          `json:"entry_mode"`
	SrID         string          `json:"sr_id"`
	ClockDecimal string          `json:"clock_decimal"`
	Home         TeamBoxScoreDto `json:"home"`
	Away         TeamBoxScoreDto `json:"away"`
}

type TeamBoxScoreDto struct {
	Alias   string              `json:"alias"`
	Points  int                 `json:"points"`
	Players []PlayerBoxScoreDto `json:"players"`
}

type PlayerBoxScoreDto struct {
	FullName        string                 `json:"full_name"`
	JerseyNumber    string                 `json:"jersey_number"`
	ID              string                 `json:"id"`
	FirstName       string                 `json:"first_name"`
	LastName        string                 `json:"last_name"`
	Position        string                 `json:"position"`
	PrimaryPosition string                 `json:"primary_position"`
	Played          bool                   `json:"played"`
	Active          bool                   `json:"active"`
	Starter         bool                   `json:"starter"`
	OnCourt         bool                   `json:"on_court"`
	SrID            string                 `json:"sr_id"`
	Reference       string                 `json:"reference"`
	Statistics      PlayerStatsBoxScoreDto `json:"statistics"`
}

type PlayerStatsBoxScoreDto struct {
	Minutes              string  `json:"minutes"`
	FieldGoalsMade       int     `json:"field_goals_made"`
	FieldGoalsAtt        int     `json:"field_goals_att"`
	FieldGoalsPct        float64 `json:"field_goals_pct"`
	ThreePointsMade      int     `json:"three_points_made"`
	ThreePointsAtt       int     `json:"three_points_att"`
	ThreePointsPct       float64 `json:"three_points_pct"`
	TwoPointsMade        int     `json:"two_points_made"`
	TwoPointsAtt         int     `json:"two_points_att"`
	TwoPointsPct         float64 `json:"two_points_pct"`
	BlockedAtt           int     `json:"blocked_att"`
	FreeThrowsMade       int     `json:"free_throws_made"`
	FreeThrowsAtt        int     `json:"free_throws_att"`
	FreeThrowsPct        float64 `json:"free_throws_pct"`
	OffensiveRebounds    int     `json:"offensive_rebounds"`
	DefensiveRebounds    int     `json:"defensive_rebounds"`
	Rebounds             int     `json:"rebounds"`
	Assists              int     `json:"assists"`
	Turnovers            int     `json:"turnovers"`
	Steals               int     `json:"steals"`
	Blocks               int     `json:"blocks"`
	AssistsTurnoverRatio float64 `json:"assists_turnover_ratio"`
	PersonalFouls        int     `json:"personal_fouls"`
	TechFouls            int     `json:"tech_fouls"`
	FlagrantFouls        int     `json:"flagrant_fouls"`
	PlsMin               int     `json:"pls_min"`
	Points               int     `json:"points"`
	DoubleDouble         bool    `json:"double_double"`
	TripleDouble         bool    `json:"triple_double"`
	EffectiveFgPct       float64 `json:"effective_fg_pct"`
	Efficiency           float64 `json:"efficiency"`
	EfficiencyGameScore  float64 `json:"efficiency_game_score"`
	FoulsDrawn           int     `json:"fouls_drawn"`
	OffensiveFouls       int     `json:"offensive_fouls"`
	PointsInPaint        int     `json:"points_in_paint"`
	PointsInPaintAtt     int     `json:"points_in_paint_att"`
	PointsInPaintMade    int     `json:"points_in_paint_made"`
	PointsInPaintPct     float64 `json:"points_in_paint_pct"`
	PointsOffTurnovers   int     `json:"points_off_turnovers"`
	TrueShootingAtt      float64 `json:"true_shooting_att"`
	TrueShootingPct      float64 `json:"true_shooting_pct"`
	CoachEjections       int     `json:"coach_ejections"`
	CoachTechFouls       int     `json:"coach_tech_fouls"`
	SecondChancePts      int     `json:"second_chance_pts"`
	SecondChancePct      float64 `json:"second_chance_pct"`
	FastBreakPts         int     `json:"fast_break_pts"`
	FastBreakAtt         int     `json:"fast_break_att"`
	FastBreakMade        int     `json:"fast_break_made"`
	FastBreakPct         float64 `json:"fast_break_pct"`
	DefensiveRating      float64 `json:"defensive_rating"`
	OffensiveRating      float64 `json:"offensive_rating"`
	Minus                int     `json:"minus"`
	Plus                 int     `json:"plus"`
	DefensiveReboundsPct float64 `json:"defensive_rebounds_pct"`
	OffensiveReboundsPct float64 `json:"offensive_rebounds_pct"`
	ReboundsPct          float64 `json:"rebounds_pct"`
	StealsPct            float64 `json:"steals_pct"`
	TurnoversPct         float64 `json:"turnovers_pct"`
	SecondChanceAtt      int     `json:"second_chance_att"`
	SecondChanceMade     int     `json:"second_chance_made"`
}
