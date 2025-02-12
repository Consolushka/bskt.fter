package api

import (
	"IMP/app/internal/handlers"
	"IMP/app/internal/modules/teams/api/internal"
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
	r.HandleFunc("/teams", handlers.BindAndValidateRequestHandler(router.controller.GetTeams)).Methods("GET")

	r.HandleFunc("/teams/{id}", handlers.BindAndValidateRequestHandler(router.controller.GetTeam)).Methods("GET")
	r.HandleFunc("/teams/{id}/games", handlers.BindAndValidateRequestHandler(router.controller.GetTeamGames)).Methods("GET")
	r.HandleFunc("/teams/{id}/games/metrics", handlers.BindAndValidateRequestHandler(router.controller.GetTeamGamesMetrics)).Methods("GET")
}
