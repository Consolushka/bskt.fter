package resources

import "IMP/app/internal/modules/leagues"

type LeagueResponse struct {
	Id           int    `json:"id"`
	NameLocal    string `json:"name_local"`
	NameEnglish  string `json:"name_english"`
	AliasEnglish string `json:"alias_english"`
	AliasLocal   string `json:"alias_local"`
}

// todo: use domain folder for each module
func NewLeagueResponse(leagueModel *leagues.League) LeagueResponse {
	return LeagueResponse{
		Id:           leagueModel.ID,
		NameLocal:    leagueModel.NameLocal,
		NameEnglish:  leagueModel.NameEn,
		AliasLocal:   leagueModel.AliasLocal,
		AliasEnglish: leagueModel.AliasEn,
	}
}
