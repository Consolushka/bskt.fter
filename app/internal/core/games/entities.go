package games

import "IMP/app/internal/core/teams"

type GameStatEntity struct {
	GameModel          GameModel
	ExternalGameId     string
	HomeTeamExternalId int
	AwayTeamExternalId int
	HomeTeamStat       teams.TeamStatEntity
	AwayTeamStat       teams.TeamStatEntity
}

// AssignTeams подставляет в игру команды из базы и собирает из их алиасов название.
// Название входит в ключ игры, поэтому собирается только из алиасов, лежащих в базе:
// сокращения провайдера меняются, и одна игра записалась бы дважды.
func (e *GameStatEntity) AssignTeams(home, away teams.TeamModel) {
	e.HomeTeamStat.TeamModel = home
	e.AwayTeamStat.TeamModel = away
	e.GameModel.Title = home.Alias + " - " + away.Alias
}
