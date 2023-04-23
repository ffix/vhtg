package sources

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type DockerProcessor struct {
	handler EventHandler
	label   string
	logger  logger
}

func NewDockerProcessor(handler EventHandler, label string, logger logger) *DockerProcessor {
	return &DockerProcessor{handler: handler, label: label, logger: logger}
}

func (d *DockerProcessor) Process() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	for {
		containerID, err := getContainerIDByLabel(ctx, cli, d.label)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}

		options := types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Timestamps: true,
		}
		logsReader, err := cli.ContainerLogs(ctx, containerID, options)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}

		//io.Copy(os.Stdout, logsReader)
		d.processLogs(logsReader)
		logsReader.Close()

		// Sleep before retrying in case of container restart
		time.Sleep(5 * time.Second)
	}
}

func (d *DockerProcessor) processLogs(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Bytes()
		// Decode the Docker log header.
		header := line[:8]
		length := binary.BigEndian.Uint32(header[4:8])

		// Extract the log content.
		content := string(line[8 : 8+length])

		timestamp, message, err := parseTimestampAndMessage(content)
		if err != nil {
			d.logger.Warnf("Error parsing line: %w\n", err.Error())
			continue
		}

		d.handler.ProcessLine(message, &timestamp)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error while reading logs: %v\n", err)
	}
}

func parseTimestampAndMessage(line string) (time.Time, string, error) {
	// Find the first space.
	messageStart := 0
	for ; messageStart < len(line); messageStart++ {
		if line[messageStart] == ' ' {
			break
		}
	}

	// Docker uses RFC3339Nano as the timestamp format by default.
	timestamp, err := time.Parse(time.RFC3339Nano, line[:messageStart])
	if err != nil {
		return time.Time{}, "", err
	}

	// Everything after the space is the message.
	message := line[messageStart+1:]

	return timestamp, message, nil
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
