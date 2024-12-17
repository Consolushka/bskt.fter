package repository

import (
	"FTER/internal/models"
	"FTER/internal/statistics/repositories/nba/client"
	"FTER/internal/statistics/repositories/nba/dtos"
	"encoding/json"
)

type NbaRepository struct {
	client *client.NbaClient
}

func (n *NbaRepository) GameBoxScore(gameId string) (*models.GameModel, error) {
	var gameDto dtos.GameDTO

	homeJSON := n.client.BoxScore(gameId)
	homeRaw, _ := json.Marshal(homeJSON)

	err := json.Unmarshal(homeRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	return gameDto.ToFterModel(), nil
}

func NewNbaRepository() *NbaRepository {
	return &NbaRepository{
		client: client.NewNbaClient(),
	}
}
