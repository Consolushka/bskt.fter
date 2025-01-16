package persistence

import (
	boxscore2 "IMP/app/internal/infrastructure/cdn_nba/dtos/boxscore"
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/teams"
)

type teamPersistenceService struct {
	playerPersistenceService *playerPersistenceService

	teamsRepository *teams.Repository
}

// saveTeam returns or create team
func (p *teamPersistenceService) saveTeam(dto boxscore2.TeamDTO, leagueId int) teams.Team {
	teamModel, _ := p.teamsRepository.FirstOrCreate(teams.Team{
		Alias:    dto.TeamTricode,
		LeagueID: leagueId,
		Name:     dto.TeamName,
	})

	return teamModel
}

// saveTeamPlayers save players and their statistics for a team
func (p *teamPersistenceService) saveTeamPlayers(teamDto boxscore2.TeamDTO, gameModel games.GameModel, teamModel teams.Team) {
	for _, player := range teamDto.Players {
		p.playerPersistenceService.savePlayerGameStats(player, gameModel, teamModel)
	}
}

func newTeamPersistenceService() *teamPersistenceService {
	return &teamPersistenceService{
		playerPersistenceService: newPlayerPersistenceService(),
		teamsRepository:          teams.NewRepository(),
	}
}
