package ports

import "IMP/app/internal/core/teams"

type TeamsRepo interface {
	FirstOrCreate(model teams.TeamModel) (teams.TeamModel, error)
	FirstOrCreateStats(model teams.GameTeamStatModel) (teams.GameTeamStatModel, error)
}
