package main

import (
	"fmt"
	"github.com/berkayersoyy/go-arabic-light-stemmer/arabic/stemmer"
)

func main() {
	arStemmer := stemmer.NewArabicLightStemmer()
	stem := arStemmer.LightStem("أفتضاربانني")
	fmt.Println("Stemmed word:", stem)
}
