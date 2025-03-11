package infobasket

import "math/rand"

type GameBoxScoreResponse struct {
	IsOnline          bool              `json:"IsOnline"`
	GameStatus        int               `json:"GameStatus"`
	MaxPeriod         int               `json:"MaxPeriod"`
	FromDate          interface{}       `json:"FromDate"`
	GameDate          string            `json:"GameDate"`
	HasTime           bool              `json:"HasTime"`
	GameTime          string            `json:"GameTime"`
	GameTimeMsk       string            `json:"GameTimeMsk"`
	HasVideo          bool              `json:"HasVideo"`
	GameTeams         []TeamBoxScoreDto `json:"GameTeams"`
	CompNameRu        string            `json:"CompNameRu"`
	CompNameEn        string            `json:"CompNameEn"`
	LeagueNameRu      string            `json:"LeagueNameRu"`
	LeagueNameEn      string            `json:"LeagueNameEn"`
	LeagueShortNameRu string            `json:"LeagueShortNameRu"`
	LeagueShortNameEn string            `json:"LeagueShortNameEn"`
	Gender            int               `json:"Gender"`
	CompID            int               `json:"CompID"`
	LeagueID          int               `json:"LeagueID"`
	Is3x3             bool              `json:"Is3x3"`
}

type TeamScheduleResponse struct {
	Games []GameScheduleDto `json:"Games"`
}

type SeasonScheduleResponse struct {
	Games []GameScheduleDto `json:"Games"`
}

func CreateMockGameBoxScoreResponse(homeTeamId, awayTeamId, gameStatus, maxPeriods int, gameDate, gameTime string) GameBoxScoreResponse {
	if homeTeamId == 0 {
		homeTeamId = rand.Intn(1000)
	}
	if awayTeamId == 0 {
		awayTeamId = rand.Intn(1000)
	}

	if awayTeamId == homeTeamId {
		awayTeamId = homeTeamId - 1
	}

	return GameBoxScoreResponse{
		IsOnline:    false,
		GameStatus:  gameStatus, // 2 typically represents a completed game
		MaxPeriod:   maxPeriods, // Standard basketball game has 4 periods
		FromDate:    nil,
		GameDate:    gameDate,
		HasTime:     true,
		GameTime:    gameTime,
		GameTimeMsk: "21.30",
		HasVideo:    true,
		GameTeams: []TeamBoxScoreDto{
			CreateMockTeamBoxScoreDto("PHX", "Pheonix", 100, homeTeamId, 12),
			CreateMockTeamBoxScoreDto("LAL", "Lakers", 101, awayTeamId, 12),
		},
		CompNameRu:        "Регулярный чемпионат",
		CompNameEn:        "Regular Championship",
		LeagueNameRu:      "Тамбовская Баскетбольная Лига",
		LeagueNameEn:      "Tambov Basketball League",
		LeagueShortNameRu: "ТБЛ",
		LeagueShortNameEn: "TBL",
		Gender:            1, // 1 for men, 2 for women
		CompID:            89960,
		LeagueID:          123,
		Is3x3:             false,
	}
}
