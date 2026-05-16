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

type validatedChirp struct {
	old string
	new string
	profane bool
	valid bool
}

func validate_chirpFunc(s string) (v validatedChirp) {
	const chirpy_size = 140
	v.old = s
	if len(s) > chirpy_size{
		v.valid = false
		return v
	}
	v.valid = true
	v.new, v.profane = unprofaneChirp(s)
	return v
}
