package api_nba

import (
	"IMP/app/pkg/http"
	"context"
	"os"
	"strconv"
	"time"

	"golang.org/x/time/rate"
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
	limiter     *rate.Limiter
}

func NewClient(baseUrl string, token string) ClientInterface {
	rateLimitStr := os.Getenv("API_NBA_RATE_LIMIT_PER_MINUTE")
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil || rateLimit <= 0 {
		rateLimit = 10
	}

	return &Client{
		baseUrl: baseUrl,
		token:   token,
		baseHeaders: map[string]string{
			"x-rapidapi-host": "v2.nba.api-sports.io",
			"x-rapidapi-key":  token,
		},
		limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1),
	}
}

func (c Client) Games(id int, date string, leagueId string, season string, teamId string, timezone string) (GamesResponse, error) {
	_ = c.limiter.Wait(context.Background())
	return http.Get[GamesResponse](c.baseUrl+"/games?date="+date, &c.baseHeaders)
}

func (c Client) PlayersStatistics(playerId int, gameId int, teamId int, season string) (PlayerStatisticResponse, error) {
	//todo: ignore zero values
	//
	// For free plan limit is 10 requests/minute
	_ = c.limiter.Wait(context.Background())
	return http.Get[PlayerStatisticResponse](c.baseUrl+"/players/statistics?game="+strconv.Itoa(gameId), &c.baseHeaders)
}

func (c Client) PlayerInfo(playerId int, name string, teamId int, season int, country string, search string) (PlayersResponse, error) {
	//todo: ignore zero values
	//
	// For free plan limit is 10 requests/minute
	_ = c.limiter.Wait(context.Background())
	return http.Get[PlayersResponse](c.baseUrl+"/players?id="+strconv.Itoa(playerId), &c.baseHeaders)
}
