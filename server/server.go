package server

import (
	"fmt"
	"net/http"
	"os"
	"slices"
)

func ListAndServeServer()error{
	setPlatform()
	mux := http.NewServeMux()
	mux.HandleFunc("/app/", rootFunc)
	mux.HandleFunc("GET /admin/healthz", healthzFunc)
	mux.HandleFunc("GET /admin/metrics", metricsFunc)
	mux.HandleFunc("POST /admin/reset", resetFunc)
	// mux.HandleFunc("POST /api/validate_chirp", validate_chirpFunc)
	mux.HandleFunc("POST /api/users", setNewUserFunc)

	mux.HandleFunc("POST /api/chirps", postChirpFunc)
	mux.HandleFunc("GET /api/chirps", getChirpFunc)

	fmt.Printf("FuncsHandler set to Server\n")
	return http.ListenAndServe(":8080", mux)
}

func setPlatform(){
	platforms := []string{"dev", "user", "admin"}
	defaultPlatform := "user"
	platform := os.Getenv("PLATFORM")
	config := getApiConfig()
	if slices.Contains(platforms, platform){
		config.setPlatform(platform)
	} else {
		config.setPlatform(defaultPlatform)
	}
}