package main

import (
	"IMP/app/cmd"
	"IMP/app/database"
	"github.com/joho/godotenv"
	"time"
)

// todo: add tests
// todo: maybe create facades for some structs (translator)
func main() {
	loc, _ := time.LoadLocation("Europe/Moscow")
	time.Local = loc

	err := godotenv.Load()
	if err != nil {
		panic("Couldn't load env file")
	}

	database.OpenDbConnection()

	cmd.Execute()
}
