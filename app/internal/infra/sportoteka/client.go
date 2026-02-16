package sportoteka

import (
	"IMP/app/pkg/http"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

type ClientInterface interface {
	BoxScore(gameId string) (GameBoxScoreResponse, error)
	Calendar(tag string, year int, from time.Time, to time.Time) (CalendarResponse, error)
}

type Client struct {
	baseUrl string
	limiter *rate.Limiter
}

func (c Client) Calendar(tag string, year int, from time.Time, to time.Time) (CalendarResponse, error) {
	_ = c.limiter.Wait(context.Background())
	return http.Get[CalendarResponse](fmt.Sprintf(c.baseUrl+"/api/abc/comps/calendar?tag=%s&season=%d&maxResultCount=1000&from=%s&to=%s", tag, year, from.Format("2006-01-02"), to.Format("2006-01-02")), nil)
}

func (c Client) BoxScore(gameId string) (GameBoxScoreResponse, error) {
	_ = c.limiter.Wait(context.Background())
	return http.Get[GameBoxScoreResponse](fmt.Sprintf(c.baseUrl+"/api/abc/games/game?id=%s&lang=en", gameId), nil)
}

func NewClient() ClientInterface {
	rateLimitStr := os.Getenv("SPORTOTEKA_RATE_LIMIT_PER_MINUTE")
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil || rateLimit <= 0 {
		rateLimit = 25
	}

	return Client{
		baseUrl: "https://basket.sportoteka.org",
		limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1),
	}
}
