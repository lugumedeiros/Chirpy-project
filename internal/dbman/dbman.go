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