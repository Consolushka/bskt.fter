package results

import (
	"FTER/app/internal/fter/enums"
	"FTER/app/internal/models"
	"FTER/app/internal/pdf/mappers"
	"strconv"
)

type PlayerFterResult struct {
	mappers.TableMapper
	Player         models.PlayerModel
	TimeBases      []enums.TimeBasedImpCoefficient
	ImpPersResults []float64
}

// Headers returns headers for table
func (t *PlayerFterResult) Headers() []string {
	res := []string{"Player", "Minutes Played"}

	for _, timeBase := range t.TimeBases {
		res = append(res, "Imp"+timeBase.Title())
	}

	return res
}

func (t *PlayerFterResult) ToTable() []string {
	result := []string{t.Player.FullName, t.Player.MinutesPlayed}

	for _, impValue := range t.ImpPersResults {
		result = append(result, strconv.FormatFloat(impValue, 'f', 2, 64))
	}

	return result

}
