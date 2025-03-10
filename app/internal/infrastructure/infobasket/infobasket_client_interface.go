package infobasket

type ClientInterface interface {
	BoxScore(gameId string) GameBoxScoreResponse
	ScheduledGames(compId int) []GameScheduleDto
	TeamGames(teamId string) TeamScheduleResponse
}
