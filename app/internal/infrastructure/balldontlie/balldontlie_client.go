package balldontlie

import (
	"IMP/app/internal/abstract"
)

type Client struct {
	baseUrl string
	apiKey  string

	httpClient *abstract.HttpClient
}

func (c *Client) GetAllPlayers(firstNameSearch string, lastNameSearch string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+"/players?first_name="+firstNameSearch+"&last_name="+lastNameSearch, &c.apiKey)

	return result["data"].(map[string]interface{})
}

func NewClient() *Client {
	return &Client{
		baseUrl:    "https://api.balldontlie.io/v1",
		apiKey:     "791fa210-2f28-46c7-9265-16209d0ee3e2",
		httpClient: abstract.NewHttpClient(),
	}
}
