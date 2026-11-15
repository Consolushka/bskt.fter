package stats_provider

import (
	"IMP/app/internal/adapters/cached_remote_resource"
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/infra/infobasket"
	"IMP/app/internal/ports"
	"IMP/app/internal/service/remote_cache_loader"
	"strconv"
	"time"

	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
)

// Начало обязано совпадать с именем провайдера INFOBASKET в tournament_providers.provider_name.
// Константу из service/providers взять нельзя: тот пакет сам импортирует этот, получился бы цикл.
const teamIdScopePrefix = "INFOBASKET:"

var _ ports.TeamIdScopeProvider = (*InfobasketStatsProviderAdapter)(nil)

type InfobasketStatsProviderAdapter struct {
	client infobasket.ClientInterface

	transformer infobasket.EntityTransformer
	compId      int
	leadHost    string
}

func (i InfobasketStatsProviderAdapter) GetPlayerBio(id string) (players.PlayerBioEntity, error) {
	// Infobasket doesn't provide detailed player bio by ID in current client
	// Return empty bio for now to satisfy interface and avoid panics
	return players.PlayerBioEntity{}, nil
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
			compositelogger.Error("There was an error while parsing game gameDate", map[string]interface{}{
				"gameDate": game.GameDate,
				"error":    err,
			})
			continue
		}
		if (gameDate.After(from) || gameDate.Equal(from)) && (gameDate.Before(to) || gameDate.Equal(to)) {
			gameBoxScore, err := i.client.BoxScore(strconv.Itoa(game.GameID))
			if err != nil {
				compositelogger.Error("There was an error while fetching game box score", map[string]interface{}{
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
				compositelogger.Error("There was an error while transforming game box score", map[string]interface{}{
					"gameId":       game.GameID,
					"gameBoxScore": gameBoxScore,
					"error":        err,
				})
				continue
			}

			gamesEntities = append(gamesEntities, transform)
		}
	}
	return gamesEntities, nil
}

func (i InfobasketStatsProviderAdapter) EnrichGameStats(game games.GameStatEntity) (games.GameStatEntity, error) {
	return game, nil
}

func NewInfobasketStatsProviderAdapter(client infobasket.ClientInterface, compId int, leadHost string) InfobasketStatsProviderAdapter {
	return InfobasketStatsProviderAdapter{
		client:      client,
		transformer: infobasket.NewEntityTransformer(),
		compId:      compId,
		leadHost:    leadHost,
	}
}

// TeamIdScope сообщает, чьи id команд отдаёт адаптер. У Infobasket своя база на каждом хосте,
// и один и тот же id на двух хостах может означать разные команды.
func (i InfobasketStatsProviderAdapter) TeamIdScope() string {
	return teamIdScopePrefix + i.leadHost
}
