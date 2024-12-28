package infobasket

import (
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/modules/statistics/leagues/mlbl/repositories/infobasket/client"
	boxscore2 "IMP/app/internal/modules/statistics/leagues/mlbl/repositories/infobasket/dtos/boxscore"
	"encoding/json"
)

type Repository struct {
	client *client.InfobasketClient
}

func (i Repository) GameBoxScore(gameId string) (*models.GameModel, error) {
	var gameDto boxscore2.GameInfo

	boxscoreJson := i.client.BoxScore(gameId)
	boxscoreRaw, _ := json.Marshal(boxscoreJson)

	err := json.Unmarshal(boxscoreRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	return gameDto.ToImpModel(), nil
}

func (i Repository) TodayGames() (string, []string, error) {
	panic("implement me")
}

func NewRepository() *Repository {
	return &Repository{
		client: client.NewInfobasketClient(),
	}
}
