package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/lugumedeiros/Chirpy-project/internal"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

var apicfg apiConfig

func (cfg *apiConfig) middleWareMetricInc(next http.Handler) http.Handler {
	cfg.fileserverHits.Add(1)
	return next
}

func (cfg *apiConfig) getHits() int {
	return int(cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetHits() {
	cfg.fileserverHits.Store(0)
}

func rootFunc(w http.ResponseWriter, r *http.Request) {
	handler := http.FileServer(http.Dir("."))
	handler = apicfg.middleWareMetricInc(handler)
	handler.ServeHTTP(w, r)
}

func healthzFunc(w http.ResponseWriter, r *http.Request) {
	writeTextToServer(w, "OK", "text/plain", http.StatusOK)
}

func metricsFunc(w http.ResponseWriter, r *http.Request) {
	writeTextToServer(w, fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`,
		apicfg.getHits()), "text/html", http.StatusOK)
}

func resetFunc(w http.ResponseWriter, r *http.Request) {
	apicfg.resetHits()
	writeTextToServer(w, "Metric Reseted", "text/plain", http.StatusOK)
}

func validate_chirpFunc(w http.ResponseWriter, r *http.Request) {
	const chirpy_size = 140
	type parameters struct {
		Body string `json:"body"`
	}
	type returnParameters struct {
		Valid        bool   `json:"valid"`
		Error        string `json:"error"`
		Cleaned_body string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	status := 400
	returnParams := returnParameters{false, "Something went wrong", ""}
	var data []byte
	if err == nil {
		if len(params.Body) < chirpy_size {
			clean_chirpy, _ := unprofaneChirp(params.Body)
			returnParams.Error = "None"
			returnParams.Valid = true
			returnParams.Cleaned_body = clean_chirpy
			status = 200
		} else {
			returnParams.Error = "Chirp is too long"
		}

		data, err = json.Marshal(returnParams)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

///

func main() {
	//Initalization
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}
	internal.DBConnect()


	fmt.Printf("Starting Server.")
	mux := http.NewServeMux()
	mux.HandleFunc("/app/", rootFunc)
	mux.HandleFunc("GET /admin/healthz", healthzFunc)
	mux.HandleFunc("GET /admin/metrics", metricsFunc)
	mux.HandleFunc("POST /admin/reset", resetFunc)
	mux.HandleFunc("POST /api/validate_chirp", validate_chirpFunc)
	http.ListenAndServe(":8080", mux)
}

///

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
