package client

import (
	"strings"

	"github.com/shirou/gopsutil/cpu"
)

// CPUInfo contains information about the CPU
type CPUInfo struct {
	Error     string `json:"error"`
	Cores     int32  `json:"cores"`
	ModelName string `json:"model_name"`
}

// CPUUsage contains about the CPU usage
type CPUUsage struct {
	Error   string  `json:"error"`
	Procent float64 `json:"procent"`
}

// CPUCheck retrieves information about the CPU
func CPUCheck(cmd Command) interface{} {
	info := false
	for _, args := range cmd.Params {
		if strings.ToLower(args.Name) == "-info" {
			info = true
		}
	}

	if info {
		info, err := cpu.Info()
		if err != nil {
			return CPUInfo{Error: err.Error()}
		}

		if len(info) <= 0 {
			return CPUInfo{Error: "Can't find CPU info"}
		}

		i := info[0]
		return CPUInfo{
			Cores:     i.Cores,
			ModelName: i.ModelName,
		}
	}

	// Check the CPU percentage from previous checks
	// there is a chance you get 0 procent if this is the first time checking
	// TODO: Is there a better way to do this?
	f, err := cpu.Percent(0, false)
	if err != nil {
		return CPUUsage{Error: err.Error()}
	}

	if len(f) <= 0 {
		return CPUInfo{Error: "Can't find CPU usage"}
	}

	i := f[0]

	return CPUUsage{Procent: i}
}
