package regex

import (
	"go-arabic-stemmer/arabic/constant"
	"regexp"
	"strings"
)

// CreatePattern generates a regular expression pattern from a list of characters.
func CreatePattern(chars ...string) *regexp.Regexp {
	return regexp.MustCompile("[" + strings.Join(chars, "") + "]")
}

func CreateHarakatPattern() *regexp.Regexp {
	return CreatePattern(
		constant.FATHATAN,
		constant.DAMMATAN,
		constant.KASRATAN,
		constant.FATHA,
		constant.DAMMA,
		constant.KASRA,
		constant.SUKUN,
		constant.SHADDA,
	)
}

func CreateHamzatPattern() *regexp.Regexp {
	return CreatePattern(
		constant.WAW_HAMZA,
		constant.YEH_HAMZA,
	)
}

func CreateAlefatPattern() *regexp.Regexp {
	return CreatePattern(
		constant.ALEF_MADDA,
		constant.ALEF_HAMZA_ABOVE,
		constant.ALEF_HAMZA_BELOW,
		constant.HAMZA_ABOVE,
		constant.HAMZA_BELOW,
	)
}

func CreateLamAlefatPattern() *regexp.Regexp {
	return CreatePattern(
		constant.LAM_ALEF,
		constant.LAM_ALEF_HAMZA_ABOVE,
		constant.LAM_ALEF_HAMZA_BELOW,
		constant.LAM_ALEF_MADDA_ABOVE,
	)
}

func CreateTatwaalPattern() *regexp.Regexp {
	return CreatePattern(constant.TATWEEL)
}

func CreateTehMarbutaPattern() *regexp.Regexp {
	return CreatePattern(constant.TEH_MARBUTA)
}

func CreateAlefMaksuraPattern() *regexp.Regexp {
	return CreatePattern(constant.ALEF_MAKSURA)
}
