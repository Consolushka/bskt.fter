package infobasket

import (
	"IMP/app/pkg/http"
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

const (
	boxScoreEndpointPattern  = "/Widget/GetOnline/%v?format=json&lang=ru"
	teamGamesEndpointPattern = "/Widget/TeamGames/%v?format=json"
	scheduleEndpointPattern  = "/Comp/GetCalendar/?comps=%v&format=json"
)

type ClientInterface interface {
	BoxScore(gameId string) (GameBoxScoreResponse, error)
	ScheduledGames(compId int) ([]GameScheduleDto, error)
}

type Client struct {
	baseUrl string
	limiter *rate.Limiter
}

func (c *Client) BoxScore(gameId string) (GameBoxScoreResponse, error) {
	_ = c.limiter.Wait(context.Background())
	return http.Get[GameBoxScoreResponse](c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)
}

func (c *Client) ScheduledGames(compId int) ([]GameScheduleDto, error) {
	_ = c.limiter.Wait(context.Background())
	return http.Get[[]GameScheduleDto](c.baseUrl+fmt.Sprintf(scheduleEndpointPattern, compId), nil)
}

func NewInfobasketClient(leadHost string, rateLimit int) ClientInterface {
	return &Client{
		baseUrl: "https://" + leadHost + ".infobasket.su",
		limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1),
	}
}
