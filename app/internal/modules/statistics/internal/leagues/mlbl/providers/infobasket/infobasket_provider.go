package infobasket

import (
	"IMP/app/internal/infrastructure/infobasket"
	"IMP/app/internal/infrastructure/infobasket/dtos/boxscore"
	models2 "IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
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

func (i *Provider) GamesByTeam(teamId string) ([]string, error) {
	scheduleJson := i.client.TeamGames(teamId)

	gamesIds := array_utils.Map(scheduleJson, func(game map[string]interface{}) string {
		if game["GameStatus"].(float64) == 1 {
			return strconv.FormatFloat(game["GameID"].(float64), 'f', 0, 64)
		}

		return ""
	})

	return array_utils.Filter(gamesIds, func(gameId string) bool {
		return gameId != ""
	}), nil
}

func NewProvider() *Provider {
	return &Provider{
		client: infobasket.NewInfobasketClient(),
		mapper: newMapper(),
	}
}
