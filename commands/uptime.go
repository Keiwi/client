package commands

import (
	"github.com/shirou/gopsutil/host"
)

//region UptimeOutput

// UptimeOutput contains information about uptime of a client
type UptimeOutput struct {
	Uptime uint64 `json:"uptime"`
}

func (UptimeOutput) Error() string { return "" }

func (u UptimeOutput) Message() OutputMessage {
	return map[string]interface{}{
		"uptime": u.Uptime,
	}
}

//endregion

//region UptimeCommand

// UptimeCommand checks how long a client has been on
type UptimeCommand struct {
}

func (u UptimeCommand) Run(cmd Command) Output {
	if cmd.HasArgument("-boot") {
		b, err := host.BootTime()
		if err != nil {
			return ErrorOutput{err: err.Error()}
		}
		return UptimeOutput{Uptime: b}
	}

	up, err := host.Uptime()
	if err != nil {
		return ErrorOutput{err: err.Error()}
	}
	return UptimeOutput{Uptime: up}
}

func (UptimeCommand) Name() string { return "uptime" }
func (UptimeCommand) Description() string {
	return "returns the clients boot time or uptime"
}
func (UptimeCommand) Usage() string { return `[-boot]` }

//endregion
