package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func writeTextToServer(w http.ResponseWriter, s string, content string, status int) {
	w.Header().Set("Content-Type", fmt.Sprintf("%v; charset=utf-8", content))
	w.WriteHeader(status)
	w.Write([]byte(s))
}

func unprofaneChirp(s string) (string, bool) {
	profane_words := []string{"kerfuffle", "sharbert", "fornax"}
	replacement := `${1}****${3}`
	regex := regexp.MustCompile(`(?i)(^|\s)(` + strings.Join(profane_words, "|") + `)($|\s)`)
	replaced := regex.ReplaceAllString(s, replacement)
	has_profane := replaced != s
	return replaced, has_profane
}
