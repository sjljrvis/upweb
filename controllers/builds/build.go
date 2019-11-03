package controller

import (
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
	builds := []models.Build{}
	DB.Find(&repository).Related(&builds)
	Helper.RespondWithJSON(w, 200, builds)
}
