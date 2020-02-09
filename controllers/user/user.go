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
	users := []models.User{}
	DB.Select("email, uuid, id, user_name, password, md5").Find(&users)
	Helper.RespondWithJSON(w, 200, &users)
}

//Get controller
func GetGithubAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := models.User{}
	github_account := models.GithubAccount{}
	query := make(map[string]interface{})
	query["uuid"] = params["uuid"]
	if err := DB.First(&user, query).Related(&github_account).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	Helper.RespondWithJSON(w, 200, github_account)
}

func RemoveGithubAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := models.User{}
	github_account := models.GithubAccount{}
	query := make(map[string]interface{})
	query["uuid"] = params["uuid"]
	DB.Find(&user, query)
	if err := DB.Delete(&github_account, models.GithubAccount{UserID: user.ID}).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	Helper.RespondWithJSON(w, 200, github_account)
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
	defer r.Body.Close()
	ctx := r.Context()
	user := ctx.Value("user").(models.User)
	userObj := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userObj); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := DB.Model(user).Update(userObj).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusOK, user)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	user := ctx.Value("user").(models.User)
	userObj := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userObj); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	ok := Helper.CheckPasswordHash(userObj["old_password"], user.Password)
	if !ok {
		Helper.RespondWithError(w, http.StatusBadRequest, "Old Password did not match")
		return
	}
	if err := DB.Model(&user).Update("password", userObj["new_password"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusOK, user)

}

func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := models.User{}
	if err := DB.First(&user, params["id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if err := DB.Delete(&user).Error; err != nil {
		Helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	Helper.RespondWithJSON(w, http.StatusOK, nil)
}
