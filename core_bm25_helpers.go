package main

import (
	"strings"
	"unicode"
)

// Tokenizar limpia y normaliza el texto en palabras individuales ignorando puntuación.
func Tokenizar(texto string) []string {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	words := strings.FieldsFunc(strings.ToLower(texto), f)
	tokens := make([]string, 0, len(words))
	for _, w := range words {
		if len(w) > 2 {
			tokens = append(tokens, w)
		}
	}
	return tokens
}
