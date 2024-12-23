package repository

import (
	"FTER/app/internal/models"
	"FTER/app/internal/statistics/repositories/nba/client"
	"FTER/app/internal/statistics/repositories/nba/dtos/boxscore"
	"FTER/app/internal/statistics/repositories/nba/dtos/todays_games"
	"FTER/app/internal/utils/arrays"
	"encoding/json"
)

type NbaRepository struct {
	client *client.NbaClient
}

func (n *NbaRepository) TodayGames() (string, []string, error) {
	var scoreboard todays_games.ScoreboardDTO

	scoreBoardJson := n.client.TodaysGames()
	raw, _ := json.Marshal(scoreBoardJson)

	err := json.Unmarshal(raw, &scoreboard)

	if err != nil {
		return "", nil, err
	}

	return scoreboard.GameDate, arrays.Map(scoreboard.Games, func(game todays_games.GameDTO) string {
		return game.GameId
	}), nil
}

func (n *NbaRepository) GameBoxScore(gameId string) (*models.GameModel, error) {
	var gameDto boxscore.GameDTO

	homeJSON := n.client.BoxScore(gameId)
	homeRaw, _ := json.Marshal(homeJSON)

	err := json.Unmarshal(homeRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	return gameDto.ToFterModel(), nil
}

func NewNbaRepository() *NbaRepository {
	return &NbaRepository{
		client: client.NewNbaClient(),
	}
}
