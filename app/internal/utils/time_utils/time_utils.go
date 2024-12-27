package time_utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	PlayedTimeFormat = "%d:%02d"
)

func FormattedMinutesToSeconds(timeStr string, pattern string) int {
	// Заменяем спецсимволы на regex паттерны
	regexPattern := strings.NewReplacer(
		"%m", `(\d+)`,
		"%s", `(\d+)`,
	).Replace(pattern)

	// Экранируем остальные символы в паттерне
	regexPattern = regexp.QuoteMeta(regexPattern)

	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(timeStr)

	if len(matches) != 3 {
		return 0
	}

	minutes, _ := strconv.Atoi(matches[1])
	seconds, _ := strconv.Atoi(matches[2])

	return minutes*60 + seconds
}

func SecondsToFormat(totalSeconds int, format string) string {
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf(format, minutes, seconds)
}
