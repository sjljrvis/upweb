package controller

import (
	"encoding/json"
	"net/http"

	// "github.com/gorilla/mux"
	. "github.com/sjljrvis/deploynow/db"
	Helper "github.com/sjljrvis/deploynow/helpers"
	github "github.com/sjljrvis/deploynow/lib/github"
	models "github.com/sjljrvis/deploynow/models"
)

type auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type githubAuth struct {
	Code string `json:"code"`
}

//Login controller
func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var authData auth
	if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user := models.User{}
	if err := DB.First(&user, models.User{Email: authData.Email}).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	hashCheck := Helper.CheckPasswordHash(authData.Password, user.Password)
	if hashCheck {
		token, _ := Helper.GenerateJWT(user.UUID, user.Email)
		Helper.RespondWithJSON(w, 200, map[string]string{"message": "login success", "token": token, "email": user.Email})
	} else {
		Helper.RespondWithError(w, 403, "Fatal Authentication")
		return
	}
}

// Register controller
func Register(w http.ResponseWriter, r *http.Request) {
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

// Register controller
func GithubAuth(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var githubAuthData githubAuth
	if err := json.NewDecoder(r.Body).Decode(&githubAuthData); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	github.AccessToken(githubAuthData.Code, "")
	Helper.RespondWithJSON(w, 200, map[string]string{"message": "login success"})
}
