package balldontlie

import (
	"IMP/app/internal/abstract/http"
	"os"
)

type Client struct {
	baseUrl string
	apiKey  string
}

func (c *Client) GetAllPlayers(firstNameSearch string, lastNameSearch string) Player {
	result := http.Get[PlayersResponse](c.baseUrl+"/players?first_name="+firstNameSearch+"&last_name="+lastNameSearch, &c.apiKey)

	return result.Data[0]
}

func NewClient() *Client {
	return &Client{
		baseUrl: "https://api.balldontlie.io/v1",
		apiKey:  os.Getenv("BALLDONTLIE_API_KEY"),
	}
}
