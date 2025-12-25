package models

type ProcessInfo struct {
	PID           int32
	Name          string
	CPUPercent    float64
	MemoryPercent float64
	Status        string
}

func NewProcessInfo(pid int32, name string) *ProcessInfo {
	// - Create new ProcessInfo with parameters
	// - Initialize CPUPercent and MemoryPercent to 0
	// - Set Status to "unknown"
	// - Return pointer to struct
	return &ProcessInfo{
		PID:           pid,
		Name:          name,
		CPUPercent:    0,
		MemoryPercent: 0,
		Status:        "unknown",
	}
}
