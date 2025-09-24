package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/infra/infobasket"
	"IMP/app/pkg/logger"
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
				logger.Error("There was an error while fetching game box score", map[string]interface{}{
					"gameId": game.GameID,
					"error":  err,
				})
				continue
			}

			transform, err := i.transformer.Transform(gameBoxScore)
			if err != nil {
				logger.Error("There was an error while transforming game box score", map[string]interface{}{
					"gameId":       game.GameID,
					"gameBoxScore": gameBoxScore,
					"error":        err,
				})
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
