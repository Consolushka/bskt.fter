package main

import (
	"IMP/app/cmd"
	"IMP/app/database"
)

func main() {
	database.OpenDbConnection()
	//cmd.GamePdf("nba", "0022400439")
	//cmd.TodayGamesPdf()
	cmd.Execute()
	//leagues.Seed()
}
