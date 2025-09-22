package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/infra/infobasket"
	"fmt"
	"strconv"
	"time"
)

type InfobasketStatsProviderAdapter struct {
	client infobasket.ClientInterface

	transformer infobasket.EntityTransformer
	compId      int
}

func (i InfobasketStatsProviderAdapter) GetGamesStatsByDate(date time.Time) ([]games.GameStatEntity, error) {
	var gamesEntities []games.GameStatEntity
	schedule, err := i.client.ScheduledGames(i.compId)
	if err != nil {
		return nil, err
	}

	for _, game := range schedule {
		if game.GameDate == date.Format("02.01.2006") {
			gameBoxScore, err := i.client.BoxScore(strconv.Itoa(game.GameID))
			if err != nil {
				// todo: log
				fmt.Println("There was an error while fetching game box score", ". GameID: ", strconv.Itoa(game.GameID), ". Error: ", err)
			}

			transform, err := i.transformer.Transform(gameBoxScore)
			if err != nil {
				// todo: log
				fmt.Println("There was an error while transforming game box score", ". GameID: ", strconv.Itoa(game.GameID), ". Error: ", err)
			}

			gamesEntities = append(gamesEntities, transform)
		}
	}
	return gamesEntities, nil
}

func NewInfobasketStatsProviderAdapter(client infobasket.ClientInterface, compId int) InfobasketStatsProviderAdapter {
	return InfobasketStatsProviderAdapter{
		client: client,
		compId: compId,
	}
}
