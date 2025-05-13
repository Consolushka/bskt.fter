package string_utils

// mockStringUtils - мок-реализация StringUtilsInterface для тестирования
type mockStringUtils struct {
	stringUtils                    StringUtilsInterface
	HasNonLanguageCharsFunc        func(text string, language Language) (bool, error)
	RemovePunctuationAndSpacesFunc func(text string) string
}

// HasNonLanguageChars - мок-метод, который вызывает настроенную функцию
func (m *mockStringUtils) HasNonLanguageChars(text string, language Language) (bool, error) {
	if m.HasNonLanguageCharsFunc == nil {
		return m.stringUtils.HasNonLanguageChars(text, language)
	}
	return m.HasNonLanguageCharsFunc(text, language)
}

// RemovePunctuationAndSpaces - мок-метод, который вызывает настроенную функцию
func (m *mockStringUtils) RemovePunctuationAndSpaces(text string) string {
	if m.RemovePunctuationAndSpacesFunc == nil {
		return m.stringUtils.RemovePunctuationAndSpaces(text)
	}
	return m.RemovePunctuationAndSpacesFunc(text)
}

// NewMockStringUtils создает новый экземпляр мока с настройками по умолчанию
func NewMockStringUtils(hasNonLanguageCharsPointer *func(text string, language Language) (bool, error), removePunctuationAndSpacesPointer *func(text string) string) *mockStringUtils {
	def := NewStringUtils()

	var hasNonLanguageChars func(text string, language Language) (bool, error)
	var removePunctuationAndSpaces func(text string) string

	if hasNonLanguageCharsPointer == nil {
		hasNonLanguageChars = def.HasNonLanguageChars
	} else {
		hasNonLanguageChars = *hasNonLanguageCharsPointer
	}
	if removePunctuationAndSpacesPointer == nil {
		removePunctuationAndSpaces = def.RemovePunctuationAndSpaces
	} else {
		removePunctuationAndSpaces = *removePunctuationAndSpacesPointer
	}
	return &mockStringUtils{
		HasNonLanguageCharsFunc:        hasNonLanguageChars,
		RemovePunctuationAndSpacesFunc: removePunctuationAndSpaces,
	}
}
