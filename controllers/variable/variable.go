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
	if err := DB.Where("uuid = ?", params["repository_id"]).First(&repository).Error; err != nil {
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
	ctx := r.Context()
	user := ctx.Value("user").(models.User)
	params := mux.Vars(r)
	activity := models.Activity{}
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
	go activity.Log(DB, repository.ID, user, "configuration" , "Added Configuration variables")
	Helper.RespondWithJSON(w, http.StatusCreated, variable)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// user := ctx.Value("user").(models.User)
	params := mux.Vars(r)
	// activity := models.Activity{}
	variable := models.Variable{}

	query := make(map[string]interface{})
	query["repository_id"] = params["repository_id"]
	query["id"] = params["id"]

	if err := DB.Find(&variable, query).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
	}

	if err := DB.Delete(&variable).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// repo_id, err := strconv.ParseUint(params["repository_id"], 10)
	// go activity.Log(DB, repo_id, user, "configuration" , "Removed Configuration variables")
	Helper.RespondWithJSON(w, http.StatusOK, nil)
}
