package controllers

import (
	"OMUS/server/models"
	"OMUS/server/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	formaterror "OMUS/server/utils"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

/*
	[ POST ] - create SHORT URL
*/
func (server *Server) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	url := models.URL{}
	err = json.Unmarshal(body, &url)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	url.Prepare()
	err = url.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	urlCreated, err := url.SaveURL(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, urlCreated.ID))

	responses.JSON(w, http.StatusCreated, urlCreated)
}

/*
	[ GET ] - REDIRECT TO SHORT

*/
func (server *Server) RedirectByShort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["encodedURL"]
	url := models.URL{}

	// Check if the URL exist
	model, err := url.GetEntityByEncodedURL(server.DB, shortLink)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	// Redirect
	http.Redirect(w, r, model.OriginalURL, http.StatusMovedPermanently)
}

/*
	[ GET ] - Get all short URLs in BASE
*/
func (server *Server) GetURLs(w http.ResponseWriter, r *http.Request) {

	url := models.URL{}

	urls, err := url.FindAllURLs(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, urls)
}

/*
	[ GET ] - get URL by ID
*/
func (server *Server) GetURL(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := uuid.FromString(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	url := models.URL{}

	urlReceived, err := url.FindURLbyID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, urlReceived)
}

/*
	[ POST ] - UPDATE URL entity
	TODO
*/

/*
	[ DELETE ] - delete URL entity by ID
*/
func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// Is a valid URL id given to us?
	pid, err := uuid.FromString(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if the URL exist
	url := models.URL{}
	err = server.DB.Debug().Model(models.URL{}).Where("id = ?", pid).Take(&url).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Don't exists"))
		return
	}

	_, err = url.DeleteURL(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")

}
