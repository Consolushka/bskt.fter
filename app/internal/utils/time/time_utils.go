package time

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	PlayedTimeFormat = "%d:%02d"
)

func FromFormatedMinutesToSeconds(minutesFormat string, separator string) int {
	splitted := strings.Split(minutesFormat, separator)

	minutes, err := strconv.Atoi(splitted[0])
	if err != nil {
		panic("can't convert " + splitted[0] + " to int")
	}

	seconds, err := strconv.Atoi(splitted[1])
	if err != nil {
		panic("can't convert " + splitted[1] + " to int")
	}
	return minutes*60 + seconds

}

func SecondsToFormat(totalSeconds int, format string) string {
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf(format, minutes, seconds)
}
