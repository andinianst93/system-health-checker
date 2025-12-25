package output

import (
	"fmt"
	"os"

	"github.com/andinianst93/system-health-checker/internal/models"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// PrintTable displays metrics in table format
func PrintTable(metrics *models.SystemMetrics, thresholds *models.Thresholds) {
	// Print header box
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║   SYSTEM HEALTH CHECK REPORT                                   ║")
	fmt.Printf("║   %s", metrics.CheckTime.Format("2006-01-02 15:04:05"))
	// Pad the rest of the line so the box looks okay
	fmt.Print("                                                     ║\n")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")

	// Create table writer
	table := tablewriter.NewWriter(os.Stdout)
	// Some versions of the tablewriter used here may not expose SetHeader/SetRowLine.
	// Append a header row manually so the header shows even without those helpers.
	table.Append([]string{"Metric", "Value", "Status", "Threshold"})
	// Add a simple divider row to visually separate the header from content.
	table.Append([]string{"------", "------", "------", "------"})

	// Track overall severity: 0=OK,1=WARNING,2=CRITICAL
	overallSeverity := 0
	incSeverity := func(s int) {
		if s > overallSeverity {
			overallSeverity = s
		}
	}

	// CPU row
	cpuValue := fmt.Sprintf("%.2f%%", metrics.CPUPercent)
	cpuStatusColored := getStatus(metrics.CPUPercent, thresholds.CPUWarning, thresholds.CPUCritical)
	// Determine raw CPU severity for overall computation
	if metrics.CPUPercent >= thresholds.CPUCritical {
		incSeverity(2)
	} else if metrics.CPUPercent >= thresholds.CPUWarning {
		incSeverity(1)
	}
	cpuThreshold := fmt.Sprintf("< %.0f%%", thresholds.CPUWarning)
	table.Append([]string{"CPU Usage", cpuValue, cpuStatusColored, cpuThreshold})

	// Memory row
	memPercent := metrics.GetMemoryPercent()
	usedGB := bytesToGB(metrics.MemoryUsed)
	totalGB := bytesToGB(metrics.MemoryTotal)
	memValue := fmt.Sprintf("%.2fGB / %.2fGB (%.1f%%)", usedGB, totalGB, memPercent)
	memStatusColored := getStatus(memPercent, thresholds.MemWarning, thresholds.MemCritical)
	// Raw memory severity
	if memPercent >= thresholds.MemCritical {
		incSeverity(2)
	} else if memPercent >= thresholds.MemWarning {
		incSeverity(1)
	}
	memThreshold := fmt.Sprintf("< %.0f%%", thresholds.MemWarning)
	table.Append([]string{"Memory Usage", memValue, memStatusColored, memThreshold})

	// Disks
	for _, d := range metrics.Disks {
		used := bytesToGB(d.UsedBytes)
		total := bytesToGB(d.TotalBytes)
		usedPercent := d.GetUsedPercent()
		diskValue := fmt.Sprintf("%.2fGB / %.2fGB (%.1f%% used)", used, total, usedPercent)
		// Disk.GetStatus uses free percent vs thresholds, so it returns raw status string ("OK","WARNING","CRITICAL")
		diskRawStatus := d.GetStatus(thresholds)
		diskStatusColored := colorizeStatus(diskRawStatus)
		// Determine severity from raw status
		switch diskRawStatus {
		case "CRITICAL":
			incSeverity(2)
		case "WARNING":
			incSeverity(1)
		}
		diskThreshold := fmt.Sprintf("< %.0f%% free", thresholds.DiskWarning)
		table.Append([]string{fmt.Sprintf("Disk %s", d.MountPoint), diskValue, diskStatusColored, diskThreshold})
	}

	// Render table
	table.Render()

	// Overall status
	var overall string
	switch overallSeverity {
	case 2:
		overall = colorizeStatus("CRITICAL")
	case 1:
		overall = colorizeStatus("WARNING")
	default:
		overall = colorizeStatus("OK")
	}
	fmt.Printf("\nOverall Status: %s\n", overall)
}

// getStatus returns a colorized status string based on numeric value and thresholds.
// value is compared against warning and critical thresholds (higher is worse).
func getStatus(value, warning, critical float64) string {
	if value >= critical {
		return colorizeStatus("CRITICAL")
	} else if value >= warning {
		return colorizeStatus("WARNING")
	} else {
		return colorizeStatus("OK")
	}
}

func colorizeStatus(status string) string {
	switch status {
	case "CRITICAL":
		return color.RedString("🔴 CRITICAL")
	case "WARNING":
		return color.YellowString("⚠️  WARNING")
	case "OK":
		return color.GreenString("✅ OK")
	default:
		return status
	}
}

func bytesToGB(bytes uint64) float64 {
	if bytes == 0 {
		return 0.0
	}
	return float64(bytes) / 1024.0 / 1024.0 / 1024.0
}
