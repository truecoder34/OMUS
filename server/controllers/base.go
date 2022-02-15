package controllers

/*
	Here will be initializztion of DB cinnection information
	routes initialization
	server start up
*/

import (
	// "fmt"
	// "log"
	// "net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}
