package repository

import (
	"FTER/app/internal/models"
	"FTER/app/internal/statistics/repositories/sport_radar/client"
	"FTER/app/internal/statistics/repositories/sport_radar/dtos"
	"encoding/json"
)

type SportRadarRepository struct {
	client *client.SportRadarClient
}

func (r SportRadarRepository) TodayGames() (string, []string, error) {
	panic("implement me")
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
