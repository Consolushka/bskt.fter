package dtos

type TeamStatsDTO struct {
	Alias   string      `json:"alias"`
	Points  int         `json:"points"`
	Players []PlayerDTO `json:"players"`
}
