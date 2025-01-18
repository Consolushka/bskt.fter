package games

import (
	"IMP/app/internal/handlers"
	"github.com/gorilla/mux"
)

type Router struct {
	controller *controller
}

func NewRouter() *Router {
	return &Router{
		controller: newController(),
	}
}

func (router *Router) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/game/{id}", handlers.BindAndValidateRequestHandler(router.controller.getGame)).Methods("GET")
	r.HandleFunc("/game/{id}/metrics", handlers.BindAndValidateRequestHandler(router.controller.getGameMetrics)).Methods("GET")
}
