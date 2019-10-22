package controller

import (
	"encoding/json"
	"net/http"

	// "log"
	"github.com/gorilla/mux"
	. "github.com/sjljrvis/deploynow/db"
	Helper "github.com/sjljrvis/deploynow/helpers"
	models "github.com/sjljrvis/deploynow/models"
)

// GetAll controller
func GetAll(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repository := models.Repository{}
	if err := DB.First(&repository, params["repository_id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	variables := []models.Variable{}
	DB.Find(&repository).Related(&variables)
	Helper.RespondWithJSON(w, 200, variables)
}

//Get controller
func Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repository := models.Repository{}
	if err := DB.First(&repository, params["repository_id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	variable := models.Variable{}
	err := DB.Where("repository_id = ?", params["repository_id"]).First(&variable, params["id"]).Error
	if err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	Helper.RespondWithJSON(w, 200, variable)
}

func Create(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repository := models.Repository{}
	if err := DB.First(&repository, params["repository_id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	variable := models.Variable{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&variable); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	variable.RepositoryID = repository.ID
	if err := DB.Save(&variable).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusCreated, variable)
}
