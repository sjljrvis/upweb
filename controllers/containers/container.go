package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/phayes/freeport"
	Helper "github.com/sjljrvis/deploynow/helpers"
	container "github.com/sjljrvis/deploynow/lib/container"
	fs "github.com/sjljrvis/deploynow/lib/fs"
	git "github.com/sjljrvis/deploynow/lib/git"
	nginx "github.com/sjljrvis/deploynow/lib/nginx"

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

//Build controller
func BuildFromGithub(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	// var buildData build
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

	if err := DB.First(&repository, params["id"]).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	go buildContainerFromGithub(repository, "node", buildlogChannel)
	for {
		data := <-buildlogChannel
		fmt.Fprintf(w, "data: %s\n\n", string(data)+"\n\n")
		flusher.Flush()
	}
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

//Create controller
func Create(w http.ResponseWriter, r *http.Request) {

}

//Stop controller
func Stop(w http.ResponseWriter, r *http.Request) {

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
		log.Info().Msgf("[BUILD] stopping container %s", repository.ContainerID)
		container.Stop(repository.ContainerID)
		container.Remove(repository.ContainerID)
	}
	log.Info().Msg("[BUILD] Copying Dockerfile to project directory")
	buildPackPath := path.Join(os.Getenv("PROJECT_DIR"), "buildpacks", build_pack, "Dockerfile")
	projectDockerFile := path.Join(repository.PathDocker, "Dockerfile")
	err = fs.Copy(buildPackPath, projectDockerFile)
	err = fs.ReplaceStr(projectDockerFile, "dnow_replace_me", repository.RepositoryName)
	fmt.Sprint(err)
	log.Info().Msg("[BUILD] Building Docker image")
	buildlogChannel <- []byte("------> Installing Dependencies \n")
	container.BuildImage(repository.PathDocker, repository.RepositoryName)
	log.Info().Msg("[BUILD] Docker image build complete")

	buildlogChannel <- []byte("------> Building Dependencies \n")
	log.Info().Msg("[BUILD ]Building Docker Container")
	containerID := container.Create(repository.RepositoryName, port, nil)
	nginx.WriteConfig(repository.RepositoryName, strconv.Itoa(port))
	nginx.Reload()
	log.Info().Msgf("[BUILD] Starting Docker Container %s", containerID)
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

func buildContainerFromGithub(repository models.Repository, build_pack string, buildlogChannel chan []byte) {
	// Remove contents from repo_docker dir here
	variables := []models.Variable{}
	envs := []string{}
	appication_url := fmt.Sprintf("------> https://%s.upweb.io", repository.RepositoryName)
	buildlogChannel <- []byte("------> Cleaning workspace\n\n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")
	fs.RemoveDir(repository.PathDocker)
	fs.CreateDir(repository.PathDocker)
	buildlogChannel <- []byte("------> Re-Initializing workspace \n\n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("------> Cloning code from github\n")
	git.Clone(repository.PathDocker, repository.GithubURL)
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")

	if repository.ContainerID != "" {
		log.Info().Msgf("[BUILD] stopping container %s", repository.ContainerID)
		container.Stop(repository.ContainerID)
		container.Remove(repository.ContainerID)
	}

	buildlogChannel <- []byte("------> Get application settings | variables \n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte(fmt.Sprintf("        Using BUILD_PACK as %s", build_pack))
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")

	DB.Find(&repository).Related(&variables)
	port, err := freeport.GetFreePort()

	buildlogChannel <- []byte("------> Creating application runtime \n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("        [PORT] will be set as environment variable dynamically \n")
	buildlogChannel <- []byte("\n")

	for _, v := range variables {
		vars := fmt.Sprintf("         %s=%s \n", v.Key, v.Value)
		envs = append(envs, fmt.Sprintf("%s=%s", v.Key, v.Value))
		buildlogChannel <- []byte(vars)
	}

	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")

	buildlogChannel <- []byte("------> Installing Dependencies \n")
	start_time := time.Now()
	log.Info().Msg("[BUILD] Copying Dockerfile to project directory")
	buildPackPath := path.Join(os.Getenv("PROJECT_DIR"), "buildpacks", build_pack, "Dockerfile")
	projectDockerFile := path.Join(repository.PathDocker, "Dockerfile")
	err = fs.Copy(buildPackPath, projectDockerFile)
	err = fs.ReplaceStr(projectDockerFile, "dnow_replace_me", repository.RepositoryName)
	fmt.Sprint(err)
	log.Info().Msg("[BUILD] Building Docker image")

	container.BuildImage(repository.PathDocker, repository.RepositoryName)
	log.Info().Msg("[BUILD] Docker image build complete")
	end_time := time.Since(start_time)
	time_diff := float32(end_time / time.Second)
	buildlogChannel <- []byte(fmt.Sprintf("        took %.2f seconds", time_diff))
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("------> Building Dependencies \n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")

	log.Info().Msg("[BUILD ]Building Docker Container")
	containerID := container.Create(repository.RepositoryName, port, envs)
	nginx.WriteConfig(repository.RepositoryName, strconv.Itoa(port))
	nginx.Reload()
	log.Info().Msgf("Building Docker Container ... done %s", containerID)

	buildlogChannel <- []byte("------> Build succeeded! \n")
	buildlogChannel <- []byte("\n")
	buildlogChannel <- []byte("\n")

	repository.ContainerID = containerID
	DB.Save(repository)
	buildlogChannel <- []byte("------> Launching... \n")
	buildlogChannel <- []byte("------> Application deployed on upweb")
	buildlogChannel <- []byte(appication_url)
	buildlogChannel <- []byte("EOF")
	close(buildlogChannel)
}
