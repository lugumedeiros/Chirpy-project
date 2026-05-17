package dbman

import (
	"database/sql"
	"os"
	"context"

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

func CreateUser(email string) (database.User, error){
	return config.queries.CreateUser(config.context, email)
}

func ResetUsers()(error){
	return config.queries.ResetUsers(config.context)
}

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