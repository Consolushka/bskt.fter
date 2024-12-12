package commands

import (
	"NBATrueEfficency/internal/fter/results"
	"NBATrueEfficency/internal/pdf/mappers"
)
import "NBATrueEfficency/internal/pdf"

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

	pdfFile.Save()
}

// toTableMapperSlice converts slice of PlayerFterResult to slice of TableMapper
func toTableMapperSlice(players []results.PlayerFterResult) []mappers.TableMapper {
	result := make([]mappers.TableMapper, len(players))
	for i, player := range players {
		result[i] = &player
	}
	return result
}
