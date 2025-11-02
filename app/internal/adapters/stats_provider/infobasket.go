package stats_provider

import (
	"IMP/app/internal/adapters/cached_remote_resource"
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/infra/infobasket"
	"IMP/app/internal/service/remote_cache_loader"
	"IMP/app/pkg/logger"
	"strconv"
	"time"
)

type InfobasketStatsProviderAdapter struct {
	client infobasket.ClientInterface

	transformer infobasket.EntityTransformer
	compId      int
}

func (i InfobasketStatsProviderAdapter) GetPlayerBio(id string) (players.PlayerBioEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (i InfobasketStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	var gamesEntities []games.GameStatEntity

	cacher := cached_remote_resource.NewInfobasketCachedResource(i.client.ScheduledGames, i.compId)
	schedule, err := remote_cache_loader.LoadLocalFile[[]infobasket.GameScheduleDto](cacher)
	if err != nil {
		return nil, err
	}

	for _, game := range schedule {
		// game haven't been scheduled yet
		if game.GameTime == "--:--" {
			continue
		}

		gameDate, err := time.Parse("02.01.2006 15:04", game.GameDate+" "+game.GameTime)
		if err != nil {
			logger.Error("There was an error while parsing game gameDate", map[string]interface{}{
				"gameDate": game.GameDate,
				"error":    err,
			})
			continue
		}
		if gameDate.After(from) && gameDate.Before(to) {
			gameBoxScore, err := i.client.BoxScore(strconv.Itoa(game.GameID))
			if err != nil {
				logger.Error("There was an error while fetching game box score", map[string]interface{}{
					"gameId": game.GameID,
					"error":  err,
				})
				continue
			}
			if gameBoxScore.GameStatus != 1 {
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
