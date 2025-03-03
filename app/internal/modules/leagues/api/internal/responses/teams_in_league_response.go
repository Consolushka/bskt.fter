package responses

import "IMP/app/internal/modules/leagues/domain/resources"

type TeamsInLeagueResponse struct {
	League resources.LeagueResource `json:"league"`
	Teams  []TeamInLeagueResponse   `json:"teams"`
}

type TeamInLeagueResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}
