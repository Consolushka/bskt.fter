package cdn_nba

import (
	"IMP/app/internal/modules/statistics/models"
	cdn_nba2 "IMP/app/internal/statistics/cdn_nba"
	"IMP/app/log"
	"IMP/app/pkg/array_utils"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const playedTimeFormat = "PT%mM%sS"

type Provider struct {
	cdnNbaClient *cdn_nba2.Client
	mapper       *mapper
}

func (n *Provider) GamesByDate(date time.Time) ([]string, error) {
	schedule := n.cachedSeasonSchedule()

	formattedSearchedDate := date.Format("01/02/2006 00:00:00")

	for _, gameDate := range schedule.Games {
		if gameDate.GameDate == formattedSearchedDate {
			return array_utils.Map(gameDate.Games, func(game cdn_nba2.GameSeasonScheduleDto) string {
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

// cachedSeasonSchedule returns the cached season schedule from a file if the file exists and is not older than 7 days
//
// Otherwise it makes a request to the CDN NBA API
func (n *Provider) cachedSeasonSchedule() cdn_nba2.SeasonScheduleDto {
	cacheFilePath := filepath.Join(os.TempDir(), "nba_schedule_cache.json")

	// Checks if cached file exists and is not older than 7 days
	if info, err := os.Stat(cacheFilePath); err == nil {
		if time.Since(info.ModTime()) < 7*time.Hour*24 {
			data, err := os.ReadFile(cacheFilePath)
			if err == nil {
				var schedule cdn_nba2.SeasonScheduleDto
				if json.Unmarshal(data, &schedule) == nil {
					return schedule
				}
			}
		}
	}

	log.Info("There is no cached file or it is outdated, making a request...")

	// Making request to get the schedule
	schedule := n.cdnNbaClient.ScheduleSeason()

	data, err := json.Marshal(schedule)
	if err == nil {
		log.Info("Saving schedule to cache...")
		// Even if there is an error, we still return the schedule from response
		_ = os.WriteFile(cacheFilePath, data, 0644)
	}

	return schedule
}

func NewProvider() *Provider {
	return &Provider{
		cdnNbaClient: cdn_nba2.NewCdnNbaClient(),
		mapper:       newMapper(),
	}
}
