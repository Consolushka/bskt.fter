package models

type PlayerImpRanking struct {
	Rank             int     `json:"rank" db:"rank"`
	ID               int     `json:"id" db:"id"`
	FullNameLocal    string  `json:"full_name_local" db:"full_name_local"`
	AvgImpClean      float64 `json:"avg_imp_clean" db:"avg_imp_clean"`
	AvgPlayedSeconds float64 `json:"avg_played_seconds" db:"avg_played_seconds"`
	GamesPlayed      int     `json:"games_played" db:"games_played"`
}
