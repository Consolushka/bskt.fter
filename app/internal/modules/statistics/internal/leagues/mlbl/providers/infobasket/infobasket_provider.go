package infobasket

import (
	"IMP/app/internal/infrastructure/infobasket"
	models2 "IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
	"strconv"
	"time"
)

type Provider struct {
	client *infobasket.Client
	mapper *mapper
}

func (i *Provider) GameBoxScore(gameId string) (*models2.GameBoxScoreDTO, error) {
	gameDto := i.client.BoxScore(gameId)

	game := i.mapper.mapGame(gameDto)
	game.Id = gameId
	return game, nil
}

func (i *Provider) GamesByDate(date time.Time) ([]string, error) {
	var result []string
	compIds := []int{89960, 89962}

	for _, compId := range compIds {
		seasonGames := i.client.ScheduledGames(compId)

		for _, game := range seasonGames {
			if game.GameDate == date.Format("02.01.2006") {
				result = append(result, strconv.Itoa(game.GameID))
			}
		}
	}

	return result, nil
}

func (i *Provider) GamesByTeam(teamId string) ([]string, error) {
	scheduleJson := i.client.TeamGames(teamId)

	gamesIds := array_utils.Map(scheduleJson.Games, func(game infobasket.GameScheduleDto) string {
		if game.GameStatus == 1 {
			return strconv.Itoa(game.GameID)
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
