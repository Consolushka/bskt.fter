package main

import (
	"IMP/app/cmd"
	"IMP/app/database"
	"IMP/app/log"
	"github.com/joho/godotenv"
)

// todo: add tests
func main() {
	log.Init()

	err := godotenv.Load()
	if err != nil {
		panic("Couldn't load env file")
	}

	database.OpenDbConnection()

	cmd.Execute()
}
