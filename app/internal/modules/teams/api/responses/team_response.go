package responses

import (
	leaguesResources "IMP/app/internal/modules/leagues/domain/resources"
	"IMP/app/internal/modules/teams/domain/models"
)

type TeamResponse struct {
	Id     int                             `json:"id"`
	Name   string                          `json:"name"`
	Alias  string                          `json:"alias"`
	League leaguesResources.LeagueResource `json:"league"`
}

func NewTeamResponse(teamModel models.Team) TeamResponse {
	return TeamResponse{
		Id:     teamModel.ID,
		Name:   teamModel.Name,
		Alias:  teamModel.Alias,
		League: leaguesResources.NewLeagueResponse(&teamModel.League),
	}

}
