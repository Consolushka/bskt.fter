package dtos

import (
	"FTER/internal/models"
	"fmt"
	"math"
	"strconv"
	"strings"
)

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
		MinutesPlayed: minutesStrToCorrectFormat(p.Statistics.Minutes),
		PlsMin:        p.Statistics.Plus - p.Statistics.Minus,
	}
}

func minutesStrToCorrectFormat(minutesStr string) string {
	timeStr := strings.Trim(minutesStr, "PTS")
	parts := strings.Split(timeStr, "M")
	if len(parts) != 2 {
		return "0:00"
	}

	minutes := parts[0]
	secondsFloat, _ := strconv.ParseFloat(parts[1], 64)
	seconds := fmt.Sprintf("%.0f", math.Floor(secondsFloat))

	// Ensure seconds are padded with leading zero if needed
	if len(seconds) == 1 {
		seconds = "0" + seconds
	}

	return fmt.Sprintf("%s:%s", minutes, seconds)
}
