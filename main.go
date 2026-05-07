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
	writeToServer(w, "OK")
}

func metricsFunc(w http.ResponseWriter, r *http.Request) {
	writeToServer(w, fmt.Sprintf("Hits: %v",apicfg.getHits()))
}

func resetFunc(w http.ResponseWriter, r *http.Request) {
	apicfg.resetHits()
	writeToServer(w, "Metric Reseted")
}

///

func main(){
	fmt.Printf("Starting Server.")
	mux := http.NewServeMux()
	mux.HandleFunc("/app/", rootFunc)
	mux.HandleFunc("GET /api/healthz", healthzFunc)
	mux.HandleFunc("GET /api/metrics", metricsFunc)
	mux.HandleFunc("POST /api/reset", resetFunc)
	http.ListenAndServe(":8080", mux)
}

///

func writeToServer(w http.ResponseWriter, s string){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}