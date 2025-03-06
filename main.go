package main

import (
	"IMP/app/cmd"
	"IMP/app/database"
	"IMP/app/log"
	"github.com/joho/godotenv"
)

// todo: custom hostname
// todo: abstract controller
// todo: add tests
// todo: swagger documentation
// todo: elasticsearch to search players by name
func main() {
	log.Init()

	godotenv.Load()

	database.OpenDbConnection()

	cmd.Execute()
}
