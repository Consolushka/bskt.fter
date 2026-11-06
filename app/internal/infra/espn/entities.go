package espn

type EventEntity struct {
	Id           string              `json:"id"`
	Date         string              `json:"date"` // "2026-10-29T00:00Z"
	Season       SeasonEntity        `json:"season"`
	Competitions []CompetitionEntity `json:"competitions"`
}

// SeasonEntity.Type: 1 preseason, 2 regular, 3 postseason, 5 play-in.
type SeasonEntity struct {
	Type int `json:"type"`
}

type CompetitionEntity struct {
	Type        CompetitionTypeEntity `json:"type"` // Abbreviation "STD" | "ALLSTAR" | ...
	Status      StatusEntity          `json:"status"`
	Competitors []CompetitorEntity    `json:"competitors"`
}

type CompetitionTypeEntity struct {
	Abbreviation string `json:"abbreviation"`
}

type StatusEntity struct {
	Type   StatusTypeEntity `json:"type"`
	Period int              `json:"period"`
}

type StatusTypeEntity struct {
	Completed bool `json:"completed"`
}

type CompetitorEntity struct {
	HomeAway string     `json:"homeAway"`
	Score    string     `json:"score"`
	Team     TeamEntity `json:"team"`
}

type TeamEntity struct {
	Id           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
	Location     string `json:"location"`
	Name         string `json:"name"`
}

type BoxscoreEntity struct {
	Players []BoxscorePlayersEntity `json:"players"`
}

type BoxscorePlayersEntity struct {
	Team       TeamEntity                `json:"team"`
	Statistics []PlayersStatisticsEntity `json:"statistics"`
}

type PlayersStatisticsEntity struct {
	Labels   []string             `json:"labels"`
	Athletes []AthleteStatsEntity `json:"athletes"`
}

type AthleteStatsEntity struct {
	Athlete    AthleteEntity `json:"athlete"`
	DidNotPlay bool          `json:"didNotPlay"`
	Stats      []string      `json:"stats"`
}

type AthleteEntity struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}
