package nba_com

import (
	"IMP/app/internal/abstract/http"
	"IMP/app/internal/infrastructure/crawler"
	"github.com/PuerkitoBio/goquery"
)

const playerInfoPage = "/player/"

type Client struct {
	baseUrl    string
	httpClient *http.Client
}

func (c *Client) PlayerInfoPage(playerId string) *goquery.Document {
	url := c.baseUrl + playerInfoPage + playerId

	doc, err := crawler.FetchPage(url)
	if err != nil {
		panic(err)
	}

	return doc
}

func NewClient() *Client {
	return &Client{
		baseUrl:    "https://www.nba.com",
		httpClient: http.NewHttpClient(),
	}
}
