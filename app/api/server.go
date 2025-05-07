package api

import (
	"IMP/app/api/middleware"
	players "IMP/app/internal/modules/players/api"
	teams "IMP/app/internal/modules/teams/api"
	"github.com/gorilla/mux"
	"log"
)

// Serve register all application routes and serve HTTP server
func Serve() *mux.Router {
	serverInstance := newServer()

	serverInstance.router.Use(middleware.ContentTypeMiddleware)
	return serverInstance.setupRoutes()
}

type server struct {
	router *mux.Router

	teamsRouter   *teams.Router
	playersRouter *players.Router
}

func newServer() *server {
	return &server{
		router:        mux.NewRouter(),
		teamsRouter:   teams.NewRouter(),
		playersRouter: players.NewRouter(),
	}
}

func (s *server) setupRoutes() *mux.Router {
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	// Register module routes
	s.teamsRouter.RegisterRoutes(apiRouter)
	s.playersRouter.RegisterRoutes(apiRouter)

	s.printRoutes()

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
