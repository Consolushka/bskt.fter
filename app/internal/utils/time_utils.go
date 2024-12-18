package utils

import (
	"strconv"
	"strings"
)

func TimeToDecimal(time string) (float64, error) {
	splitted := strings.Split(time, ":")
	minutes, err := strconv.Atoi(splitted[0])
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.Atoi(splitted[1])
	if err != nil {
		return 0, err
	}
	return float64(minutes) + (float64(seconds) / 60), nil

}
