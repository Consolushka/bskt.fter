package main

import (
	ftercommands "FTER/internal/fter/commands"
	pdfcommands "FTER/internal/pdf/commands"
	"FTER/internal/statistics/factories"
	"log"
)

func main() {
	gamePdf()
}

// gamePdf takes sportRadar gameId and generates pdf with players FTER
func gamePdf() {
	//gameId := "2aa29340-f4ca-4e43-be10-02a7415eece4"
	gameId := "0022401229"
	repo, err := factories.NewStatsRepository()
	if err != nil {
		log.Fatal(err)
		return
	}
	game, err := repo.GameBoxScore(gameId)

	gameRes := ftercommands.CalculateFullGame(game)
	pdfcommands.PrintGame(gameRes)
}
