package api_nba

import (
	"IMP/app/pkg/http"
	"strconv"
)

type ClientInterface interface {
	Games(id int, date string, leagueId string, season string, teamId string, timezone string) (GamesResponse, error)
	PlayersStatistics(playerId int, gameId int, teamId int, season string) (PlayerStatisticResponse, error)
	PlayerInfo(playerId int, name string, teamId int, season int, country string, search string) (PlayersResponse, error)
}

type Client struct {
	baseUrl     string
	token       string
	baseHeaders map[string]string
}

func NewClient(baseUrl string, token string) ClientInterface {
	return &Client{
		baseUrl: baseUrl,
		token:   token,
		baseHeaders: map[string]string{
			"x-rapidapi-host": "v2.nba.api-sports.io",
			"x-rapidapi-key":  token,
		},
	}
}

func (c Client) Games(id int, date string, leagueId string, season string, teamId string, timezone string) (GamesResponse, error) {
	return http.Get[GamesResponse](c.baseUrl+"/games?date="+date, &c.baseHeaders)
}

func (c Client) PlayersStatistics(playerId int, gameId int, teamId int, season string) (PlayerStatisticResponse, error) {
	//todo: ignore zero values
	return http.Get[PlayerStatisticResponse](c.baseUrl+"/players/statistics?game="+strconv.Itoa(gameId), &c.baseHeaders)
}

func (c Client) PlayerInfo(playerId int, name string, teamId int, season int, country string, search string) (PlayersResponse, error) {
	//todo: ignore zero values
	return http.Get[PlayersResponse](c.baseUrl+"/players?id="+strconv.Itoa(playerId), &c.baseHeaders)
}
