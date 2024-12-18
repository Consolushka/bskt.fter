package repository

import (
	"FTER/internal/models"
	"FTER/internal/statistics/repositories/sport_radar/client"
	"FTER/internal/statistics/repositories/sport_radar/dtos"
	"encoding/json"
)

type SportRadarRepository struct {
	client *client.SportRadarClient
}

func (r SportRadarRepository) GameBoxScore(gameId string) (*models.GameModel, error) {
	var gameDto dtos.GameBoxScoreDTO

	homeJSON := r.client.GameSummary(gameId)
	homeRaw, _ := json.Marshal(homeJSON)

	err := json.Unmarshal(homeRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	return gameDto.ToFterModel(), nil
}

func NewSportRadarRepository() *SportRadarRepository {
	return &SportRadarRepository{
		client: client.NewSportRadarClient(),
	}
}
