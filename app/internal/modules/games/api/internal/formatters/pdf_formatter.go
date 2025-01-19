package formatters

import (
	"IMP/app/internal/modules/imp/models"
	"net/http"
	"os"
)

// pdfFormatter implements ResponseFormatter interface
//
// formats models.GameImpMetrics model to pdf file and returns it
type pdfFormatter struct{}

func (f *pdfFormatter) Format(w http.ResponseWriter, data *models.GameImpMetrics) error {
	filePdf, err := os.ReadFile("outputs/2025-01-16/DEN - HOU. 2025-01-16 02:00:00.pdf")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=game-metrics.pdf")
	w.Header().Set("Content-Length", string(len(filePdf)))

	_, err = w.Write(filePdf)
	return err
}
