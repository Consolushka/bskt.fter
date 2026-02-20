package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/infra/sportoteka"
	"strconv"
	"time"

	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
)

type SportotekaStatsProviderAdapter struct {
	tag  string
	year int

	client      sportoteka.ClientInterface
	transformer sportoteka.EntityTransformer
}

func (s SportotekaStatsProviderAdapter) GetPlayerBio(id string) (players.PlayerBioEntity, error) {
	// Sportoteka doesn't provide detailed player bio by ID in current client
	return players.PlayerBioEntity{}, nil
}

func (s SportotekaStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	calendar, err := s.client.Calendar(s.tag, s.year, from, to)
	if err != nil {
		return make([]games.GameStatEntity, 0), err
	}

	gamesEntities := make([]games.GameStatEntity, 0, calendar.TotalCount)

	for _, calendarGame := range calendar.Items {
		if calendarGame.Game.GameStatus != "Result" && calendarGame.Game.GameStatus != "ResultConfirmed" {
			continue
		}
		gameBoxScore, err := s.client.BoxScore(strconv.Itoa(calendarGame.Game.Id))
		if err != nil {
			compositelogger.Error("There was an error while fetching game box score", map[string]interface{}{
				"gameId": calendarGame.Game.Id,
				"error":  err,
			})
			continue
		}

		if gameBoxScore.Result.Game.GameStatus != "ResultConfirmed" && gameBoxScore.Result.Game.GameStatus != "Result" {
			continue
		}

		entity, err := s.transformer.Transform(gameBoxScore.Result)
		if err != nil {
			compositelogger.Error("There was an error while transforming game box score", map[string]interface{}{
				"gameId": calendarGame.Game.Id,
				"error":  err,
			})
			continue
		}

		entity.GameModel.Title = calendarGame.Team1.AbcName + " - " + calendarGame.Team2.AbcName
		entity.HomeTeamStat.TeamModel.Name = calendarGame.Team1.Name
		entity.AwayTeamStat.TeamModel.Name = calendarGame.Team2.Name

		gamesEntities = append(gamesEntities, entity)
	}

	return gamesEntities, nil
}

func (s SportotekaStatsProviderAdapter) EnrichGameStats(game games.GameStatEntity) (games.GameStatEntity, error) {
	return game, nil
}

func NewSportotekaStatsProvider(client sportoteka.ClientInterface, tag string, year int) SportotekaStatsProviderAdapter {
	return SportotekaStatsProviderAdapter{
		tag:         tag,
		year:        year,
		client:      client,
		transformer: sportoteka.NewEntityTransformer(),
	}
}
