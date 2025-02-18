package main

import (
	"IMP/app/cmd"
	"IMP/app/database"
	"IMP/app/log"
	"github.com/joho/godotenv"
)

func main() {
	log.Init()

	godotenv.Load()

	database.OpenDbConnection()

	cmd.Execute()
}
