package time

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	PlayedTimeFormat = "%d:%02d"
)

func FromFormatToDecimal(time string) (float64, error) {
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

func SecondsToFormat(totalSeconds int, format string) string {
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf(format, minutes, seconds)
}
