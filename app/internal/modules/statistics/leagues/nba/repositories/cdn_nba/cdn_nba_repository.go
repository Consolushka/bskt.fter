package cdn_nba

import (
	"IMP/app/internal/infrastructure/cdn_nba"
	boxscore2 "IMP/app/internal/infrastructure/cdn_nba/dtos/boxscore"
	"IMP/app/internal/infrastructure/cdn_nba/dtos/todays_games"
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/utils/array_utils"
	"encoding/json"
)

const playedTimeFormat = "PT%mM%sS"

type Repository struct {
	cdnNbaClient       *cdn_nba.Client
	persistenceService *persistenceService
}

func (n *Repository) TodayGames() (string, []string, error) {
	var scoreboard todays_games.ScoreboardDTO

	scoreBoardJson := n.cdnNbaClient.TodaysGames()
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

	homeJSON := n.cdnNbaClient.BoxScore(gameId)
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
		cdnNbaClient:       cdn_nba.NewCdnNbaClient(),
		persistenceService: newPersistenceService(),
	}
}
