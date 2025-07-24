package api_nba

import (
	"IMP/app/pkg/http"
	"strconv"
	"time"
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
	return GamesResponse{
		Response: []GameEntity{
			{
				Id: 4516,
				Date: GameDateEntity{
					Start:    time.Date(2025, 7, 23, 22, 0, 0, 0, time.UTC),
					End:      time.Date(2025, 7, 23, 23, 50, 0, 0, time.UTC),
					Duration: "",
				},
				Teams: GameTeamsEntity{
					Visitors: TeamEntity{
						Id:       23,
						Name:     "New Orleans Pelicans",
						Nickname: "Pelicans",
						Code:     "NOP",
					},
					Home: TeamEntity{
						Id:       31,
						Name:     "San Antonio Spurs",
						Nickname: "Spurs",
						Code:     "SAS",
					},
				},
				Scores: GameTeamsScoresEntity{
					Visitors: TeamScoresEntity{
						Points: 115,
					},
					Home: TeamScoresEntity{
						Points: 120,
					},
				},
			},
		},
	}, nil

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
