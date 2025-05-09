package time_utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// FormattedMinutesToSeconds converts a time string in the given pattern to seconds.
// The pattern should contain %m for minutes and %s for seconds.
func FormattedMinutesToSeconds(timeStr string, pattern string) (int, error) {
	// Check if pattern contains required placeholders
	minutesIdx := strings.Index(pattern, "%m")
	secondsIdx := strings.Index(pattern, "%s")

	if minutesIdx == -1 || secondsIdx == -1 {
		return 0, errors.New("pattern must contain both %m and %s")
	}

	// Create regex patterns to extract minutes and seconds
	patternRegex := strings.ReplaceAll(pattern, "%m", "(\\d+)")
	patternRegex = strings.ReplaceAll(patternRegex, "%s", "(\\d+)")
	patternRegex = "^" + regexp.QuoteMeta(patternRegex) + "$"
	patternRegex = strings.ReplaceAll(patternRegex, "\\(\\\\d\\+\\)", "(\\d+)")

	re := regexp.MustCompile(patternRegex)
	matches := re.FindStringSubmatch(timeStr)

	if len(matches) != 3 {
		return 0, errors.New("time string does not match the pattern")
	}

	// Determine which group is minutes and which is seconds
	var minutesStr, secondsStr string
	if minutesIdx < secondsIdx {
		minutesStr = matches[1]
		secondsStr = matches[2]
	} else {
		minutesStr = matches[2]
		secondsStr = matches[1]
	}

	minutes, _ := strconv.Atoi(minutesStr)
	seconds, _ := strconv.Atoi(secondsStr)

	return minutes*60 + seconds, nil
}
