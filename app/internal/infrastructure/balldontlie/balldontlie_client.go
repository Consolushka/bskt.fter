package balldontlie

import (
	"IMP/app/internal/abstract/http"
	"os"
)

type Client struct {
	baseUrl string
	apiKey  string

	httpClient *http.Client
}

func (c *Client) GetAllPlayers(firstNameSearch string, lastNameSearch string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+"/players?first_name="+firstNameSearch+"&last_name="+lastNameSearch, &c.apiKey)

	//todo: refactor json response
	return result["data"].([]interface{})[0].(map[string]interface{})
}

func NewClient() *Client {
	return &Client{
		baseUrl:    "https://api.balldontlie.io/v1",
		apiKey:     os.Getenv("BALLDONTLIE_API_KEY"),
		httpClient: http.NewHttpClient(),
	}
}
