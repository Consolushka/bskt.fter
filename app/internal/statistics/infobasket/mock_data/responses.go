package mock_data

import (
	infobasket2 "IMP/app/internal/statistics/infobasket"
)

func CreateMockGameBoxScoreResponse(homeTeamId, awayTeamId, gameStatus, maxPeriods int, gameDate, gameTime string) infobasket2.GameBoxScoreResponse {
	return infobasket2.GameBoxScoreResponse{
		IsOnline:    false,
		GameStatus:  gameStatus,
		MaxPeriod:   maxPeriods,
		FromDate:    nil,
		GameDate:    gameDate,
		HasTime:     true,
		GameTime:    gameTime,
		GameTimeMsk: "21.30",
		HasVideo:    true,
		GameTeams: []infobasket2.TeamBoxScoreDto{
			CreateMockTeamBoxScoreDto("PHX", "Pheonix", 100, homeTeamId, 12),
			CreateMockTeamBoxScoreDto("LAL", "Lakers", 101, awayTeamId, 12),
		},
		CompNameRu:        "Регулярный чемпионат",
		CompNameEn:        "Regular Championship",
		LeagueNameRu:      "Тамбовская Баскетбольная Лига",
		LeagueNameEn:      "Tambov Basketball League",
		LeagueShortNameRu: "ТБЛ",
		LeagueShortNameEn: "TBL",
		Gender:            1,
		CompID:            89960,
		LeagueID:          123,
		Is3x3:             false,
	}
}
