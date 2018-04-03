package commands

import (
	"github.com/shirou/gopsutil/net"
)

//region NetworkOutput

// NetworkOutput contains information about a clients Network I/O
type NetworkOutput struct {
	Sent uint64
	Recv uint64
}

func (NetworkOutput) Error() string { return "" }

func (i NetworkOutput) Message() OutputMessage {
	return map[string]interface{}{
		"sent": i.Sent,
		"recv": i.Recv,
	}
}

//endregion

//region NetworkCommand

// NetworkCommand gathers the clients network I/O usage
type NetworkCommand struct {
	Version string
}

func (i NetworkCommand) Run(cmd Command) Output {
	io, err := net.IOCounters(false)
	if err != nil {
		return ErrorOutput{err: err.Error()}
	}
	if len(io) <= 0 {
		return ErrorOutput{err: "Can't get network I/O usage"}
	}

	return NetworkOutput{
		Sent: io[0].BytesSent,
		Recv: io[0].BytesRecv,
	}
}

func (NetworkCommand) Name() string { return "network" }
func (NetworkCommand) Description() string {
	return "returns the network I/O"
}
func (NetworkCommand) Usage() string { return `` }

//endregion
