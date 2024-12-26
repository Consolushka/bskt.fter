package boxscore

type TeamName struct {
	CompTeamNameID       int           `json:"CompTeamNameID"`
	TeamID               int           `json:"TeamID"`
	TeamType             int           `json:"TeamType"`
	CompTeamShortNameRu  string        `json:"CompTeamShortNameRu"`
	CompTeamShortNameEn  string        `json:"CompTeamShortNameEn"`
	CompTeamNameRu       string        `json:"CompTeamNameRu"`
	CompTeamNameEn       string        `json:"CompTeamNameEn"`
	CompTeamRegionNameRu string        `json:"CompTeamRegionNameRu"`
	CompTeamRegionNameEn string        `json:"CompTeamRegionNameEn"`
	CompTeamAbcNameRu    string        `json:"CompTeamAbcNameRu"`
	CompTeamAbcNameEn    string        `json:"CompTeamAbcNameEn"`
	CompTeamNameChanged  interface{}   `json:"CompTeamNameChanged"`
	CompTeamNameDefault  bool          `json:"CompTeamNameDefault"`
	SysStatus            int           `json:"SysStatus"`
	SysLastChanged       string        `json:"SysLastChanged"`
	SysUser              interface{}   `json:"SysUser"`
	SysMyUser            interface{}   `json:"SysMyUser"`
	CompTeams            []interface{} `json:"CompTeams"`
	Team                 interface{}   `json:"Team"`
	IsRealTeam           bool          `json:"IsRealTeam"`
}
