package stop_words

import "github.com/berkayersoyy/go-arabic-light-stemmer/arabic/constant"

type TashkeelChecker interface {
	IsTashkeel(char rune) bool
}

// tashkeelChecker handles checking if a character is a Tashkeel.
type tashkeelChecker struct{}

// NewTashkeelChecker creates a new instance of TashkeelChecker.
func NewTashkeelChecker() TashkeelChecker {
	return &tashkeelChecker{}
}

// IsTashkeel returns true if the given character is a Tashkeel, false otherwise.
func (t *tashkeelChecker) IsTashkeel(char rune) bool {
	return constant.TASHKEEL[char]
}
