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
	repositories := []models.Repository{}
	DB.Find(&repositories)
	Helper.RespondWithJSON(w, 200, repositories)
}

//Get controller
func Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repository := models.Repository{}
	if err := DB.First(&repository,params["id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	Helper.RespondWithJSON(w, 200, repository)
}

func Create(w http.ResponseWriter, r *http.Request) {
	repository := models.Repository{}
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
	Helper.RespondWithJSON(w, http.StatusCreated, repository)
}

func Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repository := models.Repository{}
	if err := DB.First(&repository,params["id"]).Error; err != nil {
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
	if err := DB.First(&repository,params["id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}  
	if err := DB.Delete(&repository).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusOK, nil)
}

// // Search controller
// func Search(w http.ResponseWriter, r *http.Request) {

// 	var query map[string]interface{}
// 	query = make(map[string]interface{})

// 	keys := r.URL.Query()

// 	for item := range keys {
// 		query[item] = keys[item][0]
// 	}

// 	result, err := RepositoryModel.FindByQuery(query)
// 	if err != nil {
// 		Helper.RespondWithError(w, 200, err.Error())
// 	} else {
// 		Helper.RespondWithJSON(w, 200, result)
// 	}
// }
