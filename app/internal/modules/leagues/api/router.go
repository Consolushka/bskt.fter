package api

import (
	"IMP/app/internal/handlers"
	"IMP/app/internal/modules/leagues/api/internal"
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
	r.HandleFunc("/leagues", handlers.BindAndValidateRequestHandler(router.controller.GetLeagues)).Methods("GET")

	r.HandleFunc("/leagues/{id}/games", handlers.BindAndValidateRequestHandler(router.controller.GetGamesByLeagueAndDate)).Methods("GET")
}
