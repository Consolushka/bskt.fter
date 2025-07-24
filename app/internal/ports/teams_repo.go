package ports

import "IMP/app/internal/core/teams"

type TeamsRepo interface {
	FirstOrCreateTeam(model teams.TeamModel) (teams.TeamModel, error)
	FirstOrCreateTeamStats(model teams.GameTeamStatModel) (teams.GameTeamStatModel, error)
}
