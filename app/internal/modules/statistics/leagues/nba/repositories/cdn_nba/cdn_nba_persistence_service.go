package cdn_nba

import (
	boxscore2 "IMP/app/internal/infrastructure/cdn_nba/dtos/boxscore"
	"IMP/app/internal/infrastructure/nba_com"
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/players"
	"IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/teams"
	"IMP/app/internal/utils/time_utils"
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

type persistenceService struct {
	nbaComClient *nba_com.Client

	teamsRepository   *teams.Repository
	playersRepository *players.Repository
	gamesRepository   *games.Repository
	leagueRepository  *leagues.Repository
}

func (p *persistenceService) savePlayerModel(player boxscore2.PlayerDTO) players.Player {

	playerModel, err := p.playersRepository.FirstByLeaguePlayerId(player.PersonId)
	if playerModel == nil {
		log.Println("Player not found in database: ", player.PersonId, ". Fetching from nba.com")

		playerModel, err = p.playersRepository.FirstOrCreate(players.Player{
			FullName:       player.Name,
			BirthDate:      p.getPlayerBirthdate(player),
			LeaguePlayerID: player.PersonId,
		})
	}

	if err != nil {
		panic(err)
	}

	return *playerModel
}

// getPlayerBirthdate fetch html page from nba.com and parse player birthdate
func (p *persistenceService) getPlayerBirthdate(player boxscore2.PlayerDTO) *time.Time {
	playerInfo := p.nbaComClient.PlayerInfoPage(player.PersonId)
	if playerInfo == nil {
		panic("There is no page on nba.com for player: " + player.Name)
	}

	var birthDate time.Time
	playerInfo.Find(".PlayerSummary_playerInfo__om2G4").Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		if children.First().Text() == "BIRTHDATE" {
			node := children.Get(1)
			birthDate, _ = time.Parse("January 2, 2006", node.FirstChild.Data)
		}
	})

	return &birthDate
}

func (p *persistenceService) saveTeamPlayers(teamDto boxscore2.TeamDTO, gameModel games.GameModel, teamModel teams.TeamModel) {
	for _, player := range teamDto.Players {
		playerModel := p.savePlayerModel(player)

		err := p.playersRepository.FirstOrCreateGameStat(players.PlayerGameStats{
			PlayerID:      playerModel.ID,
			GameID:        gameModel.ID,
			TeamID:        teamModel.ID,
			PlayedSeconds: time_utils.FormattedMinutesToSeconds(player.Statistics.Minutes, playedTimeFormat),
			PlsMin:        player.Statistics.Plus - player.Statistics.Minus,
			IsBench:       player.Starter != "1",
		})

		if err != nil {
			panic(err)
		}
	}
}

func (p *persistenceService) saveTeam(dto boxscore2.TeamDTO, leagueId int) teams.TeamModel {
	teamModel, _ := p.teamsRepository.FirstOrCreate(teams.TeamModel{
		Alias:    dto.TeamTricode,
		LeagueID: leagueId,
		Name:     dto.TeamName,
	})

	return teamModel
}

func (p *persistenceService) saveGame(gameDto boxscore2.GameDTO) games.GameModel {
	league := enums.NBA

	// query to get NBA league id
	leagueId := p.getLeagueId()

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

	gameModel, _ := p.gamesRepository.FirstOrCreate(games.GameModel{
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

func (p *persistenceService) getLeagueId() int {
	league, _ := p.leagueRepository.GetLeagueByAliasEn("nba")

	return league.ID
}

func newPersistenceService() *persistenceService {
	return &persistenceService{
		nbaComClient:      nba_com.NewClient(),
		teamsRepository:   teams.NewRepository(),
		gamesRepository:   games.NewRepository(),
		playersRepository: players.NewRepository(),
		leagueRepository:  leagues.NewRepository(),
	}
}
