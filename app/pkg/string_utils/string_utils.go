package string_utils

type Language int

const (
	Latin Language = iota
	Cyrillic
)

func (l Language) getBoundaries() (min, max int32) {
	switch l {
	case Latin:
		return 0, 127
	case Cyrillic:
		return 1024, 1279 // Basic Cyrillic Unicode range
	default:
		return 0, 127
	}
}

func HasNonLanguageChars(text string, language Language) bool {
	minBorder, maxBorder := language.getBoundaries()

	for _, r := range text {
		if r < minBorder || r > maxBorder {
			return true
		}
	}
	return false
}
