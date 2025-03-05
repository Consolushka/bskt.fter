package statistics

import (
	"IMP/app/internal/infrastructure/nba_com"
	"IMP/app/internal/infrastructure/translator"
	gamesDomain "IMP/app/internal/modules/games/domain"
	gamesModels "IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/modules/imp"
	leaguesDomain "IMP/app/internal/modules/leagues/domain"
	leaguesModels "IMP/app/internal/modules/leagues/domain/models"
	playersDomain "IMP/app/internal/modules/players/domain"
	playersModels "IMP/app/internal/modules/players/domain/models"
	statisticModels "IMP/app/internal/modules/statistics/models"
	teamsDomain "IMP/app/internal/modules/teams/domain"
	teamsModels "IMP/app/internal/modules/teams/domain/models"
	"IMP/app/internal/utils/string_utils"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"time"
)

type Persistence struct {
	leagueRepository  *leaguesDomain.Repository
	teamsRepository   *teamsDomain.Repository
	gamesRepository   *gamesDomain.Repository
	playersRepository *playersDomain.Repository

	nbaComClient *nba_com.Client

	league *leaguesModels.League
}

func NewPersistence() *Persistence {
	return &Persistence{
		leagueRepository:  leaguesDomain.NewRepository(),
		teamsRepository:   teamsDomain.NewRepository(),
		gamesRepository:   gamesDomain.NewRepository(),
		playersRepository: playersDomain.NewRepository(),
		nbaComClient:      nba_com.NewClient(),
	}
}

func (p *Persistence) SaveGameBoxScore(dto *statisticModels.GameBoxScoreDTO) error {
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

func (p *Persistence) saveTeamModel(dto statisticModels.TeamBoxScoreDTO) (teamsModels.Team, error) {
	teamModel, err := p.teamsRepository.FirstOrCreate(teamsModels.Team{
		Alias:      dto.Alias,
		LeagueID:   p.league.ID,
		Name:       dto.Name,
		OfficialId: dto.LeagueId,
	})

	return teamModel, err
}

func (p *Persistence) saveGameModel(dto *statisticModels.GameBoxScoreDTO, homeTeamId int, awayTeamId int) (gamesModels.Game, error) {
	gameModel, err := p.gamesRepository.FirstOrCreate(gamesModels.Game{
		HomeTeamID:    homeTeamId,
		AwayTeamID:    awayTeamId,
		LeagueID:      p.league.ID,
		ScheduledAt:   dto.ScheduledAt,
		PlayedMinutes: dto.PlayedMinutes,
		OfficialId:    dto.Id,
	})

	return gameModel, err
}

func (p *Persistence) saveTeamStats(dto statisticModels.TeamBoxScoreDTO, opponents statisticModels.TeamBoxScoreDTO, gameModel gamesModels.Game, teamModel teamsModels.Team) error {
	teamGameModel, err := p.teamsRepository.FirstOrCreateTeamGameStats(teamsModels.TeamGameStats{
		TeamId: teamModel.ID,
		GameId: gameModel.ID,
		Points: dto.Scored,
	})
	if err != nil {
		return err
	}

	for _, player := range dto.Players {
		playerModel := p.savePlayerModel(player)

		err := p.playersRepository.FirstOrCreatePlayerGameStats(playersModels.PlayerGameStats{
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

func (p *Persistence) savePlayerModel(player statisticModels.PlayerDTO) playersModels.Player {
	playerModel, err := p.playersRepository.FirstByOfficialId(player.LeaguePlayerID)
	if playerModel == nil {
		log.Println("Player not found in database: ", player.LeaguePlayerID, ". Fetching from client")

		if player.BirthDate == nil || player.FullNameLocal == "" {
			playerLocalFullName, playerEnFullName, birthdate := p.getPlayerBio(player.LeaguePlayerID)
			player.BirthDate = birthdate
			player.FullNameLocal = playerLocalFullName
			player.FullNameEn = playerEnFullName
		}
		playerModel, err = p.playersRepository.FirstOrCreate(playersModels.Player{
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

func (p *Persistence) getPlayerBio(playerId string) (string, string, *time.Time) {
	if p.league.AliasEn == strings.ToUpper(leaguesModels.NBAAlias) {

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
		if string_utils.HasNonLanguageChars(playerLocalFullName, string_utils.Latin) {
			playerEnFullName = translator.Translate(playerEnFullName, nil, "en")
		}

		return playerLocalFullName, playerEnFullName, &birthDate
	}
	panic("Unknown league")
}
