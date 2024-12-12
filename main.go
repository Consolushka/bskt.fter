package main

import (
	ftercommands "NBATrueEfficency/internal/fter/commands"
	pdfcommands "NBATrueEfficency/internal/pdf/commands"
	"NBATrueEfficency/internal/statistics/factories"
	"log"
)

func main() {
	gamePdf()
}

// gamePdf takes sportRadar gameId and generates pdf with players FTER
func gamePdf() {
	gameId := "2a51502e-cd49-4806-b324-ea98752da37b"

	repo, err := factories.NewStatsRepository()
	if err != nil {
		log.Fatal(err)
		return
	}
	game, err := repo.GetGame(gameId)

	gameRes := ftercommands.CalculateFullGame(game.ToFterModel())
	pdfcommands.PrintGame(gameRes)
}
