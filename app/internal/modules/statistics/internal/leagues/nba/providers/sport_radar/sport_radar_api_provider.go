package sport_radar

import (
	"IMP/app/internal/infrastructure/sport_radar"
	models2 "IMP/app/internal/modules/statistics/models"
)

type Provider struct {
	client *sport_radar.Client
}

func (r Provider) TodayGames() (string, []string, error) {
	panic("implement me")
}

func (r Provider) GameBoxScore(gameId string) (*models2.GameBoxScoreDTO, error) {
	panic("implement me")
}

func NewProvider() *Provider {
	return &Provider{
		client: sport_radar.NewSportRadarApiClient(),
	}
}
