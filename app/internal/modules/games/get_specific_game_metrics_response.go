package games

import (
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/internal/utils/time_utils"
	"strconv"
)

type getSpecificGameMetricsResponse struct {
	GameId       string                             `json:"game_id"`
	Scheduled    string                             `json:"scheduled"`
	FullGameTime string                             `json:"full_game_time"`
	Home         getSpecificGameTeamMetricsResponse `json:"home"`
	Away         getSpecificGameTeamMetricsResponse `json:"away"`
}

type getSpecificGameTeamMetricsResponse struct {
	Alias       string                                 `json:"alias"`
	TotalPoints int                                    `json:"total_points"`
	Players     []getSpecificGamePlayerMetricsResponse `json:"players"`
}

type getSpecificGamePlayerMetricsResponse struct {
	FullName      string  `json:"full_name"`
	MinutesPlayed string  `json:"minutes_played"`
	PlsMin        int     `json:"pls_min"`
	IMP           float64 `json:"imp"`
}

func fromGameMetricsModel(metrics *models.GameImpMetrics) getSpecificGameMetricsResponse {
	return getSpecificGameMetricsResponse{
		GameId:       strconv.Itoa(metrics.Id),
		Scheduled:    metrics.Scheduled.Format("02.01.2006 15:04"),
		FullGameTime: time_utils.SecondsToFormat(metrics.FullGameTime, time_utils.PlayedTimeFormat),
		Home:         mapTeamMetricsToResponse(metrics.Home),
		Away:         mapTeamMetricsToResponse(metrics.Away),
	}
}

func mapTeamMetricsToResponse(team models.TeamImpMetrics) getSpecificGameTeamMetricsResponse {
	return getSpecificGameTeamMetricsResponse{
		Alias:       team.Alias,
		TotalPoints: team.TotalPoints,
		Players:     mapPlayerMetricsToResponse(team.Players),
	}
}

func mapPlayerMetricsToResponse(players []models.PlayerImpMetrics) []getSpecificGamePlayerMetricsResponse {
	return array_utils.Map(players, func(player models.PlayerImpMetrics) getSpecificGamePlayerMetricsResponse {
		return getSpecificGamePlayerMetricsResponse{
			FullName:      player.FullName,
			MinutesPlayed: time_utils.SecondsToFormat(player.SecondsPlayed, time_utils.PlayedTimeFormat),
			PlsMin:        player.PlsMin,
			IMP:           player.IMP,
		}
	})
}
