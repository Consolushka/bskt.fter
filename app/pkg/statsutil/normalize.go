package statsutil

import "math"

// FromPercentage100 converts a percentage value in range 0-100 to 0.0-1.0.
// It applies guard rails to ensure the result is always within [0.0, 1.0].
func FromPercentage100(val float64) float32 {
	res := float32(val / 100.0)
	if res < 0 {
		return 0
	}
	if res > 1 {
		return 1
	}
	if math.IsNaN(float64(res)) {
		return 0
	}
	return res
}

// FromRatio calculates a percentage from made and attempted counts.
// It returns a value in range 0.0-1.0 and applies guard rails.
func FromRatio(made, attempted int) float32 {
	if attempted <= 0 {
		return 0
	}
	res := float32(made) / float32(attempted)
	if res < 0 {
		return 0
	}
	if res > 1 {
		return 1
	}
	if math.IsNaN(float64(res)) {
		return 0
	}
	return res
}
