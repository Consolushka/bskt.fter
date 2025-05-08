package infobasket

// MockClient is a mock implementation of the infobasket client
type MockClient struct {
	BoxScoreFunc       func(gameId string) GameBoxScoreResponse
	ScheduledGamesFunc func(compId int) []GameScheduleDto
	TeamGamesFunc      func(teamId string) TeamScheduleResponse
}

func (m *MockClient) BoxScore(gameId string) GameBoxScoreResponse {
	return m.BoxScoreFunc(gameId)
}

func (m *MockClient) ScheduledGames(compId int) []GameScheduleDto {
	return m.ScheduledGamesFunc(compId)
}

func (m *MockClient) TeamGames(teamId string) TeamScheduleResponse {
	return m.TeamGamesFunc(teamId)
}
