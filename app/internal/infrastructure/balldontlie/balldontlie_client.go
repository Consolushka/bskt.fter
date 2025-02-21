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
	headers := make(map[string]string)
	headers["Authorization"] = c.apiKey

	result := http.Get[PlayersResponse](c.baseUrl+"/players?first_name="+firstNameSearch+"&last_name="+lastNameSearch, &headers)

	return result.Data[0]
}

func NewClient() *Client {
	return &Client{
		baseUrl: "https://api.balldontlie.io/v1",
		apiKey:  os.Getenv("BALLDONTLIE_API_KEY"),
	}
}
