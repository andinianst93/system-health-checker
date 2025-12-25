package checker

import (
	"github.com/andinianst93/system-health-checker/internal/models"
	"github.com/shirou/gopsutil/v4/disk"
)

// CheckDisk gets disk usage for all partitions
func (hc *HealthChecker) CheckDisk() error {
	// - Call disk.Partitions(false) to get all partitions
	partitions, err := disk.Partitions(false)
	//   Parameter: all=false (only physical partitions)
	// - IF error THEN return error
	if err != nil {
		return err
	}
	// - FOR EACH partition IN partitions:
	//     - Call disk.Usage(partition.Mountpoint) to get usage
	//     - IF error THEN continue (skip this partition)
	//     - Create diskInfo = models.NewDiskInfo(
	//         partition.Mountpoint,
	//         usage.Used,
	//         usage.Total
	//       )
	//     - Append diskInfo to hc.metrics.Disks

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		diskInfo := models.NewDiskInfo(
			partition.Mountpoint,
			usage.Used,
			usage.Total,
		)
		hc.metrics.Disks = append(hc.metrics.Disks, diskInfo)
	}
	// - Return nil
	return nil
}
