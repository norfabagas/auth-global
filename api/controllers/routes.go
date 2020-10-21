package controllers

import "github.com/norfabagas/auth-global/api/middlewares"

func (s *Server) InitializeRoutes() {
	// Base route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// /v1 prefix routes
	v1 := s.Router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
}
