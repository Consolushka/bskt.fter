package nba_com_api

import (
	"IMP/app/internal/infrastructure/balldontlie"
	boxscore2 "IMP/app/internal/infrastructure/nba_com_api/dtos/boxscore"
	"IMP/app/internal/infrastructure/translator"
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/players"
	"IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/teams"
	"IMP/app/internal/utils/string_utils"
	"IMP/app/internal/utils/time_utils"
)

type PersistenceService struct {
	ballDontLieClient *balldontlie.Client
}

func (p *PersistenceService) savePlayerModel(player boxscore2.PlayerDTO) players.Player {
	// If player name has non-latin characters, translate it to english
	if string_utils.HasNonLanguageChars(player.FirstName, string_utils.Latin) {
		player.FirstName = translator.Translate(player.FirstName, nil, "en")
	}
	if string_utils.HasNonLanguageChars(player.FamilyName, string_utils.Latin) {
		player.FamilyName = translator.Translate(player.FamilyName, nil, "en")
		if player.FamilyName == "Chancar" {
			player.FamilyName = "Cancar"
		}
	}

	// fetch player info from balldontlie to get draft year
	playerInfo := p.ballDontLieClient.GetAllPlayers(player.FirstName, player.FamilyName)
	var draftYear *int
	draftYearFloat, ok := playerInfo["draft_year"]
	if ok && draftYearFloat != nil {
		draftYearInt := int(draftYearFloat.(float64))
		draftYear = &draftYearInt
	}

	playerModel, _ := players.FirstOrCreate(players.Player{
		FullName:  player.Name,
		BirthDate: nil,
		DraftYear: draftYear,
	})

	return playerModel
}

func (p *PersistenceService) saveTeamPlayers(teamDto boxscore2.TeamDTO, gameModel games.GameModel, teamModel teams.TeamModel) {
	for _, player := range teamDto.Players {
		playerModel := p.savePlayerModel(player)

		players.FirstOrCreateGameStat(players.PlayerGameStats{
			PlayerID:      playerModel.ID,
			GameID:        gameModel.ID,
			TeamID:        teamModel.ID,
			PlayedSeconds: time_utils.FormattedMinutesToSeconds(player.Statistics.Minutes, playedTimeFormat),
			PlsMin:        player.Statistics.Plus - player.Statistics.Minus,
			IsBench:       player.Starter != "1",
		})
	}
}

func (p *PersistenceService) saveTeam(dto boxscore2.TeamDTO, leagueId int) teams.TeamModel {
	teamModel, _ := teams.FirstOrCreate(teams.TeamModel{
		Alias:    dto.TeamTricode,
		LeagueID: leagueId,
		Name:     dto.TeamName,
	})

	return teamModel
}

func (p *PersistenceService) saveGame(gameDto boxscore2.GameDTO) games.GameModel {
	league := enums.NBA

	// query to get NBA league id
	leagueId := p.getNbaLeagueId()

	// save and get home team
	homeTeamModel := p.saveTeam(gameDto.HomeTeam, leagueId)

	// save and get away team
	awayTeamModel := p.saveTeam(gameDto.AwayTeam, leagueId)

	// calculate full game duration
	duration := 0
	duration = 4 * league.QuarterDuration()
	for i := 5; i < gameDto.Period; i++ {
		duration += league.OvertimeDuration()
	}

	gameModel, _ := games.FirstOrCreate(games.GameModel{
		HomeTeamID:    homeTeamModel.ID,
		AwayTeamID:    awayTeamModel.ID,
		LeagueID:      leagueId,
		ScheduledAt:   gameDto.GameTimeUTC,
		PlayedMinutes: duration,
	})

	// save each player (if not exists) and save player statistics
	p.saveTeamPlayers(gameDto.HomeTeam, gameModel, homeTeamModel)

	p.saveTeamPlayers(gameDto.AwayTeam, gameModel, awayTeamModel)

	return gameModel
}

func (p *PersistenceService) getNbaLeagueId() int {
	league, _ := leagues.LeagueByAliasEn("nba")

	return league.ID
}

func NewPersistenceService() *PersistenceService {
	return &PersistenceService{
		ballDontLieClient: balldontlie.NewClient(),
	}
}
