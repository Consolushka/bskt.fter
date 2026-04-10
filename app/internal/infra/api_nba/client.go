package api_nba

import (
	"IMP/app/pkg/http"
	"context"
	"net/url"
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
	if c.limiter != nil {
		_ = c.limiter.Wait(context.Background())
	}

	params := url.Values{}
	if id != 0 {
		params.Add("id", strconv.Itoa(id))
	}
	if date != "" {
		params.Add("date", date)
	}
	if leagueId != "" {
		params.Add("league", leagueId)
	}
	if season != "" {
		params.Add("season", season)
	}
	if teamId != "" {
		params.Add("team", teamId)
	}
	if timezone != "" {
		params.Add("timezone", timezone)
	}

	query := ""
	if len(params) > 0 {
		query = "?" + params.Encode()
	}

	return http.Get[GamesResponse](c.baseUrl+"/games"+query, &c.baseHeaders)
}

func (c Client) PlayersStatistics(playerId int, gameId int, teamId int, season string) (PlayerStatisticResponse, error) {
	// For free plan limit is 10 requests/minute
	if c.limiter != nil {
		_ = c.limiter.Wait(context.Background())
	}

	params := url.Values{}
	if playerId != 0 {
		params.Add("id", strconv.Itoa(playerId))
	}
	if gameId != 0 {
		params.Add("game", strconv.Itoa(gameId))
	}
	if teamId != 0 {
		params.Add("team", strconv.Itoa(teamId))
	}
	if season != "" {
		params.Add("season", season)
	}

	query := ""
	if len(params) > 0 {
		query = "?" + params.Encode()
	}

	return http.Get[PlayerStatisticResponse](c.baseUrl+"/players/statistics"+query, &c.baseHeaders)
}

func (c Client) PlayerInfo(playerId int, name string, teamId int, season int, country string, search string) (PlayersResponse, error) {
	// For free plan limit is 10 requests/minute
	if c.limiter != nil {
		_ = c.limiter.Wait(context.Background())
	}

	params := url.Values{}
	if playerId != 0 {
		params.Add("id", strconv.Itoa(playerId))
	}
	if name != "" {
		params.Add("name", name)
	}
	if teamId != 0 {
		params.Add("team", strconv.Itoa(teamId))
	}
	if season != 0 {
		params.Add("season", strconv.Itoa(season))
	}
	if country != "" {
		params.Add("country", country)
	}
	if search != "" {
		params.Add("search", search)
	}

	query := ""
	if len(params) > 0 {
		query = "?" + params.Encode()
	}

	return http.Get[PlayersResponse](c.baseUrl+"/players"+query, &c.baseHeaders)
}
