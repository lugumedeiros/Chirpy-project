package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lugumedeiros/Chirpy-project/internal/auth"
	"github.com/lugumedeiros/Chirpy-project/internal/dbman"
)

var apicfg apiConfig

func getApiConfig() *apiConfig {
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
	if apicfg.getPlatform() != "dev" {
		w.WriteHeader(403)
		return
	}

	apicfg.resetHits()
	writeTextToServer(w, "Metric Reseted", "text/plain", http.StatusOK)
	err := dbman.ResetUsers()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	err = dbman.DeleteAllChirps()
	if err != nil {
		w.WriteHeader(501)
		return
	}
	fmt.Print("FUNC END: RESET\n")
}

func setNewUserFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Print("FUNC START: SET USER\n")
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Id        int    `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	hash, _ := auth.HashPassword(params.Password)
	user, err_db := dbman.CreateUser(params.Email, hash)
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

func getUserFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Print("FUNC START: GET USER\n")
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Id        int    `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(502)
		return
	}

	user, err_db := dbman.GetUser(params.Email)	
	if err_db != nil {
		w.WriteHeader(500)
		w.Write([]byte(err_db.Error()))
		return
	}
	check, _ := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !check {
		w.WriteHeader(401)
		w.Write([]byte("Incorrect email or password"))
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
	w.WriteHeader(200)
	w.Write(data)
	fmt.Print("FUNC END: GET USER\n")
}

func postChirpFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Print("FUNC START: POST CHIRP\n")
	type input struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}

	type output struct {
		ID        int    `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		UserID    string `json:"user_id"`
		Body      string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	inParams := input{}
	errDec := decoder.Decode(&inParams)
	if errDec != nil {
		w.WriteHeader(500)
		w.Write([]byte(errDec.Error()))
		return
	}
	validatedChirp := validate_chirpFunc(inParams.Body)
	if !validatedChirp.valid {
		w.WriteHeader(501)
		return
	}
	userId, _ := strconv.Atoi(inParams.UserId)
	chirp, errDB := dbman.CreateChirp(userId, inParams.Body)
	if errDB != nil {
		w.WriteHeader(502)
		w.Write([]byte(errDB.Error()))
		return
	}

	var outParams output
	outParams.CreatedAt = chirp.CreatedAt.String()
	outParams.UpdatedAt = chirp.UpdatedAt.String()
	outParams.ID = int(chirp.ID)
	outParams.UserID = fmt.Sprintf("%v", chirp.UserID)
	outParams.Body = chirp.Body

	data, err_marshal := json.Marshal(outParams)
	if err_marshal != nil {
		w.WriteHeader(503)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
	fmt.Print("FUNC END: POST CHIRP\n")
}

func getChirpFunc(w http.ResponseWriter, r *http.Request) {
	type outputItem struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		UserID    string `json:"user_id"`
		Body      string `json:"body"`
	}
	var items []outputItem
	fmt.Print("FUNC START: GET CHIRPS\n")

	chirps, errDB := dbman.GetAllChirps()
	if errDB != nil {
		w.WriteHeader(500)
		return
	}

	for _, chirp := range chirps {
		new_item := outputItem{
			fmt.Sprintf("%v", chirp.ID),
			chirp.CreatedAt.String(),
			chirp.UpdatedAt.String(),
			fmt.Sprintf("%v", chirp.UserID),
			chirp.Body,
		}
		items = append(items, new_item)
	}

	data, err_marshal := json.Marshal(items)
	if err_marshal != nil {
		w.WriteHeader(501)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
	fmt.Print("FUNC END: GET CHIRPS\n")

}

func getChirpByIdFunc(w http.ResponseWriter, r *http.Request) {
	type output struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		UserID    string `json:"user_id"`
		Body      string `json:"body"`
	}
	fmt.Print("FUNC START: GET CHIRP\n")

	id, errconv := strconv.Atoi(r.PathValue("id"))
	if errconv != nil {
		w.WriteHeader(404)
		return
	}
	chirp, err := dbman.GetChirp(id)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	out := output{
		fmt.Sprintf("%v", chirp.ID),
		chirp.CreatedAt.String(),
		chirp.UpdatedAt.String(),
		fmt.Sprintf("%v", chirp.UserID),
		chirp.Body,
	}
	data, err_marshal := json.Marshal(out)
	if err_marshal != nil {
		w.WriteHeader(501)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
	fmt.Print("FUNC END: GET CHIRP\n")
}
