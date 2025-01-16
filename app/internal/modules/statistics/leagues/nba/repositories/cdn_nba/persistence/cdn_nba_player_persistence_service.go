package persistence

import (
	boxscore2 "IMP/app/internal/infrastructure/cdn_nba/dtos/boxscore"
	"IMP/app/internal/infrastructure/nba_com"
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/players"
	"IMP/app/internal/modules/teams"
	"IMP/app/internal/utils/time_utils"
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

const playedTimeFormat = "PT%mM%sS"

type playerPersistenceService struct {
	nbaComClient *nba_com.Client

	playersRepository *players.Repository
}

// savePlayerGameStats saves player game statistics if not exists by playerId+gameId+teamId
func (p *playerPersistenceService) savePlayerGameStats(player boxscore2.PlayerDTO, gameModel games.GameModel, teamModel teams.Team) {
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

// savePlayerModel saves player if not exists by leaguePlayerId
func (p *playerPersistenceService) savePlayerModel(player boxscore2.PlayerDTO) players.Player {

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
func (p *playerPersistenceService) getPlayerBirthdate(player boxscore2.PlayerDTO) *time.Time {
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

func newPlayerPersistenceService() *playerPersistenceService {
	return &playerPersistenceService{
		nbaComClient:      nba_com.NewClient(),
		playersRepository: players.NewRepository(),
	}
}
