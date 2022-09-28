package controllers

import "github.com/LordCeilan/receiptbackend/api/middlewares"

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/home", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Receipt routes
	s.Router.HandleFunc("/receipt", middlewares.SetMiddlewareJSON(s.CreateReceipt)).Methods("POST")
	s.Router.HandleFunc("/receipts", middlewares.SetMiddlewareJSON(s.GetReceipts)).Methods("GET")
	s.Router.HandleFunc("/receipts/{id}", middlewares.SetMiddlewareJSON(s.GetReceipt)).Methods("GET")
	s.Router.HandleFunc("/receipts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateReceipt))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteReceipt)).Methods("DELETE")

	//Create clients controllers/
	//Client routes
	s.Router.HandleFunc("/client", middlewares.SetMiddlewareJSON(s.CreateClient)).Methods("POST")
	s.Router.HandleFunc("/clients", middlewares.SetMiddlewareJSON(s.GetClient)).Methods("GET")
	s.Router.HandleFunc("/clients/{id}", middlewares.SetMiddlewareJSON(s.GetClient)).Methods("GET")
	s.Router.HandleFunc("/clients/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateAClient))).Methods("PUT")
	s.Router.HandleFunc("/client/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteAClient)).Methods("DELETE")

}
