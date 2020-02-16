package controller

import (
	"encoding/json"
	"net/http"

	// "log"
	"github.com/gorilla/mux"
	. "github.com/sjljrvis/deploynow/db"
	github "github.com/sjljrvis/deploynow/lib/github"

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
	query := make(map[string]interface{})
	query["user_id"] = user.ID
	query["uuid"] = params["uuid"]
	err := DB.Where(query).First(&repository).Error
	if err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	Helper.RespondWithJSON(w, 200, repository)
}

func GetGithubRepos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(models.User)
	github_account := models.GithubAccount{}
	err := DB.First(&user).Related(&github_account).Error
	if err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, "Github account not connected")
		return
	}
	result := github.Repositories(github_account.Login)
	Helper.RespondWithJSON(w, 200, result)
}

func Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(models.User)
	activity := models.Activity{}
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
	activity.Log(DB, repository, user, "Initial Release")
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

func LinkGithub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(models.User)
	activity := models.Activity{}
	repository := models.Repository{}
	payload := map[string]string{}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	query := make(map[string]interface{})
	query["user_id"] = user.ID
	query["uuid"] = payload["uuid"]
	if err := DB.Find(&repository, query).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	repository.GithubLinked = true
	repository.GithubURL = payload["clone_url"]

	if err := DB.Save(&repository).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	activity.Log(DB, repository, user, "Linked Github Repository")
	Helper.RespondWithJSON(w, http.StatusCreated, repository)
}
