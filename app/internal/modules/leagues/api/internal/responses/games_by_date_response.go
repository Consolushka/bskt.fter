package responses

import (
	gamesModels "IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/utils/array_utils"
	"time"
)

type GamesByDateResponse struct {
	Date  string         `json:"date"`
	Games []GameResponse `json:"games"`
}

type GameResponse struct {
	Id            int    `json:"id"`
	HomeTeamAlias string `json:"home_team_alias"`
	HomeTeamScore int    `json:"home_team_score"`
	AwayTeamAlias string `json:"away_team_alias"`
	AwayTeamScore int    `json:"away_team_score"`
}

func NewGamesByDateResponse(date time.Time, games []gamesModels.Game) GamesByDateResponse {
	return GamesByDateResponse{
		Date: date.Format("02-01-2006"),
		Games: array_utils.Map(games, func(game gamesModels.Game) GameResponse {
			return GameResponse{
				Id:            game.ID,
				HomeTeamAlias: game.HomeTeamStats.Team.Alias,
				AwayTeamAlias: game.AwayTeamStats.Team.Alias,
				HomeTeamScore: game.HomeTeamStats.Points,
				AwayTeamScore: game.AwayTeamStats.Points,
			}
		})}
}
