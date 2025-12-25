package checker

import "github.com/shirou/gopsutil/v4/mem"

// CheckMemory gets current memory usage
func (hc *HealthChecker) CheckMemory() error {
	// - Call mem.VirtualMemory() to get memory stats
	vmem, err := mem.VirtualMemory()
	// - IF error THEN return error
	if err != nil {
		return err
	}
	// - Set hc.metrics.MemoryUsed = vmem.Used
	hc.metrics.MemoryUsed = vmem.Used
	// - Set hc.metrics.MemoryTotal = vmem.Total
	hc.metrics.MemoryTotal = vmem.Total
	// - Return nil
	return nil
}
