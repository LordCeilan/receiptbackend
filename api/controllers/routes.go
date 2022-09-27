package controllers

import "github.com/LordCeilan/receiptbackend/api/middlewares"

func (s *Server) initializeRoutes() {
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")
}
