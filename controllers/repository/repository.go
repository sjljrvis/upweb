package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	Helper "github.com/sjljrvis/deploynow/helpers"
	RepositoryModel "github.com/sjljrvis/deploynow/models/repository"
	"gopkg.in/mgo.v2/bson"
)

// GetAll controller
func GetAll(w http.ResponseWriter, r *http.Request) {
	var results []RepositoryModel.Repository
	results, err := RepositoryModel.FindAll("5d18ec259c3c131ca4cfe759")
	if err != nil {
		Helper.RespondWithError(w, 200, err.Error())
	} else {
		Helper.RespondWithJSON(w, 200, results)
	}
}

// GetOne controller
func GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result, err := RepositoryModel.FindOneByID("123", params["id"])
	if err != nil {
		Helper.RespondWithError(w, 200, err.Error())
	} else {
		Helper.RespondWithJSON(w, 200, result)
	}
}

// Create controller
func Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var repository RepositoryModel.Repository

	if err := json.NewDecoder(r.Body).Decode(&repository); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	repository.ID = bson.NewObjectId()
	err := RepositoryModel.Create(repository)
	if err != nil {
		Helper.RespondWithError(w, 200, err.Error())
	} else {
		Helper.RespondWithJSON(w, 200, map[string]string{"message": "repository Created successfully"})
	}
}

// Search controller
func Search(w http.ResponseWriter, r *http.Request) {

	var query map[string]interface{}
	query = make(map[string]interface{})

	keys := r.URL.Query()

	for item := range keys {
		query[item] = keys[item][0]
	}

	result, err := RepositoryModel.FindByQuery(query)
	if err != nil {
		Helper.RespondWithError(w, 200, err.Error())
	} else {
		Helper.RespondWithJSON(w, 200, result)
	}
}
