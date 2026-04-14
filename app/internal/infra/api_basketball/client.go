package api_basketball

import (
	"IMP/app/pkg/http"
	"context"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

type ClientInterface interface {
	Games(id int, date string, leagueId string, season string, teamId string, timezone string) (GamesResponse, error)
	TeamStatistics(gameId int, teamId int) (TeamStatsResponse, error)
	PlayersStatistics(gameId int, teamId int, playerId int) (PlayerStatsResponse, error)
	PlayerInfo(playerId int) (PlayerInfoResponse, error)
}

type Client struct {
	baseUrl     string
	token       string
	baseHeaders map[string]string
	limiter     *rate.Limiter
}

func NewClient(baseUrl string, token string, rateLimit int) ClientInterface {
	return &Client{
		baseUrl: baseUrl,
		token:   token,
		baseHeaders: map[string]string{
			"x-apisports-key": token,
		},
		limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1),
	}
}

func (c *Client) Games(id int, date string, leagueId string, season string, teamId string, timezone string) (GamesResponse, error) {
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

func (c *Client) TeamStatistics(gameId int, teamId int) (TeamStatsResponse, error) {
	if c.limiter != nil {
		_ = c.limiter.Wait(context.Background())
	}

	params := url.Values{}
	if gameId != 0 {
		params.Add("id", strconv.Itoa(gameId))
	}
	if teamId != 0 {
		params.Add("team", strconv.Itoa(teamId))
	}

	query := ""
	if len(params) > 0 {
		query = "?" + params.Encode()
	}

	return http.Get[TeamStatsResponse](c.baseUrl+"/games/statistics/teams"+query, &c.baseHeaders)
}

func (c *Client) PlayersStatistics(gameId int, teamId int, playerId int) (PlayerStatsResponse, error) {
	if c.limiter != nil {
		_ = c.limiter.Wait(context.Background())
	}

	params := url.Values{}
	if gameId != 0 {
		params.Add("id", strconv.Itoa(gameId))
	}
	if teamId != 0 {
		params.Add("team", strconv.Itoa(teamId))
	}
	if playerId != 0 {
		params.Add("player", strconv.Itoa(playerId))
	}

	query := ""
	if len(params) > 0 {
		query = "?" + params.Encode()
	}

	return http.Get[PlayerStatsResponse](c.baseUrl+"/games/statistics/players"+query, &c.baseHeaders)
}

func (c *Client) PlayerInfo(playerId int) (PlayerInfoResponse, error) {
	if c.limiter != nil {
		_ = c.limiter.Wait(context.Background())
	}

	params := url.Values{}
	if playerId != 0 {
		params.Add("id", strconv.Itoa(playerId))
	}

	query := ""
	if len(params) > 0 {
		query = "?" + params.Encode()
	}

	return http.Get[PlayerInfoResponse](c.baseUrl+"/players"+query, &c.baseHeaders)
}
