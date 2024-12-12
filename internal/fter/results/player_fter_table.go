package results

import (
	"FTER/internal/fter/models"
	"FTER/internal/pdf/mappers"
	"strconv"
)

type PlayerFterResult struct {
	mappers.TableMapper
	Player models.PlayerModel
	FTER   float64
}

// Headers returns headers for table
func (t *PlayerFterResult) Headers() []string {
	return []string{"Player", "Minutes Played", "FTER"}
}

func (t *PlayerFterResult) ToTable() []string {
	return []string{
		t.Player.FullName,
		t.Player.MinutesPlayed,
		strconv.FormatFloat(t.FTER, 'f', 2, 64),
	}
}
