package nba_com_api

import (
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/players"
	"IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/statistics/leagues/nba/repositories/nba.com_api/client"
	"IMP/app/internal/modules/statistics/leagues/nba/repositories/nba.com_api/dtos/boxscore"
	todays_games2 "IMP/app/internal/modules/statistics/leagues/nba/repositories/nba.com_api/dtos/todays_games"
	"IMP/app/internal/modules/teams"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/internal/utils/time_utils"
	"encoding/json"
	"time"
)

const playedTimeFormat = "PT%mM%sS"

type Repository struct {
	client *client.NbaComApiClient
}

func (n *Repository) TodayGames() (string, []string, error) {
	var scoreboard todays_games2.ScoreboardDTO

	scoreBoardJson := n.client.TodaysGames()
	raw, _ := json.Marshal(scoreBoardJson)

	err := json.Unmarshal(raw, &scoreboard)

	if err != nil {
		return "", nil, err
	}

	return scoreboard.GameDate, array_utils.Map(scoreboard.Games, func(game todays_games2.GameDTO) string {
		return game.GameId
	}), nil
}

func (n *Repository) GameBoxScore(gameId string) (*models.GameModel, error) {
	var gameDto boxscore.GameDTO

	homeJSON := n.client.BoxScore(gameId)
	homeRaw, _ := json.Marshal(homeJSON)

	err := json.Unmarshal(homeRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	saveGame(gameDto)

	return gameDto.ToImpModel(), nil
}

func getNbaLeagueId() int {
	league, _ := leagues.NewRepository().LeagueByAliasEn("nba")

	return league.ID
}

func saveTeam(dto boxscore.TeamDTO, leagueId int) teams.TeamModel {
	teamModel, err := teams.FirstOrCreate(teams.TeamModel{
		Alias:    dto.TeamTricode,
		LeagueID: leagueId,
		Name:     dto.TeamName,
	})

	if err != nil {
		panic(err)
	}

	return teamModel
}

func saveGame(gameDto boxscore.GameDTO) games.GameModel {
	league := enums.NBA

	leagueId := getNbaLeagueId()

	homeTeamModel := saveTeam(gameDto.HomeTeam, leagueId)

	awayTeamModel := saveTeam(gameDto.AwayTeam, leagueId)

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

	for _, player := range gameDto.HomeTeam.Players {
		playerModel, _ := players.FirstOrCreate(players.Player{FullName: player.Name, BirthDate: time.Now()})

		players.CreateStatisticInGame(players.PlayerGameStats{
			PlayerID:      playerModel.ID,
			GameID:        gameModel.ID,
			TeamID:        homeTeamModel.ID,
			PlayedSeconds: time_utils.FormattedMinutesToSeconds(player.Statistics.Minutes, playedTimeFormat),
			PlsMin:        player.Statistics.Plus - player.Statistics.Minus,
			IsBench:       player.Starter != "1",
		})
	}

	for _, player := range gameDto.AwayTeam.Players {
		playerModel, _ := players.FirstOrCreate(players.Player{FullName: player.Name, BirthDate: time.Now()})

		players.CreateStatisticInGame(players.PlayerGameStats{
			PlayerID:      playerModel.ID,
			GameID:        gameModel.ID,
			TeamID:        awayTeamModel.ID,
			PlayedSeconds: time_utils.FormattedMinutesToSeconds(player.Statistics.Minutes, playedTimeFormat),
			PlsMin:        player.Statistics.Plus - player.Statistics.Minus,
			IsBench:       player.Starter != "1",
		})
	}

	panic(gameModel.ID)
	return gameModel
}

func NewRepository() *Repository {
	return &Repository{
		client: client.NewNbaComApiClient(),
	}
}
