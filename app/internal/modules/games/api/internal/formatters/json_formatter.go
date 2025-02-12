package formatters

import (
	"IMP/app/internal/modules/games/resources"
	"IMP/app/internal/modules/imp/domain/models"
	"encoding/json"
	"net/http"
)

// jsonFormatter implements ResponseFormatter interface
//
// formats models.GameImpMetrics model to JSON
type jsonFormatter struct{}

func (f *jsonFormatter) Format(w http.ResponseWriter, data *models.GameImpMetrics) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resources.NewMetricResource(*data))
}
