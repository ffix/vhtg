package sources

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func ProcessDocker(handler EventHandler) {
	labelKey := os.Args[1]

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	for {
		containerID, err := getContainerIDByLabel(ctx, cli, labelKey)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}

		options := types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		}
		logsReader, err := cli.ContainerLogs(ctx, containerID, options)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}

		io.Copy(os.Stdout, logsReader)
		logsReader.Close()

		// Sleep before retrying in case of container restart
		time.Sleep(5 * time.Second)
	}
}

func getContainerIDByLabel(ctx context.Context, cli *client.Client, labelKey string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("label", labelKey)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filter,
		All:     false,
	})
	if err != nil {
		return "", err
	}

	if len(containers) == 0 {
		return "", fmt.Errorf("no container found with label %s", labelKey)
	}

	if len(containers) > 1 {
		return "", fmt.Errorf("multiple containers found with label %s", labelKey)
	}

	return containers[0].ID, nil
}
