package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

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
	if err := json.NewDecoder(r.Body).Decode(&buildData); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	repository := models.Repository{}

	if err := DB.First(&repository, models.Repository{RepositoryName: buildData.RepositoryName}).Error; err != nil {
		Helper.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Info().Msgf("Building app async mode -> %s", buildData.RepositoryName)
	go buildContainer(repository, buildData.BuildPack, buildData.CommitHash)
	Helper.RespondWithJSON(w, 200, map[string]string{"message": "build success"})
}

//Create controller
func Create(w http.ResponseWriter, r *http.Request) {

}

//Stop controller
func Stop(w http.ResponseWriter, r *http.Request) {

}

//Getlogs controller
func Getlogs(w http.ResponseWriter, r *http.Request) {
	rdx, _ := container.Logs(r.Context(), "dc17da36440f2324ae47c070eaecdb8b773f6ca5c47f56f6fbb40d240178f249")
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
		w.Write(data)
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
func buildContainer(repository models.Repository, build_pack, commit_hash string) {
	log.Info().Msgf("Starting go routine to build container")
	user := models.User{}
	DB.First(&user, repository.UserID)
	port, err := freeport.GetFreePort()
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
	container.BuildImage(repository.PathDocker, repository.RepositoryName)
	log.Info().Msg("Building Docker image ... done")

	log.Info().Msg("Building Docker Container ... started")
	containerID := container.Create(repository.RepositoryName, port)
	log.Info().Msgf("Building Docker Container ... done %s", containerID)
	repository.ContainerID = containerID
	build := models.Build{
		CommitHash:   commit_hash,
		Email:        user.Email,
		UserID:       repository.UserID,
		RepositoryID: repository.ID,
	}
	DB.Save(repository)
	DB.Create(&build)
}
