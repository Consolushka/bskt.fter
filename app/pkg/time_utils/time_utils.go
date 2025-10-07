package time_utils

import "time"

var MoscowTZ = time.FixedZone("UTC+3", 3*60*60)

func ToMoscowTZ(timeInDiffTZ time.Time) time.Time {
	return timeInDiffTZ.In(MoscowTZ)
}
