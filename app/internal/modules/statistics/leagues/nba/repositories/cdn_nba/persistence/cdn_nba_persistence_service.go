package persistence

import (
	boxscore2 "IMP/app/internal/infrastructure/cdn_nba/dtos/boxscore"
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/statistics/enums"
	"strings"
)

type CdnNbaPersistenceService struct {
	gameService *gamePersistenceService
	teamService *teamPersistenceService

	leagueRepository *leagues.Repository
}

func (p *CdnNbaPersistenceService) SaveGame(gameDto boxscore2.GameDTO) games.GameModel {
	league := enums.NBA
	leagueModel, _ := p.leagueRepository.GetLeagueByAliasEn(strings.ToLower(league.String()))

	// save and get home team
	homeTeamModel := p.teamService.saveTeam(gameDto.HomeTeam, leagueModel.ID)

	// save and get away team
	awayTeamModel := p.teamService.saveTeam(gameDto.AwayTeam, leagueModel.ID)

	gameModel := p.gameService.saveGame(gameDto, league, leagueModel.ID, homeTeamModel.ID, awayTeamModel.ID)

	// save players statistics
	p.teamService.saveTeamPlayers(gameDto.HomeTeam, gameModel, homeTeamModel)

	p.teamService.saveTeamPlayers(gameDto.AwayTeam, gameModel, awayTeamModel)

	return gameModel
}

func NewCdnNbaPersistenceService() *CdnNbaPersistenceService {
	return &CdnNbaPersistenceService{
		gameService:      newGamePersistenceService(),
		teamService:      newTeamPersistenceService(),
		leagueRepository: leagues.NewRepository(),
	}
}
