package main

import (
	"IMP/app/cmd"
	"IMP/app/database"
	"IMP/app/log"
	"github.com/joho/godotenv"
)

// todo: save "best" (top-25) players by average IMP with PERS
func main() {
	log.Init()

	godotenv.Load()

	database.OpenDbConnection()

	cmd.Execute()
}
