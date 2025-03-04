package api

import (
	"IMP/app/internal/handlers"
	"IMP/app/internal/modules/players/api/internal"
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
	r.HandleFunc("/players/search", handlers.BindAndValidateRequestHandler(router.controller.Search)).Methods("GET")

	r.HandleFunc("/players/{id}/games", handlers.BindAndValidateRequestHandler(router.controller.PlayerGamesBoxScore)).Methods("GET")
	r.HandleFunc("/players/{id}/games/metrics", handlers.BindAndValidateRequestHandler(router.controller.PlayerGamesMetrics)).Methods("GET")
}
