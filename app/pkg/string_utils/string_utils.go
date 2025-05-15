package string_utils

import (
	"errors"
	"strings"
	"unicode"
)

type StringUtilsInterface interface {
	HasNonLanguageChars(text string, language Language) (bool, error)
	RemovePunctuationAndSpaces(text string) string
}

func NewStringUtils() StringUtilsInterface {
	return &stringUtils{}
}

type stringUtils struct{}

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

func (s stringUtils) HasNonLanguageChars(text string, language Language) (bool, error) {
	minBorder, maxBorder, err := language.getBoundaries()

	if err != nil {
		return false, err
	}

	trimmedText := s.RemovePunctuationAndSpaces(text)

	for _, r := range trimmedText {
		if r < minBorder || r > maxBorder {
			return true, nil
		}
	}
	return false, nil
}

func (s stringUtils) RemovePunctuationAndSpaces(text string) string {
	var result strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			result.WriteRune(r)
		}
	}

	return result.String()
}
