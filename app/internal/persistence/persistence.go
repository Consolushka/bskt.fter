package persistence

import (
	"IMP/app/internal/domain"
	"IMP/app/internal/imp"
	"IMP/app/internal/statistics"
	"IMP/app/internal/statistics/nba_com"
	"IMP/app/internal/statistics/translator"
	"IMP/app/pkg/string_utils"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"time"
)

type Service struct {
	leagueRepository  *domain.LeaguesRepository
	teamsRepository   *domain.TeamsRepository
	gamesRepository   *domain.GamesRepository
	playersRepository *domain.PlayersRepository

	nbaComClient *nba_com.Client

	league *domain.League

	stringUtils string_utils.StringUtilsInterface
	translator  translator.Interface
}

func NewService() *Service {
	return &Service{
		leagueRepository:  domain.NewLeaguesRepository(),
		teamsRepository:   domain.NewTeamsRepository(),
		gamesRepository:   domain.NewGamesRepository(),
		playersRepository: domain.NewPlayersRepository(),
		nbaComClient:      nba_com.NewClient(),
		stringUtils:       string_utils.NewStringUtils(),
		translator:        translator.NewTranslator(),
	}
}

func (p *Service) SaveGameBoxScore(dto *statistics.GameBoxScoreDTO) error {
	var err error
	p.league, err = p.leagueRepository.FirstByAliasEn(strings.ToLower(dto.LeagueAliasEn))
	if err != nil {
		return err
	}

	homeTeamModel, err := p.saveTeamModel(dto.HomeTeam)
	if err != nil {
		return err
	}

	awayTeamModel, err := p.saveTeamModel(dto.AwayTeam)
	if err != nil {
		return err
	}

	gameModel, err := p.saveGameModel(dto, homeTeamModel.ID, awayTeamModel.ID)
	if err != nil {
		return err
	}

	// save players statistics
	err = p.saveTeamStats(dto.HomeTeam, dto.AwayTeam, gameModel, homeTeamModel)

	err = p.saveTeamStats(dto.AwayTeam, dto.HomeTeam, gameModel, awayTeamModel)

	return nil
}

func (p *Service) saveTeamModel(dto statistics.TeamBoxScoreDTO) (domain.Team, error) {
	teamModel, err := p.teamsRepository.FirstOrCreate(domain.Team{
		Alias:      dto.Alias,
		LeagueID:   p.league.ID,
		Name:       dto.Name,
		OfficialId: dto.LeagueId,
	})

	return teamModel, err
}

func (p *Service) saveGameModel(dto *statistics.GameBoxScoreDTO, homeTeamId int, awayTeamId int) (domain.Game, error) {
	gameModel, err := p.gamesRepository.FirstOrCreate(domain.Game{
		HomeTeamID:    homeTeamId,
		AwayTeamID:    awayTeamId,
		LeagueID:      p.league.ID,
		ScheduledAt:   dto.ScheduledAt,
		PlayedMinutes: dto.PlayedMinutes,
		OfficialId:    dto.Id,
	})

	return gameModel, err
}

func (p *Service) saveTeamStats(dto statistics.TeamBoxScoreDTO, opponents statistics.TeamBoxScoreDTO, gameModel domain.Game, teamModel domain.Team) error {
	teamGameModel, err := p.teamsRepository.FirstOrCreateTeamGameStats(domain.TeamGameStats{
		TeamId: teamModel.ID,
		GameId: gameModel.ID,
		Points: dto.Scored,
	})
	if err != nil {
		return err
	}

	for _, player := range dto.Players {
		playerModel := p.savePlayerModel(player)

		err := p.playersRepository.FirstOrCreatePlayerGameStats(domain.PlayerGameStats{
			PlayerID:      playerModel.ID,
			TeamGameId:    teamGameModel.Id,
			PlayedSeconds: player.Statistic.PlayedSeconds,
			PlsMin:        player.Statistic.PlsMin,
			IsBench:       player.Statistic.IsBench,
			IMPClean:      imp.EvaluateClean(player.Statistic.PlayedSeconds, player.Statistic.PlsMin, dto.Scored-opponents.Scored, gameModel.PlayedMinutes),
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Service) savePlayerModel(player statistics.PlayerDTO) domain.Player {
	playerModel, err := p.playersRepository.FirstByOfficialId(player.LeaguePlayerID)
	if playerModel == nil {
		log.Println("Player not found in database: ", player.LeaguePlayerID, ". Fetching from client")

		if player.BirthDate == nil || player.FullNameLocal == "" {
			playerLocalFullName, playerEnFullName, birthdate := p.getPlayerBio(player.LeaguePlayerID)
			player.BirthDate = birthdate
			player.FullNameLocal = playerLocalFullName
			player.FullNameEn = playerEnFullName
		}
		playerModel, err = p.playersRepository.FirstOrCreate(domain.Player{
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

func (p *Service) getPlayerBio(playerId string) (string, string, *time.Time) {
	if p.league.AliasEn == strings.ToUpper(domain.NBAAlias) {

		playerInfo := p.nbaComClient.PlayerInfoPage(playerId)
		if playerInfo == nil {
			panic("There is no page on nba.com for player id: " + playerId)
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
		hasNonLatinChars, err := p.stringUtils.HasNonLanguageChars(playerLocalFullName, string_utils.Latin)
		if hasNonLatinChars || err != nil {
			playerEnFullName = p.translator.Translate(playerEnFullName, nil, "en")
		}

		return playerLocalFullName, playerEnFullName, &birthDate
	}
	panic("Unknown league")
}
