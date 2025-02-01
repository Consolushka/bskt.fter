package cdn_nba

import (
	"IMP/app/internal/infrastructure/cdn_nba"
	"IMP/app/internal/infrastructure/cdn_nba/dtos/boxscore"
	"IMP/app/internal/infrastructure/cdn_nba/dtos/schedule_league"
	"IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/log"
	"encoding/json"
	"time"
)

const league = enums.NBA
const playedTimeFormat = "PT%mM%sS"

type Provider struct {
	cdnNbaClient *cdn_nba.Client
	mapper       *mapper
}

func (n *Provider) GamesByDate(date time.Time) ([]string, error) {
	var schedule schedule_league.SeasonScheduleDTO

	scheduleJson := n.cdnNbaClient.ScheduleSeason()
	raw, _ := json.Marshal(scheduleJson)

	err := json.Unmarshal(raw, &schedule)
	if err != nil {
		log.GetLogger().Fatalln(err)
	}

	formattedSearchedDate := date.Format("01/02/2006 00:00:00")

	for _, gameDate := range schedule.Games {
		if gameDate.GameDate == formattedSearchedDate {
			return array_utils.Map(gameDate.Games, func(game schedule_league.GameDTO) string {
				return game.GameId
			}), nil
		}
	}

	return make([]string, 0), err
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

func (n *Provider) GamesByTeam(teamId string) ([]string, error) {
	panic("implement me")
}

func NewProvider() *Provider {
	return &Provider{
		cdnNbaClient: cdn_nba.NewCdnNbaClient(),
		mapper:       newMapper(),
	}
}
