package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

var apicfg apiConfig

func (cfg *apiConfig)middleWareMetricInc(next http.Handler) http.Handler{
	cfg.fileserverHits.Add(1)
	return next
}

func (cfg *apiConfig)getHits() int{
	return int(cfg.fileserverHits.Load())
}

func (cfg *apiConfig)resetHits(){
	cfg.fileserverHits.Store(0)
}

func rootFunc(w http.ResponseWriter, r *http.Request) {
	handler := http.FileServer(http.Dir("."))
	handler = apicfg.middleWareMetricInc(handler)
	handler.ServeHTTP(w, r)
}

func healthzFunc(w http.ResponseWriter, r *http.Request) {
	writeStrToServer(w, "OK", "text/plain")
}

func metricsFunc(w http.ResponseWriter, r *http.Request) {
	writeStrToServer(w, fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`,
	apicfg.getHits()), "text/html")
}

func resetFunc(w http.ResponseWriter, r *http.Request) {
	apicfg.resetHits()
	writeStrToServer(w, "Metric Reseted", "text/plain")
}

///

func main(){
	fmt.Printf("Starting Server.")
	mux := http.NewServeMux()
	mux.HandleFunc("/app/", rootFunc)
	mux.HandleFunc("GET /admin/healthz", healthzFunc)
	mux.HandleFunc("GET /admin/metrics", metricsFunc)
	mux.HandleFunc("POST /admin/reset", resetFunc)
	http.ListenAndServe(":8080", mux)
}

///

func writeStrToServer(w http.ResponseWriter, s string, c string){
	w.Header().Set("Content-Type", fmt.Sprintf("%v; charset=utf-8", c))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}
