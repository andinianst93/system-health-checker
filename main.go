package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/andinianst93/system-health-checker/internal/checker"
	"github.com/andinianst93/system-health-checker/internal/models"
	"github.com/andinianst93/system-health-checker/internal/output"
)

func main() {
	// CLI flags
	format := flag.String("format", "table", "Output format (table|json)")
	processName := flag.String("process", "", "Check specific process by name (optional)")

	cpuWarning := flag.Float64("cpu-warning", -1.0, "CPU warning threshold (percent, optional)")
	cpuCritical := flag.Float64("cpu-critical", -1.0, "CPU critical threshold (percent, optional)")
	memWarning := flag.Float64("mem-warning", -1.0, "Memory warning threshold (percent, optional)")
	memCritical := flag.Float64("mem-critical", -1.0, "Memory critical threshold (percent, optional)")
	diskWarning := flag.Float64("disk-warning", -1.0, "Disk warning threshold (free percent, optional)")
	diskCritical := flag.Float64("disk-critical", -1.0, "Disk critical threshold (free percent, optional)")

	flag.Parse()

	// Build thresholds starting from defaults, then override any provided flags
	thresholds := models.NewDefaultThresholds()
	if *cpuWarning >= 0 {
		thresholds.CPUWarning = *cpuWarning
	}
	if *cpuCritical >= 0 {
		thresholds.CPUCritical = *cpuCritical
	}
	if *memWarning >= 0 {
		thresholds.MemWarning = *memWarning
	}
	if *memCritical >= 0 {
		thresholds.MemCritical = *memCritical
	}
	if *diskWarning >= 0 {
		thresholds.DiskWarning = *diskWarning
	}
	if *diskCritical >= 0 {
		thresholds.DiskCritical = *diskCritical
	}

	// Normalize format
	f := strings.ToLower(strings.TrimSpace(*format))

	// Create health checker
	hc, err := checker.NewHealthChecker(thresholds, f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create health checker:", err)
		os.Exit(3)
	}

	// If process flag provided, try to check that process (don't fail hard)
	if *processName != "" {
		if err := hc.CheckProcess(*processName); err != nil {
			// Print warning but continue
			fmt.Fprintf(os.Stderr, "process check warning: %v\n", err)
		}
	}

	// Run all checks
	if err := hc.CheckAll(); err != nil {
		fmt.Fprintln(os.Stderr, "health checks failed:", err)
		os.Exit(3)
	}

	metrics := hc.GetMetrics()
	overallStatus := hc.GetOverallStatus()

	// Output according to selected format
	switch f {
	case "json":
		output.PrintJSON(metrics, thresholds)
	default:
		// default to table
		output.PrintTable(metrics, thresholds)
	}

	// Exit code based on overall status
	switch overallStatus {
	case "CRITICAL":
		os.Exit(2)
	case "WARNING":
		os.Exit(1)
	default:
		os.Exit(0)
	}
}
