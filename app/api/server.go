package api

import (
	games "IMP/app/internal/modules/games/api"
	teams "IMP/app/internal/modules/teams/api"
	"github.com/gorilla/mux"
	"log"
)

// Serve register all application routes and serve HTTP server
func Serve() *mux.Router {
	serverInstance := newServer()
	return serverInstance.setupRoutes()
}

type server struct {
	router *mux.Router

	gamesRouter *games.Router
	teamsRouter *teams.Router
}

func newServer() *server {
	return &server{
		router:      mux.NewRouter(),
		gamesRouter: games.NewRouter(),
		teamsRouter: teams.NewRouter(),
	}
}

func (s *server) setupRoutes() *mux.Router {
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	// Register module routes
	s.gamesRouter.RegisterRoutes(apiRouter)
	s.teamsRouter.RegisterRoutes(apiRouter)

	return s.router
}

// printRoutes prints all registered routes
func (s *server) printRoutes() {
	s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Route: %s [%v]\n", path, methods)
		return nil
	})
}
