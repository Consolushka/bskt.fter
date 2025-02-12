package resources

import (
	"IMP/app/internal/modules/players"
	"IMP/app/internal/modules/players/domain/resources"
	"IMP/app/internal/modules/teams/models"
	"IMP/app/internal/utils/array_utils"
)

type Team struct {
	FullName string             `json:"full_name"`
	Alias    string             `json:"alias"`
	Score    int                `json:"score"`
	Players  []resources.Player `json:"players"`
}

func NewTeamResource(teamModel models.TeamGameStats) Team {
	return Team{
		FullName: teamModel.Team.Name,
		Alias:    teamModel.Team.Alias,
		Score:    teamModel.Points,
		Players: array_utils.Map(teamModel.PlayerGameStats, func(playerGameStats players.PlayerGameStats) resources.Player {
			return resources.NewPlayerResource(playerGameStats)
		}),
	}
}
