package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/lugumedeiros/Chirpy-project/internal/dbman"
)

var apicfg apiConfig

func getApiConfig() *apiConfig{
	return &apicfg
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
	fmt.Print("FUNC START: METRIC\n")
	writeTextToServer(w, fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`,
	apicfg.getHits()), "text/html", http.StatusOK)
	fmt.Print("FUNC END: METRIC\n")
}

func resetFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Print("FUNC START: RESET\n")
	if apicfg.getPlatform() != "dev"{
		w.WriteHeader(403)
		return
	}

	apicfg.resetHits()
	writeTextToServer(w, "Metric Reseted", "text/plain", http.StatusOK)
	err := dbman.ResetUsers()
	if err != nil {
		w.WriteHeader(500)
	}
	fmt.Print("FUNC END: RESET\n")
}

func validate_chirpFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Print("FUNC START: VALIDATE\n")
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
	fmt.Print("FUNC END: VALIDATE\n")
}

func setNewUserFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Print("FUNC START: SET USER\n")
	type parameter struct {
		Email string `json:"email"`
	}
	type response struct {
		Id int `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email string `json:"email"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	user, err_db := dbman.CreateUser(params.Email)
	if err_db != nil {
		w.WriteHeader(501)
		w.Write([]byte(err_db.Error()))
		return
	}
	resp := response{
		int(user.ID),
		user.CreatedAt.String(),
		user.UpgradedAt.String(),
		user.Email,
	}
	data, err_marshal := json.Marshal(resp)
	if err_marshal != nil {
		w.WriteHeader(502)
		w.Write([]byte(err_marshal.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
	fmt.Print("FUNC END: SET USER\n")
}