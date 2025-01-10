package nba_com

import (
	"IMP/app/internal/abstract"
	"IMP/app/internal/infrastructure/crawler"
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

const playerInfoPage = "/player/"

type Client struct {
	baseUrl    string
	httpClient *abstract.HttpClient
}

func (c *Client) PlayerInfoPage(playerId int) *goquery.Document {
	url := c.baseUrl + playerInfoPage + strconv.Itoa(playerId)

	doc, err := crawler.FetchPage(url)
	if err != nil {
		panic(err)
	}

	return doc
}

func NewClient() *Client {
	return &Client{
		baseUrl:    "https://www.nba.com",
		httpClient: abstract.NewHttpClient(),
	}
}
