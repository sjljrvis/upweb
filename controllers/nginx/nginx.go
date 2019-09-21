package controller

import (
	"encoding/json"
	"net/http"

	Helper "github.com/sjljrvis/deploynow/helpers"
	// "github.com/gorilla/mux"
	// . "github.com/sjljrvis/deploynow/db"
)

/*
	TO-DO
	Check if request is coming from dnow server /script
	else drop the requests
*/

type nginx struct {
	domain string
	conf   string
}

//Create controller
func Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	payload := nginx{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

}

//Update controller
func Update(w http.ResponseWriter, r *http.Request) {

}
