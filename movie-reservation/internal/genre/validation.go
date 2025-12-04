package genre

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

var allowedGenres = map[string]bool{
	"ação":    true,
	"drama":   true,
	"terror":  true,
	"comédia": true,
}

func NameValid(name string) error {
	if utf8.RuneCountInString(name) < 4 {
		return fmt.Errorf("gênero precisa ter ao menos 4 letras")
	}
	return nil
}

func SanitizeGenres(input []string) []string {
	var out []string
	for _, g := range input {
		normalized := strings.ToLower(strings.TrimSpace(g))
		if allowedGenres[normalized] {
			out = append(out, normalized)
		}
	}
	return out
}
