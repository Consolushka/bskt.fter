package statistics

import (
	"IMP/app/internal/infrastructure/nba_com"
	"IMP/app/internal/infrastructure/translator"
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/players"
	"IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/modules/teams"
	"IMP/app/internal/utils/string_utils"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"time"
)

type Persistence struct {
	leagueRepository  *leagues.Repository
	teamsRepository   *teams.Repository
	gamesRepository   *games.Repository
	playersRepository *players.Repository

	nbaComClient *nba_com.Client

	league enums.League
}

func NewPersistence() *Persistence {
	return &Persistence{
		leagueRepository:  leagues.NewRepository(),
		teamsRepository:   teams.NewRepository(),
		gamesRepository:   games.NewRepository(),
		playersRepository: players.NewRepository(),
		nbaComClient:      nba_com.NewClient(),
	}
}

func (p *Persistence) SaveGameBoxScore(dto *models.GameBoxScoreDTO) error {
	p.league = dto.League

	leagueModel, err := p.leagueRepository.GetLeagueByAliasEn(strings.ToLower(dto.League.String()))
	if err != nil {
		return err
	}

	homeTeamModel, err := p.saveTeamModel(dto.HomeTeam, leagueModel.ID)
	if err != nil {
		return err
	}

	awayTeamModel, err := p.saveTeamModel(dto.AwayTeam, leagueModel.ID)
	if err != nil {
		return err
	}

	gameModel, err := p.saveGameModel(dto, leagueModel.ID, homeTeamModel.ID, awayTeamModel.ID)
	if err != nil {
		return err
	}

	// save players statistics
	err = p.saveTeamStats(dto.HomeTeam, gameModel, homeTeamModel)

	err = p.saveTeamStats(dto.AwayTeam, gameModel, awayTeamModel)

	return nil
}

func (p *Persistence) saveTeamModel(dto models.TeamBoxScoreDTO, leagueId int) (teams.Team, error) {
	teamModel, err := p.teamsRepository.FirstOrCreate(teams.Team{
		Alias:    dto.Alias,
		LeagueID: leagueId,
		Name:     dto.Name,
	})

	return teamModel, err
}

func (p *Persistence) saveGameModel(dto *models.GameBoxScoreDTO, leagueId int, homeTeamId int, awayTeamId int) (games.GameModel, error) {
	gameModel, err := p.gamesRepository.FirstOrCreate(games.GameModel{
		HomeTeamID:    homeTeamId,
		AwayTeamID:    awayTeamId,
		LeagueID:      leagueId,
		ScheduledAt:   dto.ScheduledAt,
		PlayedMinutes: dto.PlayedMinutes,
	})

	return gameModel, err
}

func (p *Persistence) saveTeamStats(dto models.TeamBoxScoreDTO, gameModel games.GameModel, teamModel teams.Team) error {
	teamGameModel, err := p.teamsRepository.FirstOrCreateGameStats(teams.TeamGameStats{
		TeamId: teamModel.ID,
		GameId: gameModel.ID,
		Points: dto.Scored,
	})
	if err != nil {
		return err
	}

	for _, player := range dto.Players {
		playerModel := p.savePlayerModel(player)

		err := p.playersRepository.FirstOrCreateGameStat(players.PlayerGameStats{
			PlayerID:      playerModel.ID,
			TeamGameId:    teamGameModel.Id,
			PlayedSeconds: player.Statistic.PlayedSeconds,
			PlsMin:        player.Statistic.PlsMin,
			IsBench:       player.Statistic.IsBench,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Persistence) savePlayerModel(player models.PlayerDTO) players.Player {
	playerModel, err := p.playersRepository.FirstByLeaguePlayerId(player.LeaguePlayerID)
	if playerModel == nil {
		log.Println("Player not found in database: ", player.LeaguePlayerID, ". Fetching from client")

		if player.BirthDate == nil || player.FullNameLocal == "" {
			playerLocalFullName, playerEnFullName, birthdate := p.getPlayerBio(p.league, player.LeaguePlayerID)
			player.BirthDate = birthdate
			player.FullNameLocal = playerLocalFullName
			player.FullNameEn = playerEnFullName
		}
		playerModel, err = p.playersRepository.FirstOrCreate(players.Player{
			FullNameLocal: player.FullNameLocal,
			FullNameEn:    player.FullNameEn,
			BirthDate:     player.BirthDate,
			OfficialId:    player.LeaguePlayerID,
		})
	}

	if err != nil {
		panic(err)
	}

	return *playerModel
}

func (p *Persistence) getPlayerBio(league enums.League, playerId int) (string, string, *time.Time) {
	if league == enums.NBA {

		playerInfo := p.nbaComClient.PlayerInfoPage(playerId)
		if playerInfo == nil {
			panic("There is no page on nba.com for player id: " + strconv.Itoa(playerId))
		}

		var birthDate time.Time
		var playerLocalFullName string
		playerInfo.Find(".PlayerSummary_playerInfo__om2G4").Each(func(i int, s *goquery.Selection) {
			children := s.Children()
			if children.First().Text() == "BIRTHDATE" {
				node := children.Get(1)
				birthDate, _ = time.Parse("January 2, 2006", node.FirstChild.Data)
			}
		})
		playerLocalFullName = playerInfo.Find(".PlayerSummary_playerNameText___MhqC").Text()
		playerLocalFullName = strings.ReplaceAll(playerLocalFullName, "\n", " ")
		// If player name contains non-latin characters - translates name to EN
		playerEnFullName := playerLocalFullName
		if string_utils.HasNonLanguageChars(playerLocalFullName, string_utils.Latin) {
			playerEnFullName = translator.Translate(playerEnFullName, nil, "en")
		}

		return playerLocalFullName, playerEnFullName, &birthDate
	}
	panic("Unknown league")
}
