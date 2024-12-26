package infobasket

import (
	"FTER/app/internal/models"
	"FTER/app/internal/statistics/leagues/mlbl/repositories/infobasket/client"
	boxscore2 "FTER/app/internal/statistics/leagues/mlbl/repositories/infobasket/dtos/boxscore"
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

	return gameDto.ToFterModel(), nil
}

func (i Repository) TodayGames() (string, []string, error) {
	panic("implement me")
}

func NewRepository() *Repository {
	return &Repository{
		client: client.NewInfobasketClient(),
	}
}
