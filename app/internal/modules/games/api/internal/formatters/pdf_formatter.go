package formatters

import (
	"IMP/app/internal/modules/imp/domain/models"
	"IMP/app/internal/modules/pdf"
	"IMP/app/internal/modules/pdf/domain"
	"net/http"
	"os"
	"strconv"
)

// pdfFormatter implements ResponseFormatter interface
//
// formats models.GameImpMetrics model to pdf file and returns it
type pdfFormatter struct{}

func (f *pdfFormatter) Format(w http.ResponseWriter, data *models.GameImpMetrics) error {
	innerFolder := data.Scheduled.Format("2006-01-02")
	fileName := data.Home.Alias + " - " + data.Away.Alias + ". " + data.Scheduled.Format("02.01.2006 15:04")

	pdfFile := pdf.NewBuilder(fileName, &innerFolder)

	pdfFile.PrintLn(strconv.Itoa(data.Id))
	pdfFile.PrintLn(data.Scheduled.Format("02.01.2006 15:04"))

	pdfFile.PrintLn(data.Home.Alias)
	pdfFile.PrintTable(toTableMapperSlice(data.Home.Players))

	pdfFile.PrintLn(data.Away.Alias)
	pdfFile.PrintTable(toTableMapperSlice(data.Away.Players))

	err := pdfFile.Save()
	if err != nil {
		return err
	}

	filePath := "outputs/" + innerFolder + "/" + fileName + ".pdf"
	filePdf, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName+".pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(filePdf)))

	_, err = w.Write(filePdf)
	return err
}

// toTableMapperSlice converts slice of PlayerImpResult to slice of TableMapper
func toTableMapperSlice(players []models.PlayerImpMetrics) []domain.TableMapper {
	result := make([]domain.TableMapper, len(players))
	for i, player := range players {
		result[i] = &player
	}
	return result
}
