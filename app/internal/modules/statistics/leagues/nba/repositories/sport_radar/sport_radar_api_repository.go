package sport_radar

import (
	"FTER/app/internal/modules/imp/models"
	"FTER/app/internal/modules/statistics/leagues/nba/repositories/sport_radar/client"
	"FTER/app/internal/modules/statistics/leagues/nba/repositories/sport_radar/dtos"
	"encoding/json"
)

type Repository struct {
	client *client.SportRadarApiClient
}

func (r Repository) TodayGames() (string, []string, error) {
	panic("implement me")
}

func (r Repository) GameBoxScore(gameId string) (*models.GameModel, error) {
	var gameDto dtos.GameBoxScoreDTO

	homeJSON := r.client.GameSummary(gameId)
	homeRaw, _ := json.Marshal(homeJSON)

	err := json.Unmarshal(homeRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	return gameDto.ToFterModel(), nil
}

func NewRepository() *Repository {
	return &Repository{
		client: client.NewSportRadarApiClient(),
	}
}
