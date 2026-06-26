package ports

import "IMP/app/internal/core/leagues"

type LeaguesRepo interface {
	FirstOrCreate(model leagues.LeagueModel) (leagues.LeagueModel, error)
}
