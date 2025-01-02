package infobasket

import (
	"IMP/app/internal/infrastructure/infobasket"
	boxscore2 "IMP/app/internal/infrastructure/infobasket/dtos/boxscore"
	"IMP/app/internal/modules/imp/models"
	"encoding/json"
)

type Repository struct {
	client *infobasket.Client
}

func (i Repository) GameBoxScore(gameId string) (*models.GameModel, error) {
	var gameDto boxscore2.GameInfo

	boxscoreJson := i.client.BoxScore(gameId)
	boxscoreRaw, _ := json.Marshal(boxscoreJson)

	err := json.Unmarshal(boxscoreRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	//todo: save to db
	return gameDto.ToImpModel(), nil
}

func (i Repository) TodayGames() (string, []string, error) {
	panic("implement me")
}

func NewRepository() *Repository {
	return &Repository{
		client: infobasket.NewInfobasketClient(),
	}
}
