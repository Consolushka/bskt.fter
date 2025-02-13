package resources

import (
	"IMP/app/internal/modules/leagues/domain/models"
)

type LeagueResource struct {
	Id           int    `json:"id"`
	NameLocal    string `json:"name_local"`
	NameEnglish  string `json:"name_english"`
	AliasEnglish string `json:"alias_english"`
	AliasLocal   string `json:"alias_local"`
}

func NewLeagueResponse(leagueModel *models.League) LeagueResource {
	return LeagueResource{
		Id:           leagueModel.ID,
		NameLocal:    leagueModel.NameLocal,
		NameEnglish:  leagueModel.NameEn,
		AliasLocal:   leagueModel.AliasLocal,
		AliasEnglish: leagueModel.AliasEn,
	}
}
