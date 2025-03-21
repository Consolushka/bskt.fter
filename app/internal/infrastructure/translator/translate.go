package translator

import (
	gt "github.com/bas24/googletranslatefree"
)

// Translate translates the given text to the target language.
// If the source language is not specified, it will be automatically detected.
func Translate(source string, sourceLangPoint *string, targetLang string) string {
	var sourceLang string

	if sourceLangPoint == nil {
		sourceLang = "auto"
	} else {
		sourceLang = *sourceLangPoint
	}

	result, _ := gt.Translate(source, sourceLang, targetLang)
	return result
}
