package lib

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sjljrvis/deploynow/log"
)

// GenerateDefault will spin up default container
func GenerateDefault(name string, port int) string {
	imageName := "dnow-default"
	image_path := path.Join(os.Getenv("PROJECT_DIR"), "DockerFiles")
	BuildImage(image_path, imageName)
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

	logConfig := container.LogConfig{
		Type: "syslog",
		Config: map[string]string{
			"syslog-address":  "udp://" + os.Getenv("UPWEB_ECHO"),
			"tag":             name,
			"syslog-facility": "23",
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBinding,
		LogConfig:    logConfig,
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
	return _container.ID
}

func Create(image string, port int, variables []string) string {
	portString := strconv.Itoa(port)
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: portString,
	}

	portBinding := nat.PortMap{
		nat.Port("80/tcp"): []nat.PortBinding{hostBinding},
	}

	envs := append(variables, "PORT="+portString)

	containerConfig := &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			nat.Port("80/tcp"): struct{}{},
		},
		Env: envs,
	}

	logConfig := container.LogConfig{
		Type: "syslog",
		Config: map[string]string{
			"syslog-address":  "udp://" + os.Getenv("UPWEB_ECHO"),
			"tag":             image,
			"syslog-facility": "23",
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBinding,
		LogConfig:    logConfig,
	}

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error().Msgf("Unable to create docker client")
	}
	_container, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, image)
	if err != nil {
		log.Error().Msgf("Unable to create docker container %s", err.Error())
	}
	if err := cli.ContainerStart(context.Background(), _container.ID, types.ContainerStartOptions{}); err != nil {
		log.Error().Msgf("Unable to start docker container %s", err.Error())
	}
	log.Info().Msgf("Started container %s", _container.ID)
	return _container.ID
}

func Stop(containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error().Msgf(err.Error())
	}
	if err := cli.ContainerStop(ctx, containerID, nil); err != nil {
		log.Error().Msgf(err.Error())
	}
}

func Remove(containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error().Msgf(err.Error())
	}
	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		log.Error().Msgf(err.Error())
	}
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

func Logs(ctx context.Context, containerID string) (io.ReadCloser, error) {
	log.Info().Msgf("Reading logs %s", containerID)
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error().Msgf("Unable read logs", err.Error())
	}
	reader, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
	})
	return reader, err
}
