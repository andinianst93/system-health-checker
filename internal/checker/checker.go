package checker

import (
	"errors"
	"fmt"

	"github.com/andinianst93/system-health-checker/internal/models"
)

// HealthChecker performs system health checks
type HealthChecker struct {
	metrics      *models.SystemMetrics
	thresholds   *models.Thresholds
	outputFormat string
}

// Constructor with validation
func NewHealthChecker(thresholds *models.Thresholds, format string) (*HealthChecker, error) {
	// - IF format is not "table" AND not "json" THEN
	//     return nil, error with message "invalid format: must be 'table' or 'json'"

	if format != "table" && format != "json" {
		return nil, errors.New("invalid format: must be table or json")
	}
	// - IF thresholds is nil THEN
	//     set thresholds = models.NewDefaultThresholds()
	if thresholds == nil {
		thresholds = models.NewDefaultThresholds()
	}

	// - Create new HealthChecker struct
	// - Set metrics = models.NewSystemMetrics()
	// - Set thresholds from parameter
	// - Set outputFormat from parameter
	// - Return pointer to HealthChecker, nil error
	return &HealthChecker{
		models.NewSystemMetrics(),
		thresholds,
		format,
	}, nil
}

// CheckAll runs all health checks
func (hc *HealthChecker) CheckAll() error {
	// - Call hc.CheckCPU()
	err := hc.CheckCPU()
	//   IF error THEN return wrapped error "CPU check failed: %w"
	if err != nil {
		return fmt.Errorf("CPU check failed: %w", err)
	}
	// - Call hc.CheckMemory()
	err = hc.CheckMemory()
	//   IF error THEN return wrapped error "memory check failed: %w"
	if err != nil {
		return fmt.Errorf("memory check failed: %w", err)
	}
	// - Call hc.CheckDisk()
	err = hc.CheckDisk()
	//   IF error THEN return wrapped error "disk check failed: %w"
	if err != nil {
		return fmt.Errorf("disk check failed: %w", err)
	}
	// - Return nil (success)
	return nil
}

// GetMetrics returns pointer to metrics
func (hc *HealthChecker) GetMetrics() *models.SystemMetrics {
	// - Return hc.metrics
	return hc.metrics
}

// GetOverallStatus determines overall system health
func (hc *HealthChecker) GetOverallStatus() string {
	// - Check CPU critical: IF CPUPercent >= CPUCritical THEN return "CRITICAL"
	if hc.metrics.CPUPercent >= hc.thresholds.CPUCritical {
		return "CRITICAL"
	}
	// - Calculate memory percentage
	memPercent := hc.metrics.GetMemoryPercent()
	// - Check memory critical: IF memPercent >= MemCritical THEN return "CRITICAL"
	if memPercent >= hc.thresholds.MemCritical {
		return "CRITICAL"
	}

	// - FOR EACH disk IN metrics.Disks:
	//     IF disk.GetStatus(thresholds) equals "CRITICAL" THEN return "CRITICAL"
	for _, disk := range hc.metrics.Disks {
		if disk.GetStatus(hc.thresholds) == "CRITICAL" {
			return "CRITICAL"
		}
	}
	// - Check for warnings (same logic with Warning thresholds)
	// - IF any warning found THEN return "WARNING"

	for _, disk := range hc.metrics.Disks {
		if disk.GetStatus(hc.thresholds) == "WARNING" {
			return "WARNING"
		}
	}
	// - Return "OK"
	return "OK"
}
