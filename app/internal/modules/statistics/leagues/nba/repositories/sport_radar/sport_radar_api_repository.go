package sport_radar

import (
	"IMP/app/internal/infrastructure/sport_radar"
	"IMP/app/internal/infrastructure/sport_radar/dtos"
	"IMP/app/internal/modules/imp/models"
	"encoding/json"
)

type Repository struct {
	client *sport_radar.Client
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

	//todo: save to db
	return gameDto.ToImpModel(), nil
}

func NewRepository() *Repository {
	return &Repository{
		client: sport_radar.NewSportRadarApiClient(),
	}
}
