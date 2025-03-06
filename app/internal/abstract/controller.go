package abstract

import (
	"IMP/app/log"
	"encoding/json"
	"net/http"
)

type BaseController struct {
}

func (b *BaseController) Ok(w http.ResponseWriter, response interface{}) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		b.InternalServerError(w, err)
	}
}

func (b *BaseController) InternalServerError(w http.ResponseWriter, err error) {
	log.Error(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
