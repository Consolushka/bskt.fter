package results

import (
	"FTER/app/internal/modules/imp/enums"
	"FTER/app/internal/modules/imp/models"
	"FTER/app/internal/modules/pdf/mappers"
	"FTER/app/internal/utils/time_utils"
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
	result := []string{t.Player.FullName, time_utils.SecondsToFormat(t.Player.SecondsPlayed, time_utils.PlayedTimeFormat)}

	for _, impValue := range t.ImpPersResults {
		result = append(result, strconv.FormatFloat(impValue, 'f', 2, 64))
	}

	return result

}
