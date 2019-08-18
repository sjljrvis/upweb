package controller

import (
	"encoding/json"
	"net/http"
// "log"
	"github.com/gorilla/mux"
	Helper "github.com/sjljrvis/deploynow/helpers"
	models "github.com/sjljrvis/deploynow/models"
	."github.com/sjljrvis/deploynow/db" 
)

// GetAll controller
func GetAll(w http.ResponseWriter, r *http.Request) {
	users := []models.User{}
	DB.Select("email ,uuid ,id, user_name").Find(&users)
	Helper.RespondWithJSON(w, 200, users)
}

//Get controller
func Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := models.User{}
	if err := DB.First(&user,params["id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	Helper.RespondWithJSON(w, 200, user)
}

func Create(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	if err := DB.Save(&user).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusCreated, user)
}

func Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := models.User{}
	if err := DB.First(&user,params["id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	} 
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
 
	if err := DB.Save(&user).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusOK, user)
}
 
func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := models.User{}
	if err := DB.First(&user,params["id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}  
	if err := DB.Delete(&user).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusOK, nil)
}
