package main

import (
	"IMP/app/cmd"
	"IMP/app/database"
	"IMP/app/log"
)

func main() {
	log.Init()

	database.OpenDbConnection()

	cmd.Execute()
}
