package stop_words

import (
	"strings"
	"unicode"
)

type WordProcessor interface {
	IsVocalized(word string) bool
	IsAlpha(word string) bool
	StripTashkeel(text string) string
}

// wordProcessor handles operations on words.
type wordProcessor struct {
	tashkeelChecker TashkeelChecker
}

// NewWordProcessor creates a new instance of WordProcessor with the provided TashkeelChecker.
func NewWordProcessor(tashkeelChecker TashkeelChecker) WordProcessor {
	return &wordProcessor{tashkeelChecker: tashkeelChecker}
}

// IsVocalized checks if the given word contains any Tashkeel characters.
// It returns true if the word is vocalized (contains Tashkeel), false otherwise.
func (wp *wordProcessor) IsVocalized(word string) bool {
	if wp.IsAlpha(word) {
		return false
	}
	for _, char := range word {
		if wp.tashkeelChecker.IsTashkeel(char) {
			return true
		}
	}
	return false
}

// IsAlpha checks if the given word contains only alphabetic characters.
// It returns true if all characters are alphabetic, false otherwise.
func (wp *wordProcessor) IsAlpha(word string) bool {
	for _, char := range word {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

// StripTashkeel removes all Tashkeel characters from the given text.
// It returns the text without Tashkeel characters, preserving the original order of the remaining characters.
func (wp *wordProcessor) StripTashkeel(text string) string {
	if text == "" {
		return text
	}
	if wp.IsVocalized(text) {
		var result strings.Builder
		for _, char := range text {
			if !wp.tashkeelChecker.IsTashkeel(char) {
				result.WriteRune(char)
			}
		}
		return result.String()
	}
	return text
}
