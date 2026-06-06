package views

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sadra/tui-dashboard/internal/models"
	"github.com/sadra/tui-dashboard/internal/theme"
	"github.com/sadra/tui-dashboard/internal/workers"
)

type DockerMetricsMsg models.DockerMetrics

type DockerView struct {
	width           int
	height          int
	metrics         models.DockerMetrics
	worker          *workers.DockerWorker
	ctx             context.Context
	cancel          context.CancelFunc
	metricsChad     chan models.DockerMetrics
	started         bool
	selectedIndex   int
	showLogs        bool
	selectedLogs    string
	selectedLogName string
}

func NewDockerView() *DockerView {
	ctx, cancel := context.WithCancel(context.Background())
	return &DockerView{
		worker:      workers.NewDockerWorker(2 * time.Second),
		ctx:         ctx,
		cancel:      cancel,
		metricsChad: make(chan models.DockerMetrics, 1),
		metrics:     models.DockerMetrics{Containers: []models.DockerContainer{}},
	}
}

func (dv *DockerView) StartWorker() tea.Cmd {
	if dv.started {
		return nil
	}
	dv.started = true

	go dv.worker.StartPolling(dv.ctx, dv.metricsChad)
	return dv.waitForMetricsCmd()
}

func (dv *DockerView) Update(msg interface{}) tea.Cmd {
	switch msg := msg.(type) {
	case DockerMetricsMsg:
		dv.metrics = models.DockerMetrics(msg)
		if dv.selectedIndex >= len(dv.metrics.Containers) && len(dv.metrics.Containers) > 0 {
			dv.selectedIndex = len(dv.metrics.Containers) - 1
		}
		return dv.waitForMetricsCmd()

	case tea.KeyMsg:
		if dv.showLogs {
			if msg.String() == "esc" || msg.String() == "q" {
				dv.showLogs = false
			}
			return nil
		}

		switch msg.String() {
		case "j", "down":
			if dv.selectedIndex < len(dv.metrics.Containers)-1 {
				dv.selectedIndex++
			}
		case "k", "up":
			if dv.selectedIndex > 0 {
				dv.selectedIndex--
			}
		case "s":
			if dv.selectedIndex < len(dv.metrics.Containers) {
				container := dv.metrics.Containers[dv.selectedIndex]
				if container.Status == models.StatusRunning {
					go dv.worker.StopContainer(dv.ctx, container.ID)
				} else if container.Status == models.StatusExited {
					go dv.worker.StartContainer(dv.ctx, container.ID)
				}
			}
		case "r":
			if dv.selectedIndex < len(dv.metrics.Containers) {
				container := dv.metrics.Containers[dv.selectedIndex]
				go dv.worker.RestartContainer(dv.ctx, container.ID)
			}
		case "l":
			if dv.selectedIndex < len(dv.metrics.Containers) {
				container := dv.metrics.Containers[dv.selectedIndex]
				dv.selectedLogName = container.Name
				logs, err := dv.worker.GetContainerLogs(dv.ctx, container.ID, 20)
				if err != nil {
					dv.selectedLogs = "Error fetching logs: " + err.Error()
				} else {
					dv.selectedLogs = logs
				}
				dv.showLogs = true
			}
		}
	}
	return nil
}

func (dv *DockerView) waitForMetricsCmd() tea.Cmd {
	return func() tea.Msg {
		select {
		case m := <-dv.metricsChad:
			return DockerMetricsMsg(m)
		case <-dv.ctx.Done():
			return nil
		}
	}
}

func (dv *DockerView) SetSize(width, height int) {
	dv.width = width
	dv.height = height
}

func (dv *DockerView) statusColor(status models.ContainerStatus) lipgloss.Color {
	switch status {
	case models.StatusRunning:
		return theme.Color.Success
	case models.StatusExited:
		return theme.Color.Warning
	case models.StatusPaused:
		return theme.Color.Accent
	default:
		return theme.Color.Foreground
	}
}

func (dv *DockerView) Render() string {
	if dv.showLogs {
		return dv.renderLogs()
	}

	if dv.metrics.Error != "" {
		return theme.PaneStyle.Width(dv.width - 4).Render(
			fmt.Sprintf("%s\n\n%s\n\n%s",
				theme.HeaderStyle.Render("🐳 Docker Management"),
				theme.HelpStyle.Render("⚠️  "+dv.metrics.Error),
				theme.HelpStyle.Render("Tip: Make sure Docker daemon is running"),
			),
		)
	}

	if len(dv.metrics.Containers) == 0 {
		return theme.PaneStyle.Width(dv.width - 4).Render(
			fmt.Sprintf("%s\n\n%s",
				theme.HeaderStyle.Render("🐳 Docker Management"),
				theme.HelpStyle.Render("No containers found"),
			),
		)
	}

	containerLines := []string{
		theme.HeaderStyle.Render("🐳 Docker Management"),
		"",
		fmt.Sprintf("%-4s %-25s %-12s %s", "ID", "Name", "Status", "Ports"),
		strings.Repeat("─", dv.width-6),
	}

	for i, c := range dv.metrics.Containers {
		status := string(c.Status)
		statusColor := dv.statusColor(c.Status)
		styledStatus := lipgloss.NewStyle().
			Foreground(statusColor).
			Render(fmt.Sprintf("● %s", status))

		ports := strings.Join(c.Ports, ", ")
		if ports == "" {
			ports = "—"
		}

		line := fmt.Sprintf("%-4s %-25s %s %s", c.ID[:4], truncate(c.Name, 23), styledStatus, ports)

		if i == dv.selectedIndex {
			line = lipgloss.NewStyle().
				Background(theme.Color.Primary).
				Foreground(theme.Color.Background).
				Render(line)
		}

		containerLines = append(containerLines, line)
	}

	containerLines = append(containerLines, "")
	containerLines = append(containerLines, theme.HeaderStyle.Render("Actions"))
	containerLines = append(containerLines, "  (j/k) Navigate  (s) Start/Stop  (r) Restart  (l) Logs")

	content := strings.Join(containerLines, "\n")
	return theme.PaneStyle.Width(dv.width - 4).Render(content)
}

func (dv *DockerView) renderLogs() string {
	logs := dv.selectedLogs
	if len(logs) > 2000 {
		logs = "... (truncated)\n" + logs[len(logs)-2000:]
	}

	content := fmt.Sprintf(`%s - %s

%s

%s`,
		theme.HeaderStyle.Render("📋 Container Logs"),
		dv.selectedLogName,
		logs,
		theme.HelpStyle.Render("(q or Esc to close)"),
	)

	return theme.PaneStyle.Width(dv.width - 4).Render(content)
}

func (dv *DockerView) Cleanup() {
	dv.cancel()
	dv.worker.Cleanup()
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
