package nba

type NbaRepository struct {
}

func (n NbaRepository) GetPlayerStatsByGame(playerId string, gameId string) string {
	return "stats"
}
