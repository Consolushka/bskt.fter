package infobasket

import (
	"IMP/app/internal/infrastructure/infobasket/dtos/boxscore"
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/players"
	"IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/teams"
	"time"
)

type persistenceService struct {
	teamsRepository   *teams.Repository
	playersRepository *players.Repository
	gamesRepository   *games.Repository
	leaguesRepository *leagues.Repository
}

func (p *persistenceService) savePlayerModel(player boxscore.PlayerBoxscore) *players.Player {
	birthDate, _ := time.Parse("02.01.2006", player.PersonBirth)

	playerModel, _ := p.playersRepository.FirstOrCreate(players.Player{
		FullName:       player.PersonNameRu,
		BirthDate:      &birthDate,
		LeaguePlayerID: player.PersonID,
	})

	return playerModel
}

func (p *persistenceService) saveTeamPlayers(teamDto boxscore.TeamBoxscore, gameModel games.GameModel, teamModel teams.TeamModel) {
	for _, player := range teamDto.Players {
		playerModel := p.savePlayerModel(player)

		err := p.playersRepository.FirstOrCreateGameStat(players.PlayerGameStats{
			PlayerID:      playerModel.ID,
			GameID:        gameModel.ID,
			TeamID:        teamModel.ID,
			PlayedSeconds: player.Seconds,
			PlsMin:        player.PlusMinus,
			IsBench:       !player.IsStart,
		})

		if err != nil {
			panic(err)
		}
	}
}

func (p *persistenceService) saveTeam(dto boxscore.TeamBoxscore, leagueId int) teams.TeamModel {
	teamModel, _ := p.teamsRepository.FirstOrCreate(teams.TeamModel{
		Alias:    dto.TeamName.CompTeamAbcNameEn,
		LeagueID: leagueId,
		Name:     dto.TeamName.CompTeamNameRu,
	})

	return teamModel
}

func (p *persistenceService) saveGame(gameDto boxscore.GameInfo) games.GameModel {
	league := enums.MLBL
	homeTeam := gameDto.GameTeams[0]
	awayTeam := gameDto.GameTeams[1]

	// query to get NBA league id
	leagueId := p.getLeagueId()

	// save and get home team
	homeTeamModel := p.saveTeam(homeTeam, leagueId)

	// save and get away team
	awayTeamModel := p.saveTeam(awayTeam, leagueId)

	// calculate full game duration
	duration := 0
	duration = 4 * league.QuarterDuration()
	for i := 0; i < gameDto.MaxPeriod-4; i++ {
		duration += league.OvertimeDuration()
	}

	scheduled, _ := time.Parse("2006-01-02 23.10", gameDto.GameDate+" "+gameDto.GameTime)
	gameModel, _ := p.gamesRepository.FirstOrCreate(games.GameModel{
		HomeTeamID:    homeTeamModel.ID,
		AwayTeamID:    awayTeamModel.ID,
		LeagueID:      leagueId,
		ScheduledAt:   scheduled,
		PlayedMinutes: duration,
	})

	// save each player (if not exists) and save player statistics
	p.saveTeamPlayers(homeTeam, gameModel, homeTeamModel)

	p.saveTeamPlayers(awayTeam, gameModel, awayTeamModel)

	return gameModel
}

func (p *persistenceService) getLeagueId() int {
	league, _ := p.leaguesRepository.GetLeagueByAliasEn("mlbl")

	return league.ID
}

func newPersistenceService() *persistenceService {
	return &persistenceService{
		teamsRepository:   teams.NewRepository(),
		playersRepository: players.NewRepository(),
		gamesRepository:   games.NewRepository(),
		leaguesRepository: leagues.NewRepository(),
	}
}
