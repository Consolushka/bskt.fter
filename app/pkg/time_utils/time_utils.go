package time_utils

import (
	"strconv"
	"strings"
)

func FormattedMinutesToSeconds(timeStr string, pattern string) int {
	minutesIdx := strings.Index(pattern, "%m")
	secondsIdx := strings.Index(pattern, "%s")

	if minutesIdx == -1 || secondsIdx == -1 || len(timeStr) < len(pattern) {
		return 0
	}

	// Extract minutes
	minutesStart := minutesIdx
	minutesEnd := minutesStart
	for i := minutesStart; i < len(timeStr); i++ {
		if timeStr[i] >= '0' && timeStr[i] <= '9' {
			minutesEnd = i + 1
		} else {
			break
		}
	}

	// Extract seconds
	secondsStart := secondsIdx
	secondsEnd := secondsStart
	for i := secondsStart; i < len(timeStr); i++ {
		if timeStr[i] >= '0' && timeStr[i] <= '9' {
			secondsEnd = i + 1
		} else if timeStr[i] == '.' {
			break
		}
	}

	minutes, _ := strconv.Atoi(timeStr[minutesStart:minutesEnd])
	seconds, _ := strconv.Atoi(timeStr[secondsStart:secondsEnd])

	return minutes*60 + seconds
}
