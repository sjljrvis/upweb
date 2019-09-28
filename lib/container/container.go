package lib

import (
	"context"
	"os"
	"os/exec"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sjljrvis/deploynow/log"
)

// GenerateDefault will spin up default container
func GenerateDefault(name string, port int) {
	imageName := "dnow-default"
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: strconv.Itoa(port),
	}
	portBinding := nat.PortMap{
		nat.Port("80/tcp"): []nat.PortBinding{hostBinding},
	}

	containerConfig := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			nat.Port("80/tcp"): struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBinding,
	}

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error().Msgf("Unable to create docker client")
	}
	_container, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, name)

	if err := cli.ContainerStart(context.Background(), _container.ID, types.ContainerStartOptions{}); err != nil {
		log.Error().Msgf("Unable to start docker container", err.Error())
	}
	log.Info().Msgf("Started container", _container.ID)
}

func Create(image string, port int) string {

	portString := strconv.Itoa(port)
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: portString,
	}
	portBinding := nat.PortMap{
		nat.Port(portString + "/tcp"): []nat.PortBinding{hostBinding},
	}

	containerConfig := &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			nat.Port(portString + "/tcp"): struct{}{},
		},
		Env: []string{"PORT=" + portString},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBinding,
	}

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error().Msgf("Unable to create docker client")
	}
	_container, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, image)

	if err := cli.ContainerStart(context.Background(), _container.ID, types.ContainerStartOptions{}); err != nil {
		log.Error().Msgf("Unable to start docker container", err.Error())
	}
	log.Info().Msgf("Started container %s", _container.ID)
	return _container.ID
}

func Stop(containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Stopping container ... ")
	if err := cli.ContainerStop(ctx, containerID, nil); err != nil {
		panic(err)
	}
	log.Info().Msgf("Stopping container ... done")
}

func Remove(containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Removing container ... ")
	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
	log.Info().Msgf("Removing container ... done")
}

func BuildImage(path, name string) {
	cmd := exec.Command("docker", "build", path, "-t", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Info().Err(err)
	}
}
