package workers

import (
	"context"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sadra/tui-dashboard/internal/models"
)

type SystemWorker struct {
	interval      time.Duration
	lastNetStats  map[string]interface{}
	prevBytesSent uint64
	prevBytesRecv uint64
}

func NewSystemWorker(interval time.Duration) *SystemWorker {
	return &SystemWorker{
		interval: interval,
	}
}

func (sw *SystemWorker) FetchMetrics() (models.SystemMetrics, error) {
	metrics := models.SystemMetrics{}

	if err := sw.fetchCPU(&metrics); err != nil {
		metrics.LastErr = err.Error()
	}

	if err := sw.fetchMemory(&metrics); err != nil {
		metrics.LastErr = err.Error()
	}

	if err := sw.fetchDisk(&metrics); err != nil {
		metrics.LastErr = err.Error()
	}

	if err := sw.fetchNetwork(&metrics); err != nil {
		metrics.LastErr = err.Error()
	}

	if err := sw.fetchHost(&metrics); err != nil {
		metrics.LastErr = err.Error()
	}

	if err := sw.fetchProcesses(&metrics); err != nil {
		metrics.LastErr = err.Error()
	}

	return metrics, nil
}

func (sw *SystemWorker) fetchCPU(m *models.SystemMetrics) error {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}
	if len(percent) > 0 {
		m.CPU.Percent = percent[0]
	}

	cores, err := cpu.Counts(logical := true)
	if err == nil {
		m.CPU.Cores = cores
	}

	return nil
}

func (sw *SystemWorker) fetchMemory(m *models.SystemMetrics) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	m.Memory.UsedPercent = v.UsedPercent
	m.Memory.UsedGB = float64(v.Used) / (1024 * 1024 * 1024)
	m.Memory.TotalGB = float64(v.Total) / (1024 * 1024 * 1024)
	return nil
}

func (sw *SystemWorker) fetchDisk(m *models.SystemMetrics) error {
	usage, err := disk.Usage("/")
	if err != nil {
		return err
	}
	m.Disk.UsedPercent = usage.UsedPercent
	m.Disk.UsedGB = float64(usage.Used) / (1024 * 1024 * 1024)
	m.Disk.TotalGB = float64(usage.Total) / (1024 * 1024 * 1024)
	return nil
}

func (sw *SystemWorker) fetchNetwork(m *models.SystemMetrics) error {
	stats, err := net.IOCounters(false)
	if err != nil {
		return err
	}

	if len(stats) > 0 {
		stat := stats[0]

		if sw.prevBytesSent > 0 {
			bytesDiff := stat.BytesSent - sw.prevBytesSent
			m.Network.BytesSentPerSec = float64(bytesDiff)
		}

		if sw.prevBytesRecv > 0 {
			bytesDiff := stat.BytesRecv - sw.prevBytesRecv
			m.Network.BytesRecvPerSec = float64(bytesDiff)
		}

		sw.prevBytesSent = stat.BytesSent
		sw.prevBytesRecv = stat.BytesRecv
	}

	return nil
}

func (sw *SystemWorker) fetchHost(m *models.SystemMetrics) error {
	hostInfo, err := host.Info()
	if err == nil {
		m.Host.Uptime = hostInfo.Uptime
	}

	loadAvg, err := load.Avg()
	if err == nil {
		m.Host.LoadAvg1 = loadAvg.Load1
		m.Host.LoadAvg5 = loadAvg.Load5
		m.Host.LoadAvg15 = loadAvg.Load15
	}

	countInfo, err := process.Pids()
	if err == nil {
		m.Host.ProcessCount = len(countInfo)
	}

	return nil
}

func (sw *SystemWorker) fetchProcesses(m *models.SystemMetrics) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}

	topCPU := make([]models.ProcessInfo, 0, 5)
	topMemory := make([]models.ProcessInfo, 0, 5)

	for _, p := range processes {
		name, _ := p.Name()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()

		procInfo := models.ProcessInfo{
			PID:        p.Pid,
			Name:       name,
			CPUPercent: cpuPercent,
			MemPercent: float64(memPercent),
		}

		topCPU = append(topCPU, procInfo)
		topMemory = append(topMemory, procInfo)
	}

	sort.Slice(topCPU, func(i, j int) bool {
		return topCPU[i].CPUPercent > topCPU[j].CPUPercent
	})

	sort.Slice(topMemory, func(i, j int) bool {
		return topMemory[i].MemPercent > topMemory[j].MemPercent
	})

	if len(topCPU) > 5 {
		m.Processes.TopCPU = topCPU[:5]
	} else {
		m.Processes.TopCPU = topCPU
	}

	if len(topMemory) > 5 {
		m.Processes.TopMemory = topMemory[:5]
	} else {
		m.Processes.TopMemory = topMemory
	}

	return nil
}

func (sw *SystemWorker) StartPolling(ctx context.Context, ch chan<- models.SystemMetrics) {
	ticker := time.NewTicker(sw.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics, _ := sw.FetchMetrics()
			select {
			case ch <- metrics:
			case <-ctx.Done():
				return
			}
		}
	}
}
