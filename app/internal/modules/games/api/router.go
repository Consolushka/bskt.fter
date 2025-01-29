package api

import (
	"IMP/app/internal/handlers"
	"IMP/app/internal/modules/games/api/internal"
	"github.com/gorilla/mux"
)

type Router struct {
	controller *internal.Controller
}

func NewRouter() *Router {
	return &Router{
		controller: internal.NewController(),
	}
}

func (router *Router) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/games", handlers.BindAndValidateRequestHandler(router.controller.GetGames)).Methods("GET")

	r.HandleFunc("/games/{id}", handlers.BindAndValidateRequestHandler(router.controller.GetGame)).Methods("GET")
	r.HandleFunc("/games/{id}/metrics", handlers.BindAndValidateRequestHandler(router.controller.GetGameMetrics)).Methods("GET")
}
