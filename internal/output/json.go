package output

import (
	"encoding/json"
	"fmt"

	"github.com/andinianst93/system-health-checker/internal/models"
)

// JSONOutput represents the JSON output structure
type JSONOutput struct {
	Timestamp     string      `json:"timestamp"`
	OverallStatus string      `json:"overall_status"`
	Metrics       MetricsJSON `json:"metrics"`
}

type MetricsJSON struct {
	CPU       CPUMetric       `json:"cpu"`
	Memory    MemoryMetric    `json:"memory"`
	Disks     []DiskMetric    `json:"disks"`
	Processes []ProcessMetric `json:"processes,omitempty"`
}

type CPUMetric struct {
	Percent float64 `json:"percent"`
	Status  string  `json:"status"`
}

type MemoryMetric struct {
	UsedBytes  uint64  `json:"used_bytes"`
	TotalBytes uint64  `json:"total_bytes"`
	Percent    float64 `json:"percent"`
	Status     string  `json:"status"`
}

type DiskMetric struct {
	MountPoint  string  `json:"mount_point"`
	UsedBytes   uint64  `json:"used_bytes"`
	TotalBytes  uint64  `json:"total_bytes"`
	UsedPercent float64 `json:"used_percent"`
	FreePercent float64 `json:"free_percent"`
	Status      string  `json:"status"`
}

type ProcessMetric struct {
	Name          string  `json:"name"`
	PID           int32   `json:"pid"`
	Status        string  `json:"status"`
	MemoryPercent float64 `json:"memory_percent"`
}

// PrintJSON displays metrics in JSON format
func PrintJSON(metrics *models.SystemMetrics, thresholds *models.Thresholds) {
	// Build base JSON output
	jsonOutput := JSONOutput{
		Timestamp: metrics.CheckTime.Format("2006-01-02T15:04:05Z"),
	}

	// Determine overall status (mirror logic from checker.GetOverallStatus)
	overall := "OK"

	// CPU critical check
	if metrics.CPUPercent >= thresholds.CPUCritical {
		overall = "CRITICAL"
	} else {
		// Memory checks and disk critical checks
		memPercent := metrics.GetMemoryPercent()
		if memPercent >= thresholds.MemCritical {
			overall = "CRITICAL"
		} else {
			// Any disk critical?
			for _, d := range metrics.Disks {
				if d.GetStatus(thresholds) == "CRITICAL" {
					overall = "CRITICAL"
					break
				}
			}
			// If not critical, check warnings
			if overall != "CRITICAL" {
				if metrics.CPUPercent >= thresholds.CPUWarning || memPercent >= thresholds.MemWarning {
					overall = "WARNING"
				} else {
					for _, d := range metrics.Disks {
						if d.GetStatus(thresholds) == "WARNING" {
							overall = "WARNING"
							break
						}
					}
				}
			}
		}
	}

	jsonOutput.OverallStatus = overall

	// Build metrics JSON
	var mj MetricsJSON

	// CPU
	cpuStatus := "OK"
	if metrics.CPUPercent >= thresholds.CPUCritical {
		cpuStatus = "CRITICAL"
	} else if metrics.CPUPercent >= thresholds.CPUWarning {
		cpuStatus = "WARNING"
	}
	mj.CPU = CPUMetric{
		Percent: metrics.CPUPercent,
		Status:  cpuStatus,
	}

	// Memory
	memPercent := metrics.GetMemoryPercent()
	memStatus := "OK"
	if memPercent >= thresholds.MemCritical {
		memStatus = "CRITICAL"
	} else if memPercent >= thresholds.MemWarning {
		memStatus = "WARNING"
	}
	mj.Memory = MemoryMetric{
		UsedBytes:  metrics.MemoryUsed,
		TotalBytes: metrics.MemoryTotal,
		Percent:    memPercent,
		Status:     memStatus,
	}

	// Disks
	mj.Disks = make([]DiskMetric, 0, len(metrics.Disks))
	for _, d := range metrics.Disks {
		dm := DiskMetric{
			MountPoint:  d.MountPoint,
			UsedBytes:   d.UsedBytes,
			TotalBytes:  d.TotalBytes,
			UsedPercent: d.GetUsedPercent(),
			FreePercent: d.GetFreePercent(),
			Status:      d.GetStatus(thresholds),
		}
		mj.Disks = append(mj.Disks, dm)
	}

	// Processes (optional)
	if len(metrics.Processes) > 0 {
		mj.Processes = make([]ProcessMetric, 0, len(metrics.Processes))
		for _, p := range metrics.Processes {
			pm := ProcessMetric{
				Name:          p.Name,
				PID:           p.PID,
				Status:        p.Status,
				MemoryPercent: p.MemoryPercent,
			}
			mj.Processes = append(mj.Processes, pm)
		}
	}

	jsonOutput.Metrics = mj

	// Marshal with indentation
	out, err := json.MarshalIndent(jsonOutput, "", "  ")
	if err != nil {
		fmt.Println("error marshaling JSON:", err)
		return
	}

	fmt.Println(string(out))
}
