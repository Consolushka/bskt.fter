package commands

import (
	"FTER/internal/fter/results"
	"FTER/internal/pdf/mappers"
	"log"
)
import "FTER/internal/pdf"

// PrintGame takes game results and prints it to pdf file
// saves in ./outputs directory
func PrintGame(game *results.GameResult) {
	pdfFile := pdf.NewBuilder(game.Title)

	pdfFile.PrintLn(game.GameId)
	pdfFile.PrintLn(game.Schedule)

	pdfFile.PrintLn(game.Home.Title)
	pdfFile.PrintTable(toTableMapperSlice(game.Home.Players))

	pdfFile.PrintLn(game.Away.Title)
	pdfFile.PrintTable(toTableMapperSlice(game.Away.Players))

	err := pdfFile.Save()
	if err != nil {
		log.Fatal(err)
		return
	}
}

// toTableMapperSlice converts slice of PlayerFterResult to slice of TableMapper
func toTableMapperSlice(players []results.PlayerFterResult) []mappers.TableMapper {
	result := make([]mappers.TableMapper, len(players))
	for i, player := range players {
		result[i] = &player
	}
	return result
}
