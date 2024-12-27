package commands

import (
	results2 "IMP/app/internal/modules/imp/results"
	"IMP/app/internal/modules/pdf"
	"IMP/app/internal/modules/pdf/mappers"
	"log"
	"strings"
)

// PrintGame takes boxscore results and prints it to pdf file
// saves in ./outputs/{innerFolder} directory
// innerFolder - optional, if not provided, will be set to date of the game
func PrintGame(game *results2.GameResult, innerFolder *string) {
	if innerFolder == nil {
		innerFolder = &strings.Split(game.Schedule, " ")[0]
	}
	pdfFile := pdf.NewBuilder(game.Title, innerFolder)

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
func toTableMapperSlice(players []results2.PlayerFterResult) []mappers.TableMapper {
	result := make([]mappers.TableMapper, len(players))
	for i, player := range players {
		result[i] = &player
	}
	return result
}
