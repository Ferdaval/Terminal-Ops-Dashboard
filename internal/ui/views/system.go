package views

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sadra/tui-dashboard/internal/models"
	"github.com/sadra/tui-dashboard/internal/theme"
	"github.com/sadra/tui-dashboard/internal/utils"
	"github.com/sadra/tui-dashboard/internal/workers"
)

type SystemMetricsMsg models.SystemMetrics

type SystemView struct {
	width           int
	height          int
	metrics         models.SystemMetrics
	worker          *workers.SystemWorker
	ctx             context.Context
	cancel          context.CancelFunc
	metricsChad     chan models.SystemMetrics
	started         bool
	cpuHistory      []float64
	memoryHistory   []float64
	networkSentHist []float64
	networkRecvHist []float64
	historyMaxSize  int
}

func NewSystemView() *SystemView {
	ctx, cancel := context.WithCancel(context.Background())
	const historySize = 60
	return &SystemView{
		worker:          workers.NewSystemWorker(1 * time.Second),
		ctx:             ctx,
		cancel:          cancel,
		metricsChad:     make(chan models.SystemMetrics, 1),
		cpuHistory:      make([]float64, 0, historySize),
		memoryHistory:   make([]float64, 0, historySize),
		networkSentHist: make([]float64, 0, historySize),
		networkRecvHist: make([]float64, 0, historySize),
		historyMaxSize:  historySize,
	}
}

func (sv *SystemView) StartWorker() tea.Cmd {
	if sv.started {
		return nil
	}
	sv.started = true

	go sv.worker.StartPolling(sv.ctx, sv.metricsChad)
	return sv.waitForMetricsCmd()
}

func (sv *SystemView) Update(msg interface{}) tea.Cmd {
	switch msg := msg.(type) {
	case SystemMetricsMsg:
		sv.metrics = models.SystemMetrics(msg)
		sv.updateHistory()
		return sv.waitForMetricsCmd()
	}
	return nil
}

func (sv *SystemView) updateHistory() {
	sv.cpuHistory = append(sv.cpuHistory, sv.metrics.CPU.Percent)
	sv.memoryHistory = append(sv.memoryHistory, sv.metrics.Memory.UsedPercent)
	sv.networkSentHist = append(sv.networkSentHist, sv.metrics.Network.BytesSentPerSec/1024/1024)
	sv.networkRecvHist = append(sv.networkRecvHist, sv.metrics.Network.BytesRecvPerSec/1024/1024)

	if len(sv.cpuHistory) > sv.historyMaxSize {
		sv.cpuHistory = sv.cpuHistory[1:]
		sv.memoryHistory = sv.memoryHistory[1:]
		sv.networkSentHist = sv.networkSentHist[1:]
		sv.networkRecvHist = sv.networkRecvHist[1:]
	}
}

func (sv *SystemView) waitForMetricsCmd() tea.Cmd {
	return func() tea.Msg {
		select {
		case m := <-sv.metricsChad:
			return SystemMetricsMsg(m)
		case <-sv.ctx.Done():
			return nil
		}
	}
}

func (sv *SystemView) SetSize(width, height int) {
	sv.width = width
	sv.height = height
}

func (sv *SystemView) Render() string {
	gaugeWidth := 30
	cpuSparkline := ""
	if len(sv.cpuHistory) > 0 {
		cpuSparkline = utils.RenderSparkline("  History", sv.cpuHistory, 100)
	}

	memSparkline := ""
	if len(sv.memoryHistory) > 0 {
		memSparkline = utils.RenderSparkline("  History", sv.memoryHistory, 100)
	}

	sentSparkline := ""
	maxSent := 100.0
	if len(sv.networkSentHist) > 0 {
		for _, v := range sv.networkSentHist {
			if v > maxSent {
				maxSent = v
			}
		}
		if maxSent == 0 {
			maxSent = 1
		}
		sentSparkline = utils.RenderSparkline("  Upload", sv.networkSentHist, maxSent)
	}

	recvSparkline := ""
	maxRecv := 100.0
	if len(sv.networkRecvHist) > 0 {
		for _, v := range sv.networkRecvHist {
			if v > maxRecv {
				maxRecv = v
			}
		}
		if maxRecv == 0 {
			maxRecv = 1
		}
		recvSparkline = utils.RenderSparkline("  Download", sv.networkRecvHist, maxRecv)
	}

	uptimeStr := sv.formatUptime(sv.metrics.Host.Uptime)
	topCPUStr := sv.renderTopProcesses("CPU", sv.metrics.Processes.TopCPU, true)
	topMemStr := sv.renderTopProcesses("Memory", sv.metrics.Processes.TopMemory, false)

	content := fmt.Sprintf(`%s

%s
%s
%s (%d cores)

%s
%s
%s

%s
%s (Used: %.1f GB / %.1f GB)

Network:
  ⬆️  Current: %s
%s
  ⬇️  Current: %s
%s

System Info:
  ⏱️  Uptime: %s  |  Load: %.2f / %.2f / %.2f
  ⚙️  Processes: %d

%s

%s
`,
		theme.HeaderStyle.Render("📊 System Monitoring"),
		theme.HeaderStyle.Render("CPU"),
		utils.RenderGauge("  Current", sv.metrics.CPU.Percent, gaugeWidth),
		cpuSparkline,
		sv.metrics.CPU.Cores,
		theme.HeaderStyle.Render("Memory"),
		utils.RenderGauge("  Current", sv.metrics.Memory.UsedPercent, gaugeWidth),
		memSparkline,
		theme.HeaderStyle.Render("Disk"),
		utils.RenderGauge("  Usage", sv.metrics.Disk.UsedPercent, gaugeWidth),
		sv.metrics.Disk.UsedGB,
		sv.metrics.Disk.TotalGB,
		utils.FormatBytesPerSec(sv.metrics.Network.BytesSentPerSec),
		sentSparkline,
		utils.FormatBytesPerSec(sv.metrics.Network.BytesRecvPerSec),
		recvSparkline,
		uptimeStr,
		sv.metrics.Host.LoadAvg1,
		sv.metrics.Host.LoadAvg5,
		sv.metrics.Host.LoadAvg15,
		sv.metrics.Host.ProcessCount,
		sv.renderProcessesSection(topCPUStr, topMemStr),
		theme.HelpStyle.Render("Updates: every 1 second | 1=System 2=Docker 3=GitHub | Last 60 seconds shown in sparklines"),
	)

	return theme.PaneStyle.Width(sv.width - 4).Render(content)
}

func (sv *SystemView) formatUptime(seconds uint64) string {
	if seconds == 0 {
		return "—"
	}
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	mins := (seconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}

func (sv *SystemView) renderTopProcesses(label string, procs []models.ProcessInfo, isCPU bool) string {
	if len(procs) == 0 {
		return ""
	}

	result := fmt.Sprintf("  Top 5 by %s:\n", label)
	for i, p := range procs {
		if i >= 5 {
			break
		}
		name := p.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		if isCPU {
			result += fmt.Sprintf("    %d. %-20s %5.1f%%\n", i+1, name, p.CPUPercent)
		} else {
			result += fmt.Sprintf("    %d. %-20s %5.1f%%\n", i+1, name, p.MemPercent)
		}
	}

	return result
}

func (sv *SystemView) renderProcessesSection(topCPUStr, topMemStr string) string {
	if topCPUStr == "" && topMemStr == "" {
		return ""
	}

	section := theme.HeaderStyle.Render("Top Processes") + "\n"
	section += topCPUStr
	if topCPUStr != "" && topMemStr != "" {
		section += "\n"
	}
	section += topMemStr

	return section
}

func (sv *SystemView) Cleanup() {
	sv.cancel()
}
