package string_utils

import "errors"

type Language int

const (
	Latin Language = iota
	Cyrillic
)

func (l Language) getBoundaries() (min, max int32, err error) {
	switch l {
	case Latin:
		return 0, 127, nil
	case Cyrillic:
		return 1024, 1279, nil // Basic Cyrillic Unicode range
	default:
		return 0, 0, errors.New("invalid language")
	}
}

func HasNonLanguageChars(text string, language Language) (bool, error) {
	minBorder, maxBorder, err := language.getBoundaries()

	if err != nil {
		return false, err
	}

	for _, r := range text {
		if r < minBorder || r > maxBorder {
			return true, nil
		}
	}
	return false, nil
}
