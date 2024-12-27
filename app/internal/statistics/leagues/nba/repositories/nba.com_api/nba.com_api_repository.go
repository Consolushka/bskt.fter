package nba_com_api

import (
	"FTER/app/internal/models"
	"FTER/app/internal/statistics/leagues/nba/repositories/nba.com_api/client"
	"FTER/app/internal/statistics/leagues/nba/repositories/nba.com_api/dtos/boxscore"
	"FTER/app/internal/statistics/leagues/nba/repositories/nba.com_api/dtos/todays_games"
	"FTER/app/internal/utils/array_utils"
	"encoding/json"
)

type Repository struct {
	client *client.NbaComApiClient
}

func (n *Repository) TodayGames() (string, []string, error) {
	var scoreboard todays_games.ScoreboardDTO

	scoreBoardJson := n.client.TodaysGames()
	raw, _ := json.Marshal(scoreBoardJson)

	err := json.Unmarshal(raw, &scoreboard)

	if err != nil {
		return "", nil, err
	}

	return scoreboard.GameDate, array_utils.Map(scoreboard.Games, func(game todays_games.GameDTO) string {
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

	return gameDto.ToFterModel(), nil
}

func NewRepository() *Repository {
	return &Repository{
		client: client.NewNbaComApiClient(),
	}
}
