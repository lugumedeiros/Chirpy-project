package dbman

import (
	"database/sql"
	"os"
	"context"
	"time"
	"github.com/lugumedeiros/Chirpy-project/internal/database"
)

type dbconfig struct {
	db      *sql.DB
	queries *database.Queries
	context context.Context
}

var config dbconfig

func DBConnect() error {
	dbUrl := os.Getenv("DB_URL")
	db, errDbConnect := sql.Open("postgres", dbUrl)
	if errDbConnect != nil {
		return errDbConnect
	}
	dbQueries := database.New(db)

	config.db = db
	config.queries = dbQueries
	config.context = context.Background()
	return nil
}

// USER
func CreateUser(email string, hash string) (database.User, error){
	params := database.CreateUserParams{Email: email, HashedPassword: hash}
	return config.queries.CreateUser(config.context, params)
}

func ResetUsers()(error){
	return config.queries.ResetUsers(config.context)
}

func GetUser(email string) (database.User, error){
	return config.queries.GetUser(config.context, email)
}

// CHIRP
func CreateChirp(userId int, body string) (database.Chirp, error){
	var params database.CreateChirpParams
	params.Body = body
	params.UserID = int32(userId)
	return config.queries.CreateChirp(config.context, params)
}

func DeleteChirp(chirpId int) error{
	return config.queries.DeleteChirp(config.context, int32(chirpId))
}

func DeleteChirpByUserId(userId int) error {
	return config.queries.DeleteChirpsByUserId(config.context, int32(userId))
}

func DeleteAllChirps() error {
	return config.queries.DeleteAllChirps(config.context)
}

func GetAllChirps() ([]database.Chirp, error) {
	return config.queries.GetAllChirps(config.context)
}

func GetChirp(chirpID int) (database.Chirp, error) {
	return config.queries.GetChirp(config.context, int32(chirpID))
}

// TOKENS
func CreateRefreshToken(token string, user_id int, expires time.Time)(database.RefreshToken, error) {
	params := database.AddRefreshTokenParams{
		Token: token,
		UserID: int32(user_id),
		ExpiresAt: expires,
		RevokedAt: sql.NullTime{},
	}
	return config.queries.AddRefreshToken(config.context, params)
}

func GetToken(token string)(database.RefreshToken, error){
	return config.queries.GetRefreshTokenToken(config.context, token)
}

func RevokeToken(token string) error{
	param := database.RevokeTokenParams{Token: token, RevokedAt: sql.NullTime{Time: time.Now(), Valid: true}}
	return config.queries.RevokeToken(config.context, param)
}