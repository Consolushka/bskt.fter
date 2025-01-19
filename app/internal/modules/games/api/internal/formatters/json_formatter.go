package formatters

import (
	"IMP/app/internal/modules/games/api/internal/responses"
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/internal/utils/time_utils"
	"encoding/json"
	"net/http"
	"strconv"
)

// jsonFormatter implements ResponseFormatter interface
//
// formats models.GameImpMetrics model to JSON
type jsonFormatter struct{}

func (f *jsonFormatter) Format(w http.ResponseWriter, data *models.GameImpMetrics) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(fromGameMetricsModel(data))
}

func fromGameMetricsModel(metrics *models.GameImpMetrics) responses.GetSpecificGameMetricsResponse {
	return responses.GetSpecificGameMetricsResponse{
		GameId:       strconv.Itoa(metrics.Id),
		Scheduled:    metrics.Scheduled.Format("02.01.2006 15:04"),
		FullGameTime: time_utils.SecondsToFormat(metrics.FullGameTime, time_utils.PlayedTimeFormat),
		Home:         mapTeamMetricsToResponse(metrics.Home),
		Away:         mapTeamMetricsToResponse(metrics.Away),
	}
}

func mapTeamMetricsToResponse(team models.TeamImpMetrics) responses.GetSpecificGameTeamMetricsResponse {
	return responses.GetSpecificGameTeamMetricsResponse{
		Alias:       team.Alias,
		TotalPoints: team.TotalPoints,
		Players:     mapPlayerMetricsToResponse(team.Players),
	}
}

func mapPlayerMetricsToResponse(players []models.PlayerImpMetrics) []responses.GetSpecificGamePlayerMetricsResponse {
	return array_utils.Map(players, func(player models.PlayerImpMetrics) responses.GetSpecificGamePlayerMetricsResponse {
		return responses.GetSpecificGamePlayerMetricsResponse{
			FullName:      player.FullName,
			MinutesPlayed: time_utils.SecondsToFormat(player.SecondsPlayed, time_utils.PlayedTimeFormat),
			PlsMin:        player.PlsMin,
			IMP:           player.IMP,
		}
	})
}
