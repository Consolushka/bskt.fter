package nba_com_api

import (
	"IMP/app/internal/infrastructure/nba_com_api"
	boxscore2 "IMP/app/internal/infrastructure/nba_com_api/dtos/boxscore"
	"IMP/app/internal/infrastructure/nba_com_api/dtos/todays_games"
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/utils/array_utils"
	"encoding/json"
)

const playedTimeFormat = "PT%mM%sS"

type Repository struct {
	nbaComClient       *nba_com_api.Client
	persistenceService *PersistenceService
}

func (n *Repository) TodayGames() (string, []string, error) {
	var scoreboard todays_games.ScoreboardDTO

	scoreBoardJson := n.nbaComClient.TodaysGames()
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
	var gameDto boxscore2.GameDTO

	homeJSON := n.nbaComClient.BoxScore(gameId)
	homeRaw, _ := json.Marshal(homeJSON)

	err := json.Unmarshal(homeRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	n.persistenceService.saveGame(gameDto)

	return gameDto.ToImpModel(), nil
}

func NewRepository() *Repository {
	return &Repository{
		nbaComClient:       nba_com_api.NewNbaComApiClient(),
		persistenceService: NewPersistenceService(),
	}
}
