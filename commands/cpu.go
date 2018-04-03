package commands

import (
	"github.com/shirou/gopsutil/cpu"
)

//region CPUOutput

// CPUOutput contains information about the CPU
type CPUOutput struct {
	Cores     int32
	ModelName string
	Procent   float64
	info      bool
}

func (CPUOutput) Error() string { return "" }

func (c CPUOutput) Message() OutputMessage {
	if c.info {
		return map[string]interface{}{
			"cores":      c.Cores,
			"model_name": c.ModelName,
		}
	}
	return map[string]interface{}{
		"procent": c.Procent,
	}
}

//endregion

//region CPUCommand

// CPUCommand retrieves information about the CPU
type CPUCommand struct {
}

func (c CPUCommand) Run(cmd Command) Output {
	if cmd.HasArgument("-info") {
		info, err := cpu.Info()
		if err != nil {
			return ErrorOutput{err: err.Error()}
		}

		if len(info) <= 0 {
			return ErrorOutput{err: "can't retrieve CPU info"}
		}

		i := info[0]
		return CPUOutput{
			Cores:     i.Cores,
			ModelName: i.ModelName,
			info:      true,
		}
	}

	// Check the CPU percentage from previous checks
	// there is a chance you get 0 procent if this is the first time checking
	// TODO: Is there a better way to do this?
	f, err := cpu.Percent(0, false)
	if err != nil {
		return ErrorOutput{err: err.Error()}
	}

	if len(f) <= 0 {
		return ErrorOutput{err: err.Error()}
	}

	i := f[0]
	return CPUOutput{Procent: i}
}

func (c CPUCommand) Name() string        { return "cpu" }
func (c CPUCommand) Description() string { return "returns information about the CPU and CPU usage" }
func (c CPUCommand) Usage() string       { return "[-info]" }

//endregion
