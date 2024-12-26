package boxscore

import (
	models2 "FTER/app/internal/models"
	"FTER/app/internal/utils/arrays"
)

type TeamBoxscore struct {
	TeamNumber       int      `json:"TeamNumber"`
	TeamID           int      `json:"TeamID"`
	TeamName         TeamName `json:"TeamName"`
	Score            int      `json:"Score"`
	Points           int      `json:"Points"`
	Shot1            int      `json:"Shot1"`
	Goal1            int      `json:"Goal1"`
	Shot2            int      `json:"Shot2"`
	Goal2            int      `json:"Goal2"`
	Shot3            int      `json:"Shot3"`
	Goal3            int      `json:"Goal3"`
	PaintShot        int      `json:"PaintShot"`
	PaintGoal        int      `json:"PaintGoal"`
	Shots1           string   `json:"Shots1"`
	Shot1Percent     string   `json:"Shot1Percent"`
	Shots2           string   `json:"Shots2"`
	Shot2Percent     string   `json:"Shot2Percent"`
	Shots3           string   `json:"Shots3"`
	Shot3Percent     string   `json:"Shot3Percent"`
	PaintShots       string   `json:"PaintShots"`
	PaintShotPercent string   `json:"PaintShotPercent"`
	Assist           int      `json:"Assist"`
	Blocks           int      `json:"Blocks"`
	DefRebound       int      `json:"DefRebound"`
	OffRebound       int      `json:"OffRebound"`
	Rebound          int      `json:"Rebound"`
	Steal            int      `json:"Steal"`
	Turnover         int      `json:"Turnover"`
	TeamDefRebound   int      `json:"TeamDefRebound"`
	TeamOffRebound   *int     `json:"TeamOffRebound"`
	TeamRebound      int      `json:"TeamRebound"`
	TeamSteal        *int     `json:"TeamSteal"`
	TeamTurnover     int      `json:"TeamTurnover"`
	Foul             int      `json:"Foul"`
	OpponentFoul     int      `json:"OpponentFoul"`
	Seconds          int      `json:"Seconds"`
	PlayedTime       string   `json:"PlayedTime"`
	PlusMinus        *int     `json:"PlusMinus"`
	//Coach            map[string]interface{} `json:"Coach"`
	Players []PlayerBoxscore `json:"Players"`
	//Coaches          map[string]interface{} `json:"Coaches"`
}

func (dto TeamBoxscore) ToFterModel() models2.TeamGameResultModel {
	return models2.TeamGameResultModel{
		Team: models2.TeamModel{
			FullName: dto.TeamName.CompTeamNameEn,
			Alias:    dto.TeamName.CompTeamAbcNameEn,
		},
		TotalPoints: dto.Score,
		Players:     gameDtoToPlayers(dto.Players),
	}
}

func gameDtoToPlayers(players []PlayerBoxscore) []models2.PlayerModel {
	players = arrays.Filter(players, func(p PlayerBoxscore) bool {
		return p.Seconds > 0
	})

	playersModels := make([]models2.PlayerModel, len(players))

	for i, player := range players {
		playersModels[i] = player.ToFterModel()
	}

	return playersModels
}
