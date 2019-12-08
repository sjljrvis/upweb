package controller

import (
	"encoding/json"
	"fmt"
	"io"
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

type LogBroker struct {
	LogEvents      chan []byte
	newClients     chan chan []byte
	closingClients chan chan []byte
	clients        map[chan []byte]bool
}

func (broker *LogBroker) listen() {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.LogEvents:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan, _ := range broker.clients {
				clientMessageChan <- event
			}
		}
	}

}

//Getlogs controller
func Getlogs(w http.ResponseWriter, r *http.Request) {
	broker := &LogBroker{
		LogEvents:      make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	go broker.listen()
	rdx, _ := container.Logs(r.Context(), "dc17da36440f2324ae47c070eaecdb8b773f6ca5c47f56f6fbb40d240178f249")
	writeCmdOutput(w, r, rdx, broker)
}

func writeCmdOutput(res http.ResponseWriter, r *http.Request, reader io.ReadCloser, broker *LogBroker) {

	flusher, ok := res.(http.Flusher)

	if !ok {
		http.Error(res, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	logBuffer := make([]byte, 1024)

	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	// notify := rw.(http.CloseNotifier).CloseNotify()
	notify := r.Context().Done()

	go func() {
		<-notify
		broker.closingClients <- messageChan
	}()

	for {
		n, err := reader.Read(logBuffer)
		log.Error().Msgf("Error in reading %s", err.Error())
		if err != nil {
			reader.Close()
			break
		}
		data := logBuffer[0:n]
		log.Info().Msg(string(data))
		broker.LogEvents <- data
		//reset buffer
		for i := 0; i < n; i++ {
			logBuffer[i] = 0
		}
	}

	for {

		// Write to the ResponseWriter
		// Server Sent Events compatible
		fmt.Fprintf(res, "data: %s\n\n", <-messageChan)

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}

	// for {
	// 	n, err := reader.Read(logBuffer)
	// 	log.Error().Msgf("Error in reading %s", err.Error())
	// 	if err != nil {
	// 		reader.Close()
	// 		break
	// 	}
	// 	data := buffer[0:n]
	// 	log.Info().Msg(string(data))
	// 	res.Write(data)
	// 	if f, ok := res.(http.Flusher); ok {
	// 		f.Flush()
	// 	}
	// 	//reset buffer
	// 	for i := 0; i < n; i++ {
	// 		buffer[i] = 0
	// 	}
	// }

}

//Rebuild controller
func Rebuild(w http.ResponseWriter, r *http.Request) {

}

/*
* GO routines
 */
func buildContainer(repository models.Repository, build_pack, commit_hash string) {
	log.Info().Msgf("Starting go routine to build container")
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
		UserName:     repository.UserName,
		UserID:       repository.UserID,
		RepositoryID: repository.ID,
	}
	DB.Save(repository)
	DB.Create(&build)
}
