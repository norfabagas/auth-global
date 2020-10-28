package controllers

import "github.com/norfabagas/auth-global/api/middlewares"

func (s *Server) InitializeRoutes() {
	// Base route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// /v1 prefix routes
	v1 := s.Router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	v1.HandleFunc("/register", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	v1.HandleFunc("/user", middlewares.SetMiddlewareAuth(middlewares.SetMiddlewareJSON(s.ShowUser))).Methods("GET")
	v1.HandleFunc("/user/edit", middlewares.SetMiddlewareAuth(middlewares.SetMiddlewareJSON(s.UpdateUser))).Methods("PUT")
	v1.HandleFunc("/forget-password", middlewares.SetMiddlewareJSON(s.ForgetPassword)).Methods("POST")
	v1.HandleFunc("/mail", middlewares.SetMiddlewareJSON(s.TestMail)).Methods("GET")
}
