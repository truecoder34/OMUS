package controllers

import (
	"OMUS/server/responses"
	"net/http"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To OMUS - One More Lins Shortener")

}
