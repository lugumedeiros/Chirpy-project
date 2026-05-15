package internal

import (
	"database/sql"
	"os"

	"github.com/lugumedeiros/Chirpy-project/internal/database"
)

type dbconfig struct {
	db      *sql.DB
	queries *database.Queries
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
	return nil
}
