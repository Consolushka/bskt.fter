package sportoteka

import (
	"IMP/app/pkg/http"
	"fmt"
	"time"
)

type ClientInterface interface {
	BoxScore(gameId string) (GameBoxScoreResponse, error)
	Calendar(tag string, year int, from time.Time, to time.Time) (CalendarResponse, error)
}

type Client struct {
	baseUrl string
}

func (c Client) Calendar(tag string, year int, from time.Time, to time.Time) (CalendarResponse, error) {
	return http.Get[CalendarResponse](fmt.Sprintf(c.baseUrl+"/api/abc/comps/calendar?tag=%s&season=%d&maxResultCount=1000&from=%s&to=%s", tag, year, from.Format("2006-01-02"), to.Format("2006-01-02")), nil)
}

func (c Client) BoxScore(gameId string) (GameBoxScoreResponse, error) {
	return http.Get[GameBoxScoreResponse](fmt.Sprintf(c.baseUrl+"/api/abc/games/game?id=%s&lang=en", gameId), nil)
}

func NewClient() ClientInterface {
	return Client{
		baseUrl: "https://basket.sportoteka.org",
	}
}
