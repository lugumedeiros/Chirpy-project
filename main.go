package main

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/lugumedeiros/Chirpy-project/internal/dbman"
	myserver "github.com/lugumedeiros/Chirpy-project/server"
)

// run build: "go build -o .out && ./.out"
func main() {
	err := godotenv.Load()
	if err!=nil{
		fmt.Printf("ERROR: %v\n", err)
	}

	dbman.DBConnect()
	err = myserver.ListAndServeServer()
	fmt.Printf("%v", err)
}
