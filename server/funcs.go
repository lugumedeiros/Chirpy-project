package server

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

// USERS
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

func loginUserFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Print("FUNC START: GET USER\n")
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Id           int    `json:"id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	expireDuration := time.Duration(time.Hour * 24 * 60)
	accessToken, err := auth.MakeJWT(fmt.Sprintf("%v", user.ID), apicfg.getJWTKey(), expireDuration)
	if err != nil {
		w.WriteHeader(502)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	_, errtoken := dbman.CreateRefreshToken(refreshToken, int(user.ID), time.Now().Add(expireDuration))
	if errtoken != nil {
		w.Write([]byte(errtoken.Error()))
		w.WriteHeader(401)
	}

	resp := response{
		int(user.ID),
		user.CreatedAt.String(),
		user.UpgradedAt.String(),
		user.Email,
		accessToken,
		refreshToken,
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

func updateUserFunc(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Email    string `json"email"`
		Password string `json"password"`
	}

	decoder := json.NewDecoder(r.Body)
	inParams := input{}
	errDec := decoder.Decode(&inParams)
	if errDec != nil {
		w.WriteHeader(400)
		w.Write([]byte(errDec.Error()))
		return
	}
	token, _ := auth.GetBearerToken(r.Header)
	newHashedPass, _ := auth.HashPassword(inParams.Password)
	userId, err := auth.ValidateJWT(token, apicfg.getJWTKey())
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("Invalid token" + err.Error()))
		return
	}
	userIdInt, _ := strconv.Atoi(userId)
	errupdt := dbman.UpdateUser(userIdInt, inParams.Email, newHashedPass)
	if errupdt != nil {
		w.WriteHeader(401)
		w.Write([]byte(errDec.Error()))
		return
	}

	type out struct {
		Email string `json:"email"`
	}
	data, err_marshal := json.Marshal(out{inParams.Email})
	if err_marshal != nil {
		w.WriteHeader(503)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

// CHIRP
func postChirpFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Print("FUNC START: POST CHIRP\n")
	type input struct {
		Body string `json:"body"`
		// UserId string `json:"user_id"`
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

	tokenString, _ := auth.GetBearerToken(r.Header)
	userId, err := auth.ValidateJWT(tokenString, apicfg.getJWTKey())
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("Invalid token" + err.Error()))
		return
	}
	userIdStr, erratoi := strconv.Atoi(userId)
	if erratoi != nil {
		w.WriteHeader(500)
		return
	}

	chirp, errDB := dbman.CreateChirp(userIdStr, inParams.Body)
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

func deleteChirpByIdFunc(w http.ResponseWriter, r *http.Request){
	tokenString, _ := auth.GetBearerToken(r.Header)
	userIdstring, err := auth.ValidateJWT(tokenString, apicfg.getJWTKey())
	userId, _ := strconv.Atoi(userIdstring)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("Invalid token" + err.Error()))
		return
	}
	chirpId, errconv := strconv.Atoi(r.PathValue("id"))
	if errconv != nil {
		w.WriteHeader(401)
		return
	}
	chirp, err := dbman.GetChirp(chirpId)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	if chirp.UserID != int32(userId){
		w.WriteHeader(403)
		return
	}

	err = dbman.DeleteChirp(chirpId)
	if err != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(204)
	}
}

// TOKENS
func refreshTokenFunc(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	tokenKey, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
		return
	}
	token, errdb := dbman.GetToken(tokenKey)
	if errdb != nil {
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
		return
	}
	if token.RevokedAt.Valid || token.ExpiresAt.Compare(time.Now()) == -1 {
		w.WriteHeader(401)
		return
	}

	jwtToken, errtoken := auth.MakeJWT(
		fmt.Sprintf("%v", token.UserID),
		apicfg.getJWTKey(),
		time.Duration(time.Hour),
	)
	if errtoken != nil {
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
		return
	}

	resp := response{jwtToken}
	data, err_marshal := json.Marshal(resp)
	if err_marshal != nil {
		w.WriteHeader(502)
		w.Write([]byte(err_marshal.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

func revokeTokenFunc(w http.ResponseWriter, r *http.Request) {
	tokenKey, _ := auth.GetBearerToken(r.Header)
	err := dbman.RevokeToken(tokenKey)
	if err != nil {
		w.WriteHeader(506)
		return
	}
	w.WriteHeader(204)
}
