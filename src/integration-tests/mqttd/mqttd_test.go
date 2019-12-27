package uhppote

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	containers := setup()
	code := m.Run()
	teardown(containers)

	os.Exit(code)
}

func setup() []*types.Container {
	containers := []*types.Container{}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	cli.NegotiateAPIVersion(ctx)

	container, err := cli.ContainerCreate(ctx, &container.Config{Image: "simulator"}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	list, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, c := range list {
		if c.ID == c.ID {
			fmt.Printf(" ... started Docker 'simulator' container %s\n", strings.TrimPrefix(c.Names[0], "/"))
			containers = append(containers, &c)
		}
	}

	return containers
}

func teardown(containers []*types.Container) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	cli.NegotiateAPIVersion(ctx)

	for _, c := range containers {
		fmt.Printf(" ... stopping Docker 'simulator' container %s\n", strings.TrimPrefix(c.Names[0], "/"))

		if err := cli.ContainerStop(ctx, c.ID, nil); err != nil {
			panic(err)
		}

		if err := cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
			panic(err)
		}
	}
}

func TestMQTTD(t *testing.T) {
	t.Errorf("NOT IMPLEMENTED YET")
}
