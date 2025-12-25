package models

type DiskInfo struct {
	MountPoint string
	UsedBytes  uint64
	TotalBytes uint64
}

// Constructor
func NewDiskInfo(mountPoint string, used, total uint64) *DiskInfo {
	// PSEUDOCODE:
	// - Create new DiskInfo with parameters
	// - Return pointer to struct
	return &DiskInfo{
		MountPoint: mountPoint,
		UsedBytes:  used,
		TotalBytes: total,
	}
}

// GetUsedPercent calculates used disk percentage
func (di *DiskInfo) GetUsedPercent() float64 {
	// - IF TotalBytes equals 0 THEN return 0.0
	if di.TotalBytes == 0 {
		return 0.0
	}
	// - Calculate: (UsedBytes / TotalBytes) * 100
	diskPercentage := float64(di.UsedBytes) / float64(di.TotalBytes) * 100

	return diskPercentage
}

// GetFreePercent calculates free disk percentage
func (di *DiskInfo) GetFreePercent() float64 {
	// - Calculate: 100.0 - GetUsedPercent()
	// - Return result
	freePercentage := 100.0 - di.GetUsedPercent()
	return freePercentage
}

// GetStatus determines disk health status
func (di *DiskInfo) GetStatus(thresholds *Thresholds) string {
	// - Get free percentage
	freePercentage := di.GetFreePercent()
	// - IF free < thresholds.DiskCritical THEN return "CRITICAL"
	// - ELSE IF free < thresholds.DiskWarning THEN return "WARNING"
	// - ELSE return "OK"
	if freePercentage < thresholds.DiskCritical {
		return "CRITICAL"
	} else if freePercentage < thresholds.DiskWarning {
		return "WARNING"
	} else {
		return "OK"
	}
}
