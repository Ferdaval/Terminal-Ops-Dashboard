package models

type SystemMetrics struct {
	CPU       CPUMetrics
	Memory    MemoryMetrics
	Disk      DiskMetrics
	Network   NetworkMetrics
	Host      HostMetrics
	Processes ProcessMetrics
	LastErr   string
}

type CPUMetrics struct {
	Percent float64
	Cores   int
}

type MemoryMetrics struct {
	UsedPercent float64
	UsedGB      float64
	TotalGB     float64
}

type DiskMetrics struct {
	UsedPercent float64
	UsedGB      float64
	TotalGB     float64
}

type NetworkMetrics struct {
	BytesSentPerSec   float64
	BytesRecvPerSec   float64
}

type HostMetrics struct {
	Uptime     uint64
	LoadAvg1   float64
	LoadAvg5   float64
	LoadAvg15  float64
	Temperature float64
	ProcessCount int
}

type ProcessInfo struct {
	PID     int32
	Name    string
	CPUPercent float64
	MemPercent float64
}

type ProcessMetrics struct {
	TopCPU []ProcessInfo
	TopMemory []ProcessInfo
}
