package responses

import (
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/players"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/internal/utils/time_utils"
	"strconv"
)

type GetSpecificGameResponse struct {
	GameId    string                             `json:"game_id"`
	Scheduled string                             `json:"scheduled"`
	Home      getSpecificGameTeamDetailsResponse `json:"home"`
	Away      getSpecificGameTeamDetailsResponse `json:"away"`
}

type getSpecificGameTeamDetailsResponse struct {
	FullName string                                 `json:"full_name"`
	Alias    string                                 `json:"alias"`
	Score    int                                    `json:"score"`
	Players  []getSpecificGamePlayerDetailsResponse `json:"players"`
}

type getSpecificGamePlayerDetailsResponse struct {
	FullName      string `json:"full_name"`
	MinutesPlayed string `json:"minutes_played"`
	Start         bool   `json:"start"`
	PlsMin        int    `json:"pls_min"`
}

func NewGetSpecificGameResponse(gameModel *games.GameModel) GetSpecificGameResponse {
	return GetSpecificGameResponse{
		GameId:    strconv.Itoa(gameModel.ID),
		Scheduled: gameModel.ScheduledAt.Format("02.01.2006 15:04"),
		Home: getSpecificGameTeamDetailsResponse{
			FullName: gameModel.HomeTeamStats.Team.Name,
			Alias:    gameModel.HomeTeamStats.Team.Alias,
			Score:    gameModel.HomeTeamStats.Points,
			Players:  mapPlayersToResponse(gameModel.HomeTeamStats.PlayerGameStats),
		},
		Away: getSpecificGameTeamDetailsResponse{
			FullName: gameModel.AwayTeamStats.Team.Name,
			Alias:    gameModel.AwayTeamStats.Team.Alias,
			Score:    gameModel.AwayTeamStats.Points,
			Players:  mapPlayersToResponse(gameModel.AwayTeamStats.PlayerGameStats),
		},
	}
}

func mapPlayersToResponse(playersList []players.PlayerGameStats) []getSpecificGamePlayerDetailsResponse {
	return array_utils.Map(playersList, func(playerGameStats players.PlayerGameStats) getSpecificGamePlayerDetailsResponse {
		formattedPlayedTime := time_utils.SecondsToFormat(playerGameStats.PlayedSeconds, time_utils.PlayedTimeFormat)
		return getSpecificGamePlayerDetailsResponse{
			FullName:      playerGameStats.Player.FullName,
			MinutesPlayed: formattedPlayedTime,
			Start:         !playerGameStats.IsBench,
			PlsMin:        playerGameStats.PlsMin,
		}
	})

}
