package models

type Thresholds struct {
	CPUWarning   float64
	CPUCritical  float64
	MemWarning   float64
	MemCritical  float64
	DiskWarning  float64
	DiskCritical float64
}

func NewDefaultThresholds() *Thresholds {
	// - Create new Thresholds
	// - Set CPUWarning = 80.0
	// - Set CPUCritical = 90.0
	// - Set MemWarning = 75.0
	// - Set MemCritical = 85.0
	// - Set DiskWarning = 20.0 (20% free)
	// - Set DiskCritical = 10.0 (10% free)
	// - Return pointer to struct
	return &Thresholds{
		CPUWarning:   80.0,
		CPUCritical:  90.0,
		MemWarning:   75.0,
		MemCritical:  85.0,
		DiskWarning:  20.0,
		DiskCritical: 10.0,
	}
}
