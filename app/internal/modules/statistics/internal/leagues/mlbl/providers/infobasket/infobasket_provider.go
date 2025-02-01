package infobasket

import (
	"IMP/app/internal/infrastructure/infobasket"
	"IMP/app/internal/infrastructure/infobasket/dtos/boxscore"
	models2 "IMP/app/internal/modules/statistics/models"
	"encoding/json"
	"strconv"
	"time"
)

type Provider struct {
	client *infobasket.Client
	mapper *mapper
}

func (i *Provider) GameBoxScore(gameId string) (*models2.GameBoxScoreDTO, error) {
	var gameDto boxscore.GameInfo

	boxscoreJson := i.client.BoxScore(gameId)
	boxscoreRaw, _ := json.Marshal(boxscoreJson)

	err := json.Unmarshal(boxscoreRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	game := i.mapper.mapGame(gameDto)
	game.Id = gameId
	return game, nil
}

func (i *Provider) GamesByDate(date time.Time) ([]string, error) {
	panic("implement me")
}

func NewProvider() *Provider {
	return &Provider{
		client: infobasket.NewInfobasketClient(),
		mapper: newMapper(),
	}
}
