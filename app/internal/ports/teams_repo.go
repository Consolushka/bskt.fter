package ports

import "IMP/app/internal/core/teams"

type TeamsRepo interface {
	// created = true, если строки не было и она создана этим вызовом.
	FirstOrCreate(model teams.TeamModel) (team teams.TeamModel, created bool, err error)
	FirstOrCreateStats(model teams.GameTeamStatModel) (teams.GameTeamStatModel, error)
	// found = false, если связи нет; это не ошибка.
	FindByExternalId(scope string, externalId int) (team teams.TeamModel, found bool, err error)
	// Повторная запись той же пары (scope, externalId) ничего не меняет и не возвращает ошибку.
	LinkExternalId(scope string, externalId int, teamId uint) error
}
