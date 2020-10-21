package controllers

import "github.com/norfabagas/auth-global/api/middlewares"

func (s *Server) InitializeRoutes() {
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")
}
