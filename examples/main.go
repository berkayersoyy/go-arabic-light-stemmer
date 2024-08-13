package main

import (
	"fmt"
	go_arabic_stemmer "go-arabic-stemmer/arabic/stemmer"
)

func main() {
	arStemmer := go_arabic_stemmer.NewArabicLightStemmer()
	stem := arStemmer.LightStem("أفتضاربانني")
	fmt.Println("Stemmed word:", stem)
}
