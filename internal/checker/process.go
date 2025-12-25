package checker

import (
	"fmt"

	"github.com/andinianst93/system-health-checker/internal/models"
	"github.com/shirou/gopsutil/v4/process"
)

// CheckProcess checks if a specific process is running
func (hc *HealthChecker) CheckProcess(processName string) error {
	// - Call process.Processes() to get all processes
	processes, err := process.Processes()
	// - IF error THEN return error
	if err != nil {
		return err
	}

	// - FOR EACH proc IN processes:
	//     - Call proc.Name() to get process name
	//     - IF error THEN continue
	//     - IF name matches processName THEN
	//         - Get proc.Pid()
	//         - Get proc.CPUPercent()
	//         - Get proc.MemoryPercent()
	//         - Create processInfo = models.NewProcessInfo(pid, name)
	//         - Set processInfo.CPUPercent
	//         - Set processInfo.MemoryPercent
	//         - Set processInfo.Status = "running"
	//         - Append to hc.metrics.Processes
	//         - Return nil (found)

	for _, proc := range processes {
		name, err := proc.Name()
		if err != nil {
			continue
		}
		if name == processName {
			pid, err := proc.Ppid()
			if err != nil {
				continue
			}
			cpuPercent, err := proc.CPUPercent()
			if err != nil {
				continue
			}
			memPercent, err := proc.MemoryPercent()
			if err != nil {
				continue
			}
			processInfo := models.NewProcessInfo(pid, name)
			processInfo.CPUPercent = cpuPercent
			processInfo.MemoryPercent = float64(memPercent)
			processInfo.Status = "running"
			hc.metrics.Processes = append(hc.metrics.Processes, processInfo)
			return nil
		}
	}

	// - Return error with message "process not found: %s"
	return fmt.Errorf("process not found: %s", processName)
}
