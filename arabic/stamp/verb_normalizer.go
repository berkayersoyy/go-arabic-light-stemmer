package stamp

import (
	"go-arabic-stemmer/arabic/constant"
	"go-arabic-stemmer/arabic/stop_words"
	"regexp"
	"strings"
)

type verbNormalizer struct {
	wordProcessor stop_words.WordProcessor
}

// VerbNormalizer handles the normalization of verbs.
type VerbNormalizer interface {
	Normalize(verb string) string
}

// NewVerbNormalizer creates a new instance of VerbNormalizer with the provided WordProcessor.
// The WordProcessor is used for various word processing tasks, such as stripping Tashkeel.
// This function returns a VerbNormalizer interface that can be used to normalize verbs.
func NewVerbNormalizer(wordProcessor stop_words.WordProcessor) VerbNormalizer {
	return &verbNormalizer{wordProcessor: wordProcessor}
}

// Normalize applies a series of normalization steps to the given verb string.
// It strips Tashkeel, normalizes Hamza characters, removes weak letters, and handles double letters at the end of the verb.
func (vn *verbNormalizer) Normalize(verb string) string {
	if verb == "" {
		return ""
	}

	// Strip Tashkeel from the verb
	verb = vn.wordProcessor.StripTashkeel(verb)

	if verb == "" {
		return ""
	}

	// Normalize 4-letter verbs starting with ALEF_HAMZA_ABOVE
	if len(verb) == 4 && strings.HasPrefix(verb, constant.ALEF_HAMZA_ABOVE) {
		verb = strings.TrimPrefix(verb, constant.ALEF_HAMZA_ABOVE)
	}

	// Normalize Hamza characters in the verb
	verb = vn.normalizeHamza(verb)

	// Remove weak letters from the verb
	verb = vn.removeWeakLetters(verb)

	// Remove double letters at the end of the verb
	verb = vn.removeDoubleLetterAtEnd(verb)

	return verb
}

// normalizeHamza normalizes Hamza characters in the given verb string by replacing them with a standard Hamza ('ء').
func (vn *verbNormalizer) normalizeHamza(verb string) string {
	reHamza := regexp.MustCompile(`[أإءؤئآ]`)
	return reHamza.ReplaceAllString(verb, "ء")
}

// removeWeakLetters removes weak letters ('ا', 'و', 'ي', 'ى') from the given verb string.
func (vn *verbNormalizer) removeWeakLetters(verb string) string {
	reWeakLetters := regexp.MustCompile(`[اويى]`)
	return reWeakLetters.ReplaceAllString(verb, "")
}

// removeDoubleLetterAtEnd removes the last character of the verb if it is the same as the second-to-last character,
// which helps to standardize verbs that end in double letters.
func (vn *verbNormalizer) removeDoubleLetterAtEnd(verb string) string {
	if len(verb) > 1 && verb[len(verb)-1] == verb[len(verb)-2] {
		return verb[:len(verb)-1]
	}
	return verb
}
