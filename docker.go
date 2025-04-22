package main

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

const apiVersion = "1.45"

func StartTrainingContainer(args []string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion(apiVersion),
	)

	if err != nil {
		return "", err
	}

	jobID := uuid.New().String()
	containerName := "lora-job-" + jobID

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "training-job:latest",                   // your Go-based training image
		Cmd:   append([]string{"./lora-job"}, args...), // pass the command + args
		Tty:   false,
	}, nil, nil, nil, containerName)

	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return jobID, nil
}

func FetchContainerLogs(jobID string) (string, error) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion(apiVersion),
	)
	reader, err := cli.ContainerLogs(ctx, "lora-job-"+jobID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Tail: "100"})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	return string(logs), err
}

func GetContainerStatus(jobID string) (string, error) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion(apiVersion),
	)
	cJSON, err := cli.ContainerInspect(ctx, "lora-job-"+jobID)
	if err != nil {
		return "", err
	}
	return cJSON.State.Status, nil
}

func StopContainer(jobID string) error {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion(apiVersion),
	)
	return cli.ContainerStop(ctx, "lora-job-"+jobID, container.StopOptions{})
}
