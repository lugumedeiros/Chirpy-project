package server

import (
	"fmt"
	"net/http"
	"os"
	"slices"
)

func ListAndServeServer() error {
	setPlatform()
	mux := http.NewServeMux()
	mux.HandleFunc("/app/", rootFunc)
	mux.HandleFunc("GET /admin/healthz", healthzFunc)
	mux.HandleFunc("GET /admin/metrics", metricsFunc)
	mux.HandleFunc("POST /admin/reset", resetFunc)
	mux.HandleFunc("POST /api/login", loginUserFunc)
	mux.HandleFunc("POST /api/users", setNewUserFunc)
	mux.HandleFunc("PUT /api/users", updateUserFunc)

	mux.HandleFunc("POST /api/chirps", postChirpFunc)
	mux.HandleFunc("GET /api/chirps", getChirpFunc)
	mux.HandleFunc("GET /api/chirps/{id}", getChirpByIdFunc)
	mux.HandleFunc("DELETE /api/chirps/{id}", deleteChirpByIdFunc)

	mux.HandleFunc("POST /api/refresh", refreshTokenFunc)
	mux.HandleFunc("POST /api/revoke", revokeTokenFunc)

	mux.HandleFunc("POST /api/polka/webhooks", upgradeUserFunc)

	fmt.Printf("FuncsHandler set to Server\n")
	return http.ListenAndServe(":8080", mux)
}

func setPlatform() {
	platforms := []string{"dev", "user", "admin"}
	defaultPlatform := "user"
	platform := os.Getenv("PLATFORM")
	config := getApiConfig()
	if slices.Contains(platforms, platform) {
		config.setPlatform(platform)
	} else {
		config.setPlatform(defaultPlatform)
	}
}
