package time_utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type TimeUtilsInterface interface {
	FormattedMinutesToSeconds(timeStr string, pattern string) (int, error)
}

type timeUtils struct{}

func NewTimeUtils() TimeUtilsInterface {
	return &timeUtils{}
}

// FormattedMinutesToSeconds converts a time string in the given pattern to seconds.
// The pattern should contain %m for minutes and %s for seconds.
func (u *timeUtils) FormattedMinutesToSeconds(timeStr string, pattern string) (int, error) {
	// Check if pattern contains required placeholders
	minutesIdx := strings.Index(pattern, "%m")
	secondsIdx := strings.Index(pattern, "%s")

	if minutesIdx == -1 || secondsIdx == -1 {
		return 0, errors.New("pattern must contain both %m and %s")
	}

	// Create regex patterns to extract minutes and seconds
	patternRegex := strings.ReplaceAll(pattern, "%m", `(\d+)`)
	patternRegex = strings.ReplaceAll(patternRegex, "%s", `(\d+(\.\d+)?)`)
	patternRegex = `^` + patternRegex + "$"

	re := regexp.MustCompile(patternRegex)
	matches := re.FindStringSubmatch(timeStr)

	if len(matches) != 4 {
		return 0, errors.New("time string '" + timeStr + "' does not match the pattern '" + pattern + "'")
	}

	// Determine which group is minutes and which is seconds
	var minutesStr, secondsStr string
	if minutesIdx < secondsIdx {
		minutesStr = matches[1]
		secondsStr = matches[2]
	} else {
		minutesStr = matches[3]
		secondsStr = matches[1]
	}

	minutes, _ := strconv.Atoi(minutesStr)
	seconds, _ := strconv.ParseFloat(secondsStr, 32)

	return minutes*60 + int(seconds), nil
}
