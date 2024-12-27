package dtos

import (
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/utils/time_utils"
)

const playedTimeFormat = "%m:%s"

type PlayerDTO struct {
	FullName        string         `json:"full_name"`
	JerseyNumber    string         `json:"jersey_number"`
	ID              string         `json:"id"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Position        string         `json:"position"`
	PrimaryPosition string         `json:"primary_position"`
	Played          bool           `json:"played"`
	Active          bool           `json:"active"`
	Starter         bool           `json:"starter"`
	OnCourt         bool           `json:"on_court"`
	SrID            string         `json:"sr_id"`
	Reference       string         `json:"reference"`
	Statistics      PlayerStatsDTO `json:"statistics"`
}

func (p *PlayerDTO) ToFterModel() models.PlayerModel {
	return models.PlayerModel{
		FullName:      p.FullName,
		SecondsPlayed: time_utils.FormattedMinutesToSeconds(p.Statistics.Minutes, playedTimeFormat),
		PlsMin:        p.Statistics.PlsMin,
	}

}
