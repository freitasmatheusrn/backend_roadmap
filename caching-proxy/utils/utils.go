package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

func NormalizeRequest(r *http.Request) string {
	path := strings.ToLower(r.URL.Path)
	query := r.URL.RawQuery 
	return fmt.Sprintf("%s?%s", path, query)
}

func HashString(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
