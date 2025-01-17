package sport_radar

import (
	"IMP/app/internal/infrastructure/sport_radar"
	models2 "IMP/app/internal/modules/statistics/models"
)

type Repository struct {
	client *sport_radar.Client
}

func (r Repository) TodayGames() (string, []string, error) {
	panic("implement me")
}

func (r Repository) GameBoxScore(gameId string) (*models2.GameBoxScoreDTO, error) {
	panic("implement me")
}

func NewRepository() *Repository {
	return &Repository{
		client: sport_radar.NewSportRadarApiClient(),
	}
}
