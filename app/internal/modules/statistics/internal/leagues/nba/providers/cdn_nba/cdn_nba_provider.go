package cdn_nba

import (
	"IMP/app/internal/infrastructure/cdn_nba"
	"IMP/app/internal/infrastructure/cdn_nba/dtos/boxscore"
	"IMP/app/internal/infrastructure/cdn_nba/dtos/todays_games"
	"IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
	"encoding/json"
)

const league = enums.NBA
const playedTimeFormat = "PT%mM%sS"

type Provider struct {
	cdnNbaClient *cdn_nba.Client
	mapper       *mapper
}

func (n *Provider) TodayGames() (string, []string, error) {
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

func (n *Provider) GameBoxScore(gameId string) (*models.GameBoxScoreDTO, error) {
	var gameDto boxscore.GameDTO

	homeJSON := n.cdnNbaClient.BoxScore(gameId)
	homeRaw, _ := json.Marshal(homeJSON)

	err := json.Unmarshal(homeRaw, &gameDto)
	if err != nil {
		return nil, err
	}

	gameBoxScoreDto := n.mapper.mapGame(gameDto)

	return &gameBoxScoreDto, nil
}

func NewProvider() *Provider {
	return &Provider{
		cdnNbaClient: cdn_nba.NewCdnNbaClient(),
		mapper:       newMapper(),
	}
}
