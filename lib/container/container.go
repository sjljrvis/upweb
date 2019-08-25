package lib

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/phayes/freeport"
	"github.com/sjljrvis/deploynow/log"
)

func GenerateDefault() {

	port, err := freeport.GetFreePort()
	if err != nil {
		log.Error().Msgf("Error Occured in generating default container")
	}
	imageName := "dnow-default"
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: string(port),
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
	_container, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, "dnow-default-1")
	fmt.Println(err)
	if err := cli.ContainerStart(context.Background(), _container.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println("------")
		fmt.Println(err)
		log.Error().Msgf("Unable to start docker container")
	}

	fmt.Println(_container.ID)
}

func Create(port int, env []string) {

	imageName := "dnow-default"
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: string(port),
	}
	containerPort, _ := nat.NewPort("tcp", "80")
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	containerConfig := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			nat.Port("80/tcp"): {},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBinding,
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Error().Msgf("Unable to create docker client")
	}

	_container, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, "dnow-default-1")

	if err := cli.ContainerStart(context.Background(), _container.ID, types.ContainerStartOptions{}); err != nil {
		log.Error().Msgf("Unable to start docker container")
	}

	fmt.Println(_container.ID)
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
	fmt.Println("Success")
}
