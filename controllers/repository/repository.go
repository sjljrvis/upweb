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
	ctx := r.Context()
	user := ctx.Value("user").(models.User)
	repositories := []models.Repository{}
	DB.Find(&user).Related(&repositories)
	Helper.RespondWithJSON(w, 200, repositories)
}

//Get controller
func Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(models.User)
	params := mux.Vars(r)
	repository := models.Repository{}
	err := DB.Where("user_id = ?", user.ID).First(&repository, params["id"]).Error
	if err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	Helper.RespondWithJSON(w, 200, repository)
}

func Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(models.User)
	repository := models.Repository{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&repository); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	repository.UserID = user.ID
	repository.UserName = user.UserName
	if err := DB.Save(&repository).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusCreated, repository)
}

func Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repository := models.Repository{}
	if err := DB.First(&repository, params["id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&repository); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := DB.Save(&repository).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusOK, repository)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repository := models.Repository{}
	if err := DB.First(&repository, params["id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if err := DB.Delete(&repository).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusOK, nil)
}
