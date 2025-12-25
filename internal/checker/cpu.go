package checker

import (
	"github.com/shirou/gopsutil/v4/cpu"
)

// CheckCPU gets current CPU usage
func (hc *HealthChecker) CheckCPU() error {
	// PSEUDOCODE:
	// - Call cpu.Percent(0, false) to get CPU percentage
	percent, err := cpu.Percent(0, false)
	//   Parameters: interval=0 (instant), percpu=false (total)
	// - IF error THEN return error
	if err != nil {
		return err
	}
	// - IF percent array length > 0 THEN
	//     set hc.metrics.CPUPercent = percent[0]
	if len(percent) > 0 {
		hc.metrics.CPUPercent = percent[0]
	}
	// - Return nil
	return nil
}
