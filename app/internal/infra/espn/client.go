package espn

import (
	"IMP/app/pkg/http"
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

const (
	siteBaseUrl    = "https://site.api.espn.com/apis/site/v2/sports/basketball/nba"
	coreBaseUrl    = "https://sports.core.api.espn.com/v2/sports/basketball/leagues/nba"
	requestTimeout = 30 * time.Second
)

type ClientInterface interface {
	// Scoreboard: date — YYYYMMDD, календарный день US Eastern.
	Scoreboard(date string) (ScoreboardResponse, error)
	Summary(eventId string) (SummaryResponse, error)
	Athlete(id string) (AthleteResponse, error)
}

type Client struct {
	siteBaseUrl string
	coreBaseUrl string
	limiter     *rate.Limiter
	headers     map[string]string
}

func NewClient(rateLimit int) ClientInterface {
	// При нуле time.Minute/0 паникует (как в api_nba), поэтому без лимита — бесконечный лимитер.
	limiter := rate.NewLimiter(rate.Inf, 1)
	if rateLimit > 0 {
		limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(rateLimit)), 1)
	}

	return &Client{
		siteBaseUrl: siteBaseUrl,
		coreBaseUrl: coreBaseUrl,
		limiter:     limiter,
		headers:     map[string]string{"User-Agent": "Argus/1.0 (+bskt.imp)"},
	}
}

func (c Client) Scoreboard(date string) (ScoreboardResponse, error) {
	resp, err := get[ScoreboardResponse](c, c.siteBaseUrl+"/scoreboard?dates="+date)
	if err != nil {
		return resp, fmt.Errorf("scoreboard %s returned error: %w", date, err)
	}
	return resp, nil
}

func (c Client) Summary(eventId string) (SummaryResponse, error) {
	resp, err := get[SummaryResponse](c, c.siteBaseUrl+"/summary?event="+eventId)
	if err != nil {
		return resp, fmt.Errorf("summary %s returned error: %w", eventId, err)
	}
	return resp, nil
}

func (c Client) Athlete(id string) (AthleteResponse, error) {
	resp, err := get[AthleteResponse](c, c.coreBaseUrl+"/athletes/"+id)
	if err != nil {
		return resp, fmt.Errorf("athlete %s returned error: %w", id, err)
	}
	return resp, nil
}

func get[T any](c Client, url string) (T, error) {
	// Ожидание лимитера — на фоновом контексте, чтобы оно не съедало таймаут запроса.
	_ = c.limiter.Wait(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	return http.GetWithContext[T](ctx, url, &c.headers)
}
