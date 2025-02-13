package responses

import (
	leaguesDomain "IMP/app/internal/modules/leagues/domain/resources"
	"IMP/app/internal/modules/teams/models"
)

type TeamResponse struct {
	Id     int                          `json:"id"`
	Name   string                       `json:"name"`
	Alias  string                       `json:"alias"`
	League leaguesDomain.LeagueResponse `json:"league"`
}

func NewTeamResponse(teamModel models.Team) TeamResponse {
	return TeamResponse{
		Id:     teamModel.ID,
		Name:   teamModel.Name,
		Alias:  teamModel.Alias,
		League: leaguesDomain.NewLeagueResponse(&teamModel.League),
	}

}
