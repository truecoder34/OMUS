package controllers

import (
	"OMUS/server/models"
	"OMUS/server/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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

	// check if ENTITY ALREADY EXISTS. IF YES - INCREASE REGENERATES COUNTER
	model, err := url.GetEntityByOriginalURL(server.DB, url.OriginalURL)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		log.Printf(formattedError.Error())
		// responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		// return
	}

	if model.EncodedURL != "" {
		// update regenerations counter

		// TODO-FIX: each time adds AMP to original URL
		urlRes, err := url.UpdateURL(server.DB, 2)
		if err != nil {
			formattedError := formaterror.FormatError(err.Error())
			responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
			return
		}
		log.Printf(strconv.FormatInt(urlRes.RegeneratesCounter, 10))

	} else {
		url.Prepare()
		err = url.Validate()
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		urlRes, err := url.SaveURL(server.DB)
		if err != nil {
			formattedError := formaterror.FormatError(err.Error())
			responses.ERROR(w, http.StatusInternalServerError, formattedError)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, urlRes.ID))

		responses.JSON(w, http.StatusCreated, urlRes)
	}

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

	// Increment visits counter
	updated_url, err := url.UpdateURL(server.DB, 1)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	log.Printf(strconv.FormatInt(updated_url.VisitsCounter, 10))

	// Redirect
	http.Redirect(w, r, model.OriginalURL, http.StatusMovedPermanently)
}

/*
	[ GET ] - Get all URL entities in BASE
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
	[ GET ] get stat by short link
*/
func (server *Server) GetStatistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["encodedURL"]
	url := models.URL{}

	// Check if the URL exist
	urlEntity, err := url.GetEntityByEncodedURL(server.DB, shortLink)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, urlEntity)

}

/*
	[ POST ] - UPDATE URL entity
	TODO
*/

/*
	[ DELETE ] - delete URL entity by ID
	TODO-longterm: make it available only for ADMIN
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
		responses.ERROR(w, http.StatusNotFound, errors.New("url entity with specified ID don't exists"))
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
