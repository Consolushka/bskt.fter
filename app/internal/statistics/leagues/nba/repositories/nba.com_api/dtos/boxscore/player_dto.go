package boxscore

import (
	"FTER/app/internal/models"
	"FTER/app/internal/utils/time"
)

const playedTimeFormat = "PTS%mM%sS"

type PlayerDTO struct {
	Status     string              `json:"status"`
	Order      int                 `json:"order"`
	PersonId   int                 `json:"personId"`
	JerseyNum  string              `json:"jerseyNum"`
	Position   string              `json:"position"`
	Starter    string              `json:"starter"`
	Oncourt    string              `json:"oncourt"`
	Played     string              `json:"played"`
	Statistics PlayerEfficiencyDTO `json:"statistics"`
	Name       string              `json:"name"`
	NameI      string              `json:"nameI"`
	FirstName  string              `json:"firstName"`
	FamilyName string              `json:"familyName"`
}

func (p *PlayerDTO) ToFterModel() models.PlayerModel {
	return models.PlayerModel{
		FullName:      p.Name,
		SecondsPlayed: time.FormattedMinutesToSeconds(p.Statistics.Minutes, playedTimeFormat),
		PlsMin:        p.Statistics.Plus - p.Statistics.Minus,
	}
}
