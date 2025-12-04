package reference

import (
	"math/rand"
	"strings"
	"time"
)

const letterCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digitCharset = "0123456789"

func generateRandomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func Generate() string {
	now := time.Now()
	datePart := now.Format("020106")
	randomLetter := generateRandomString(1, letterCharset)
	randomNumberPart := generateRandomString(4, digitCharset)
	var result strings.Builder
	result.WriteString(datePart)
	result.WriteString(randomLetter)
	result.WriteString(randomNumberPart)

	finalString := result.String()
	return finalString
}
