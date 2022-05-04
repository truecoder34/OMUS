package controllers

import (
	"OMUS/server/middlewares"
)

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// URL Routes
	s.Router.HandleFunc("/urls", middlewares.SetMiddlewareJSON(s.CreateShortURL)).Methods("POST")
	s.Router.HandleFunc("/urls", middlewares.SetMiddlewareJSON(s.GetURLs)).Methods("GET")
	s.Router.HandleFunc("/urls/{id}", middlewares.SetMiddlewareJSON(s.GetURL)).Methods("GET")
	//s.Router.HandleFunc("/urls/{id}", middlewares.SetMiddlewareJSON(s.UpdateURL)).Methods("PUT")
	s.Router.HandleFunc("/urls/{id}", middlewares.SetMiddlewareJSON(s.DeletePost)).Methods("DELETE")
	s.Router.HandleFunc("/redirect/{encodedURL}", middlewares.SetMiddlewareJSON(s.RedirectByShort)).Methods("GET")

}
