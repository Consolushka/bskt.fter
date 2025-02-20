package cdn_nba

import (
	"IMP/app/internal/infrastructure/cdn_nba"
	"IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
	"time"
)

const playedTimeFormat = "PT%mM%sS"

type Provider struct {
	cdnNbaClient *cdn_nba.Client
	mapper       *mapper
}

func (n *Provider) GamesByDate(date time.Time) ([]string, error) {
	schedule := n.cdnNbaClient.ScheduleSeason()

	formattedSearchedDate := date.Format("01/02/2006 00:00:00")

	for _, gameDate := range schedule.Games {
		if gameDate.GameDate == formattedSearchedDate {
			return array_utils.Map(gameDate.Games, func(game cdn_nba.GameSeasonScheduleDto) string {
				return game.GameId
			}), nil
		}
	}

	return make([]string, 0), nil
}

func (n *Provider) GameBoxScore(gameId string) (*models.GameBoxScoreDTO, error) {
	gameDto := n.cdnNbaClient.BoxScore(gameId)

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
