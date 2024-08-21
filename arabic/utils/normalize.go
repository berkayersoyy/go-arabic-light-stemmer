package utils

import (
	"github.com/berkayersoyy/go-arabic-light-stemmer/arabic/constant"
	"github.com/berkayersoyy/go-arabic-light-stemmer/arabic/regex"
)

func StripTashkeel(text string) string {
	return regex.CreateHarakatPattern().ReplaceAllString(text, "")
}

func StripTatweel(text string) string {
	return regex.CreateTatwaalPattern().ReplaceAllString(text, "")
}

func NormalizeHamza(text string) string {
	text = regex.CreateAlefatPattern().ReplaceAllString(text, constant.ALEF)
	return regex.CreateHamzatPattern().ReplaceAllString(text, "\u0621")
}

func NormalizeLamAlef(text string) string {
	return regex.CreateLamAlefatPattern().ReplaceAllString(text, constant.LAM_ALEF+constant.ALEF)
}

func NormalizeSpellErrors(text string) string {
	text = regex.CreateTehMarbutaPattern().ReplaceAllString(text, constant.HEH)
	return regex.CreateAlefMaksuraPattern().ReplaceAllString(text, constant.YEH)
}

func NormalizeSearchText(text string) string {
	text = StripTashkeel(text)
	text = StripTatweel(text)
	text = NormalizeLamAlef(text)
	text = NormalizeHamza(text)
	text = NormalizeSpellErrors(text)
	return text
}
