package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/phayes/freeport"
	Helper "github.com/sjljrvis/deploynow/helpers"
	container "github.com/sjljrvis/deploynow/lib/container"
	fs "github.com/sjljrvis/deploynow/lib/fs"

	"github.com/sjljrvis/deploynow/log"
	models "github.com/sjljrvis/deploynow/models"

	// "github.com/gorilla/mux"
	. "github.com/sjljrvis/deploynow/db"
)

type build struct {
	BuildPack      string `json:"build_pack"`
	RepositoryName string `json:"repository_name"`
	CommitHash     string `json:"commit_hash"`
}

//Build controller
func Build(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var buildData build
	buildlogChannel := make(chan []byte)
	repository := models.Repository{}
	notify := r.Context().Done()

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Error().Msgf("Streaming unsupported! %s", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	go func() {
		<-notify
		log.Info().Msg("Stopped streaming")
	}()

	if err := json.NewDecoder(r.Body).Decode(&buildData); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := DB.First(&repository, models.Repository{RepositoryName: buildData.RepositoryName}).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	go buildContainer(repository, buildData.BuildPack, buildData.CommitHash, buildlogChannel)
	for {
		data := <-buildlogChannel
		if string(data) == "EOF" {
			break
		}
		w.Write(data)
		flusher.Flush()
	}
}

//Create controller
func Create(w http.ResponseWriter, r *http.Request) {

}

//Stop controller
func Stop(w http.ResponseWriter, r *http.Request) {

}

func BuildLogs(w http.ResponseWriter, r *http.Request) {
	_channel := make(chan []byte)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	notify := r.Context().Done()
	go func() {
		<-notify
		log.Info().Msg("Stopped streaming")
	}()

	go _buildLogs(_channel)

	for {
		data := <-_channel
		if string(data) == "EOF" {
			break
		}
		w.Write(data)
		flusher.Flush()
	}

}

func _buildLogs(_channel chan []byte) {
	_channel <- []byte("golang build pack detected \n")
	time.Sleep(2 * time.Second)
	_channel <- []byte("stopping runnning app instance \n")
	time.Sleep(2 * time.Second)
	_channel <- []byte("Building app \n")
	time.Sleep(2 * time.Second)
	_channel <- []byte("Configuring app runtime \n")
	time.Sleep(2 * time.Second)
	_channel <- []byte("Installing dependencies \n \n")
	time.Sleep(2 * time.Second)
	_channel <- []byte("Build Succeded \n")
	time.Sleep(2 * time.Second)
	_channel <- []byte("Starting app instance \n")
	time.Sleep(2 * time.Second)
	_channel <- []byte("EOF")
	close(_channel)
}

//Getlogs controller
func Getlogs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repository := models.Repository{}
	query := make(map[string]interface{})
	query["uuid"] = params["uuid"]
	err := DB.Where(query).First(&repository).Error

	if err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	rdx, _ := container.Logs(r.Context(), repository.ContainerID)
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	logBuffer := make([]byte, 1024)
	notify := r.Context().Done()

	go func() {
		<-notify
		rdx.Close()
	}()

	for {
		n, err := rdx.Read(logBuffer)
		if err != nil {
			log.Error().Msgf("Error in reading %s", err.Error())
			rdx.Close()
			break
		}
		data := logBuffer[0:n]
		fmt.Fprintf(w, "data: %s\n\n", string(data))
		flusher.Flush()
		for i := 0; i < n; i++ {
			logBuffer[i] = 0
		}
	}
}

//Rebuild controller
func Rebuild(w http.ResponseWriter, r *http.Request) {

}

/*
* GO routines
 */
func buildContainer(repository models.Repository, build_pack, commit_hash string, buildlogChannel chan []byte) {
	appication_url := fmt.Sprintf("------> https://%s.upweb.io", repository.RepositoryName)
	buildlogChannel <- []byte("------> Starting your build\n\n")
	user := models.User{}
	DB.First(&user, repository.UserID)
	buildlogChannel <- []byte("------> Get application settings | variables \n")
	port, err := freeport.GetFreePort()
	buildlogChannel <- []byte("------> Creating application runtime \n")
	if repository.ContainerID != "" {
		log.Info().Msgf("Stopping previous container")
		container.Stop(repository.ContainerID)
		container.Remove(repository.ContainerID)
	}
	log.Info().Msg("Copying Dockerfile to project directory")
	buildPackPath := path.Join(os.Getenv("PROJECT_DIR"), "buildpacks", build_pack, "Dockerfile")
	projectDockerFile := path.Join(repository.PathDocker, "Dockerfile")
	err = fs.Copy(buildPackPath, projectDockerFile)
	err = fs.ReplaceStr(projectDockerFile, "dnow_replace_me", repository.RepositoryName)
	fmt.Sprint(err)
	log.Info().Msg("Building Docker image ... started")
	buildlogChannel <- []byte("------> Installing Dependencies \n")
	container.BuildImage(repository.PathDocker, repository.RepositoryName)
	log.Info().Msg("Building Docker image ... done")

	buildlogChannel <- []byte("------> Building Dependencies \n")
	log.Info().Msg("Building Docker Container ... started")
	containerID := container.Create(repository.RepositoryName, port)
	log.Info().Msgf("Building Docker Container ... done %s", containerID)
	buildlogChannel <- []byte("------> Built Successfully! \n")
	repository.ContainerID = containerID
	build := models.Build{
		CommitHash:   commit_hash,
		Email:        user.Email,
		UserID:       repository.UserID,
		RepositoryID: repository.ID,
	}
	DB.Save(repository)
	DB.Create(&build)
	buildlogChannel <- []byte("------> Launching... \n")
	buildlogChannel <- []byte(appication_url)
	buildlogChannel <- []byte("EOF")
	close(buildlogChannel)
}
