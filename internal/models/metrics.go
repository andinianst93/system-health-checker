package models

import (
	"time"
)

type SystemMetrics struct {
	CPUPercent  float64
	MemoryUsed  uint64
	MemoryTotal uint64
	Disks       []*DiskInfo
	Processes   []*ProcessInfo
	CheckTime   time.Time
}

func NewSystemMetrics() *SystemMetrics {
	// - Create new SystemMetrics struct
	// - Initialize Disks as empty slice: make([]*DiskInfo, 0)
	// - Initialize Processes as empty slice: make([]*ProcessInfo, 0)
	// - Set CheckTime to current time
	// - Return pointer to struct
	return &SystemMetrics{
		CPUPercent:  0.0,
		MemoryUsed:  0,
		MemoryTotal: 0,
		Disks:       make([]*DiskInfo, 0),
		Processes:   make([]*ProcessInfo, 0),
		CheckTime:   time.Now(),
	}
}

// GetMemoryPercent calculates memory usage percentage
func (sm *SystemMetrics) GetMemoryPercent() float64 {
	// - IF MemoryTotal equals 0 THEN return 0.0
	if sm.MemoryTotal == 0 {
		return 0.0
	}
	// - Calculate: (MemoryUsed / MemoryTotal) * 100
	memoryPercentage := float64(sm.MemoryUsed) / float64(sm.MemoryTotal) * 100

	return memoryPercentage
}
