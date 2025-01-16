package main

import (
	"IMP/app/cmd"
	"IMP/app/database"
)

func main() {
	database.OpenDbConnection()

	cmd.Execute()
}
