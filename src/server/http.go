package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atanurdemir/gatekeeper/src/handler"
	"github.com/atanurdemir/gatekeeper/src/middlewares"
	"github.com/gorilla/mux"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start(port string) {
	router := mux.NewRouter()

	// Global Middlewares
	router.Use(middlewares.GuardMiddleware)

	// Public routes
	s.setupPublicRoutes(router)

	// Private (authenticated) routes
	privateRouter := router.PathPrefix("/").Subrouter()
	privateRouter.Use(middlewares.AuthMiddleware)
	s.setupPrivateRoutes(privateRouter)

	fmt.Printf("Server is starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) setupPublicRoutes(router *mux.Router) {
	router.HandleFunc("/restriction", handler.RestrictionHandler).Methods(http.MethodGet)
}

func (s *Server) setupPrivateRoutes(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(handler.RootHandler).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch)
}
