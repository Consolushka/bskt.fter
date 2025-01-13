package persistence

import (
	boxscore2 "IMP/app/internal/infrastructure/cdn_nba/dtos/boxscore"
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/statistics/enums"
)

type gamePersistenceService struct {
	gamesRepository  *games.Repository
	leagueRepository *leagues.Repository
}

// saveGame saves a game if not exists
func (p *gamePersistenceService) saveGame(gameDto boxscore2.GameDTO, league enums.League, leagueId int, homeTeamId int, awayTeamId int) games.GameModel {
	// calculate full game duration
	duration := 0
	duration = 4 * league.QuarterDuration()
	for i := 5; i < gameDto.Period; i++ {
		duration += league.OvertimeDuration()
	}

	gameModel, _ := p.gamesRepository.FirstOrCreate(games.GameModel{
		HomeTeamID:    homeTeamId,
		AwayTeamID:    awayTeamId,
		LeagueID:      leagueId,
		ScheduledAt:   gameDto.GameTimeUTC,
		PlayedMinutes: duration,
	})

	return gameModel
}

func newGamePersistenceService() *gamePersistenceService {
	return &gamePersistenceService{
		gamesRepository:  games.NewRepository(),
		leagueRepository: leagues.NewRepository(),
	}
}
