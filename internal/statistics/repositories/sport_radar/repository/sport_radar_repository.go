package repository

import (
	"NBATrueEfficency/internal/statistics/repositories/sport_radar/client"
	"NBATrueEfficency/internal/statistics/repositories/sport_radar/dtos"
	"encoding/json"
)

type SportRadarRepository struct {
	client *client.SportRadarClient
}

func (r SportRadarRepository) GetGame(gameId string) (*dtos.GameDTO, error) {
	var gameDto dtos.GameDTO

	homeJSON := r.client.GameSummary(gameId)
	homeRaw, _ := json.Marshal(homeJSON)

	err := json.Unmarshal(homeRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	return &gameDto, nil
}

func NewSportRadarRepository() *SportRadarRepository {
	return &SportRadarRepository{
		client: client.NewSportRadarClient(),
	}
}
