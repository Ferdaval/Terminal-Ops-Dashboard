package workers

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sadra/tui-dashboard/internal/models"
)

type DockerWorker struct {
	interval time.Duration
	client   *client.Client
}

func NewDockerWorker(interval time.Duration) *DockerWorker {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return &DockerWorker{interval: interval, client: nil}
	}
	return &DockerWorker{
		interval: interval,
		client:   cli,
	}
}

func (dw *DockerWorker) FetchContainers(ctx context.Context) (models.DockerMetrics, error) {
	metrics := models.DockerMetrics{
		Containers: []models.DockerContainer{},
	}

	if dw.client == nil {
		metrics.Error = "Docker daemon not available"
		return metrics, nil
	}

	containers, err := dw.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		metrics.Error = "Failed to list containers: " + err.Error()
		return metrics, nil
	}

	for _, c := range containers {
		status := models.ContainerStatus(c.State)
		name := strings.TrimPrefix(c.Names[0], "/")

		ports := []string{}
		for _, p := range c.Ports {
			if p.PublicPort > 0 {
				ports = append(ports, p.IP+":"+string(rune(p.PublicPort)))
			}
		}

		metrics.Containers = append(metrics.Containers, models.DockerContainer{
			ID:      c.ID[:12],
			Name:    name,
			Image:   c.Image,
			Status:  status,
			Ports:   ports,
			State:   c.State,
		})
	}

	return metrics, nil
}

func (dw *DockerWorker) StartContainer(ctx context.Context, containerID string) error {
	if dw.client == nil {
		return nil
	}
	return dw.client.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (dw *DockerWorker) StopContainer(ctx context.Context, containerID string) error {
	if dw.client == nil {
		return nil
	}
	timeout := 10
	return dw.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
}

func (dw *DockerWorker) RestartContainer(ctx context.Context, containerID string) error {
	if dw.client == nil {
		return nil
	}
	timeout := 10
	return dw.client.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout})
}

func (dw *DockerWorker) GetContainerLogs(ctx context.Context, containerID string, lines int64) (string, error) {
	if dw.client == nil {
		return "", nil
	}

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "20",
	}

	reader, err := dw.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	logs := make([]byte, 0, 4096)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			logs = append(logs, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return string(logs), nil
}

func (dw *DockerWorker) StartPolling(ctx context.Context, ch chan<- models.DockerMetrics) {
	ticker := time.NewTicker(dw.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics, _ := dw.FetchContainers(ctx)
			select {
			case ch <- metrics:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (dw *DockerWorker) Cleanup() {
	if dw.client != nil {
		dw.client.Close()
	}
}
