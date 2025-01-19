package formatters

import (
	"IMP/app/internal/modules/imp/models"
	"net/http"
)

// ResponseFormatter is an interface that defines the method to format the response data
// Can format models.GameImpMetrics model to json or pdf file
type ResponseFormatter interface {
	Format(w http.ResponseWriter, data *models.GameImpMetrics) error
}

func NewFormatter(format string) ResponseFormatter {
	switch format {
	case "pdf":
		return &pdfFormatter{}
	default:
		return &jsonFormatter{}
	}
}
