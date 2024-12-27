package translator

import gt "github.com/bas24/googletranslatefree"

// Translate translates the given text to the target language.
// If the source language is not specified, it will be automatically detected.
func Translate(name string, sourceLang *string, targetLang string) string {
	if sourceLang == nil {
		*sourceLang = "auto"
	}

	result, _ := gt.Translate(name, *sourceLang, targetLang)
	return result
}
